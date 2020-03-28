package goutils

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

var (
	host     string
	port     string
	user     string
	password string
	attr     string
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "RabbitMQ Server address")
	flag.StringVar(&port, "port", "15672", "RabbitMQ API port")
	flag.StringVar(&user, "user", "guest", "RabbitMQ user")
	flag.StringVar(&password, "password", "guest", "RabbitMQ user password")
	flag.StringVar(&attr, "attr", "", "Get's RabbitMQ metric name")
}

func main() {
	flag.Parse()
	if attr == "" {
		fmt.Println("Use -h to get help. the attr is must present")
		return
	}
	mq := NewRabbitMQ(host, port, user, password)
	switch attr {
	case "channel":
		fmt.Printf("%d", mq.GetChannelCount())
	case "message":
		fmt.Printf("%d", mq.GetMessageCount())
	case "connect":
		fmt.Printf("%d", mq.GetConnectionCount())
	case "consumer":
		fmt.Printf("%d", mq.GetConsumerCount())
	default:
		fmt.Printf("%d", mq.getAttr(attr))
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		return
	}
}

func NewRabbitMQ(host, port, user, password string) *RabbitMQ {
	return &RabbitMQ{Host: host, user: user, port: port, password: password}
}

type RabbitMQ struct {
	Host     string
	port     string
	user     string
	password string
}

func (r *RabbitMQ) get(url string) (respByte []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.user, r.password)
	client := http.Client{Timeout: 5 * time.Second}
	respStruct, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer respStruct.Body.Close()
	respByte, err = ioutil.ReadAll(respStruct.Body)
	if respStruct.StatusCode != 200 {
		return nil, errors.New(string(respByte))
	}
	if err != nil {
		return nil, err
	}
	return respByte, err
}

func (r *RabbitMQ) getAttr(attr string) (total int) {
	apiPath := fmt.Sprintf("http://%s:%s/api/overview", r.Host, r.port)
	req, err := r.get(apiPath)
	failOnError(err, "Connect RabbitMQ failed")
	type detail struct {
		Queue      int `json:"queues"`
		Channel    int `json:"channels"`
		Connection int `json:"connections"`
		Consumer   int `json:"consumers"`
	}
	var info struct {
		Total detail `json:"object_totals"`
	}
	err = json.Unmarshal(req, &info)
	failOnError(err, "Json failed")
	value := reflect.ValueOf(info.Total).FieldByName(attr)
	if (value == reflect.Value{}) {
		total = -1
	} else {
		v := fmt.Sprintf("%v", value)
		total, _ = strconv.Atoi(v)
	}
	return total
}

func (r *RabbitMQ) GetMessageCount() (total int) {
	apiPath := fmt.Sprintf("http://%s:%s/api/vhosts", r.Host, r.port)
	req, err := r.get(apiPath)
	failOnError(err, "Get data from api failed")
	var msgs []struct {
		Message int `json:"messages_ready"`
	}
	err = json.Unmarshal(req, &msgs)
	failOnError(err, "Json failed")
	for _, msg := range msgs {
		total += msg.Message
	}
	return total
}

func (r *RabbitMQ) GetChannelCount() int {
	return r.getAttr("Channel")
}

func (r *RabbitMQ) GetConnectionCount() int {
	return r.getAttr("Connection")
}

func (r *RabbitMQ) GetConsumerCount() int {
	return r.getAttr("Consumer")
}