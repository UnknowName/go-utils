package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func GetToken() string {
	return os.Getenv("DDING_TOKEN")
}

func GetSecret() string {
	return os.Getenv("DDING_SECRET")
}

// env SKIPS=hostname:ItemName
// eg: SKIPS=log:内存,es:memory
// 表示针对主机名称包含log或者es且指标为内存或memory的告警的直接忽略

type Skip struct {
	HostName string
	ItemName string
}

func (s *Skip) String() string {
	return fmt.Sprintf("Skip{HostName: %s,ItemName: %s", s.HostName, s.ItemName)
}

func GetSkips() []*Skip {
	skips := make([]*Skip, 0)
	skip := strings.Split(os.Getenv("SKIPS"), ",")
	for i := range skip {
		_skip := strings.Split(skip[i], ":")
		if len(_skip) != 2 {
			log.Println("environment SKIPS wrong: ", _skip)
			return nil
		}
		skips = append(skips, &Skip{_skip[0], _skip[1]})
	}
	return skips
}

func GetSkipKey() string {
	return os.Getenv("SKIP_KEY")
}

func CreateMsg (alert *PrometheusAlert, skips []*Skip) *TMessage {
	_status := ""
	if alert.Status == "firing" {
		_status = "告警"
	} else {
		_status = "恢复"
	}
	// 这个要放在后面定义，不然计数会有差异常，因为SKIP的不进入计数
	hosts := make([]string, 0)
	for _, _alert := range alert.Alerts {
		// 如果labels不为空，则取出所有的Labels
		hostInfo := ""
		if len(_alert.Labels) > 0 {
			// 先检查是否匹配主机名与指标，如果匹配，则此条目录跳过即可
			hostName := _alert.Labels["hostname"]
			itemName := _alert.Labels["alertname"]
			var okHostname, okItem bool
			for _, skip := range skips {
				// 为假就一直尝试匹配，直到匹配成功
				if !okHostname {
					okHostname, _ = regexp.MatchString(skip.HostName, hostName)
				}
				if !okItem {
					okItem, _ = regexp.MatchString(skip.ItemName, itemName)
				}
				// 如果两个都匹配成功，跳出匹配关键字循环
				if okHostname && okItem {
					break
				}
			}
			if okHostname && okItem {
				log.Printf("忽略主机: %s, 指标: %s 的告警", hostName, itemName)
				continue
			}
			// 继续的将相关指标显示出来
			for k, v := range _alert.Labels {
				if k == "job" || k == "severity" || k == "alertname" {
					continue
				}
				hostInfo = fmt.Sprintf("\t%s %s: %s ", hostInfo, k, v)
			}
		}
		// 尝试从annotation中取description的值，追加进主机列表信息中
		if value, ok := _alert.Annotation["description"]; ok && hostInfo != "" {
			hostInfo = fmt.Sprintf("%s 当前值: %.5s", hostInfo, value)
		}
		if hostInfo != "" {
			hosts = append(hosts, hostInfo)
		}
	}
	if len(hosts) < 1 {
		return nil
	}
	msg := fmt.Sprintf("异常名称: %s\n状态:  %s\n异常主机总数:  %d\n异常开始时间: %s\n异常主机列表: \n%s",
		alert.Alerts[0].Labels["alertname"],
		_status,
		len(hosts),
		alert.Alerts[0].Start.Local(),
		strings.Join(hosts, "\t\n"),
	)
	return NewTMessage(msg, nil, false)
}

func GetSignature(secret string) string {
	now := time.Now().UnixMilli()
	s := fmt.Sprintf("%d\n%s", now, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(s))
	sign := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	sign = url.QueryEscape(sign)
	return fmt.Sprintf("&timestamp=%d&sign=%s", now, sign)
}