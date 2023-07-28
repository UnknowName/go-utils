package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknowname/webhook-dding/utils"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	sendFmt     = "https://oapi.dingtalk.com/robot/send?access_token="
	contentType = "application/json"
	defautSlice = time.Hour * 24
)

var recorder map[string]time.Time

func init() {
	recorder = make(map[string]time.Time)
}

func main() {
	r := gin.Default()
	r.POST("/ping", send)
	r.Run("0.0.0.0:8080")
}

func send(c *gin.Context) {
	postData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Read post data error ", err)
		c.JSON(500, gin.H{"message": "Read post data error"})
		return
	}
	log.Println("接收到AlertManager的消息:", string(postData))
	// 将字节类型的告警信息转换成Golang struct类型
	alert := utils.NewPrometheusAlert()
	if err := alert.Decode(postData); err != nil {
		log.Println("Struct error ", err)
		c.JSON(500, gin.H{"message": "Decode error"})
		return
	}
	// 详细告警信息在alert.Alerts里面
	msg := utils.CreateMsg(alert, utils.GetSkips())
	if msg != nil {
		key := msg.Text.Content
		latest, ok := recorder[key]
		if ok && latest.After(time.Now()) {
			// 消息静默期，后续不执行
			return
		}
		recorder[msg.Text.Content] = time.Now().Add(defautSlice)
		go func() {
			url := fmt.Sprintf("%s%s", sendFmt, utils.GetToken())
			secret := utils.GetSecret()
			if secret != "" {
				url = fmt.Sprintf("%s%s", url, utils.GetSignature(secret))
			}
			httpClient := http.Client{Timeout: time.Second * 5}
			resp, err := httpClient.Post(url, contentType, bytes.NewBuffer(msg.Encode()))
			if err != nil {
				log.Println("send error ", err)
			} else {
				defer resp.Body.Close()
				resp, _ := io.ReadAll(resp.Body)
				log.Println(string(resp))
			}
		}()
	} else {
		log.Println("匹配到关键字", utils.GetSkips(), "此次告警将不会发送钉钉通知")
	}
	c.JSON(200, gin.H{"message": "ok"})
}