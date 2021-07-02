package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

// DDing Text Message

type TMessage struct {
	MsgType string  `json:"msgtype"`
	Text    Content `json:"text"`
	At      At      `json:"at"`
	IsAtAll bool    `json:"isAtAll"`
}

type Content struct {
	Content string `json:"content"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
}

func NewTMessage(msg string, atMobiles []string, atAll bool) *TMessage {
	if atMobiles == nil {
		atMobiles = make([]string, 0)
	}
	atUsers := At{AtMobiles: atMobiles}
	text := Content{Content: msg}
	return &TMessage{
		MsgType: "text",
		Text:    text,
		At:      atUsers,
		IsAtAll: atAll,
	}
}

// Struct to JSON format

func (tm *TMessage) Encode() []byte {
	bytes, err := json.Marshal(&tm)
	if err != nil {
		fmt.Println("Encode to json err ", err)
		return nil
	}
	return bytes
}

// AlertManager alert struct

func NewPrometheusAlert() *PrometheusAlert {
	return &PrometheusAlert{}
}

type PrometheusAlert struct {
	Status string  `json:"status"`
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Status     string            `json:"status"`
	Labels     map[string]string `json:"labels"`
	Annotation map[string]string `json:"annotations"`
	Start      time.Time         `json:"startsAt"`
	End        time.Time         `json:"endsAt"`
}

// 暂时没用上

type Label struct {
	Key   string
	Value string
}

func (pa *PrometheusAlert) Decode(data []byte) error {
	return json.Unmarshal(data, pa)
}
