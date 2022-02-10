package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func GetToken() string {
	return os.Getenv("DDING_TOKEN")
}

// env SKIPS=hostname:ItemName
// eg: SKIPS=log:内存,es:memory
// 表示针对主机名称包含log或者es且指标为内存或memory的告警的直接忽略

type Skip struct {
	HostName string
	ItemName string
}

func GetSkips() []*Skip {
	skips := make([]*Skip, 0)
	skip := strings.Split(os.Getenv("SKIPS"), ",")
	for i := range skip {
		_skip := strings.Split(skip[i], ":")
		skips = append(skips, &Skip{_skip[0], _skip[1]})
	}
	return skips
}

func GetSkipKey() string {
	return os.Getenv("SKIP_KEY")
}

// 在这里判断有些忽略的机器，如GrayLog故意让内存使用很高，当收到这个时，直接忽略掉

func CreateMsg (alert *PrometheusAlert, keywords string) *TMessage {
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
			for k, v := range _alert.Labels {
				if k == "job" || k == "severity" || k == "alertname" {
					continue
				}
				var ok bool
				for _, keyword := range strings.Split(keywords, ",") {
					ok, _ = regexp.MatchString(keyword, v)
					if ok {
						break
					}
				}
				if keywords != "" && ok {
					hostInfo = ""
					break
				} else {
					hostInfo = fmt.Sprintf("\t%s %s: %s ", hostInfo, k, v)
				}

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