package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknowname/webhook-dding/utils"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	sendFmt     = "https://oapi.dingtalk.com/robot/send?access_token="
	contentType = "application/json"
)

func main() {
	r := gin.Default()
	r.POST("/ping", send)
	r.Run("0.0.0.0:8080")
}

func send(c *gin.Context) {
	postData, err := ioutil.ReadAll(c.Request.Body)
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
	msg := utils.CreateMsg(alert, utils.GetSkipKey())
	if msg != nil {
		go func() {
			url := fmt.Sprintf("%s%s", sendFmt, utils.GetToken())
			httpClient := http.Client{Timeout: time.Second * 5}
			resp, err := httpClient.Post(url, contentType, bytes.NewBuffer(msg.Encode()))
			if err != nil {
				log.Println("send error ", err)
			} else {
				defer resp.Body.Close()
				resp, _ := ioutil.ReadAll(resp.Body)
				log.Println(string(resp))
			}
		}()
	} else {
		log.Println("匹配到关键字", utils.GetSkipKey(), "此次告警将不会发送钉钉通知")
	}
	c.JSON(200, gin.H{"message": "ok"})
}