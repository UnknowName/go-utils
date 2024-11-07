package gelf

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	/*
	// Ethernet default MTU is 1500 Byte, when great then MTU while happen fragmentation.
	// The 1500 Byte is TCP/UDP max data size
	// So, Singer UDP Data size max is LAN MTU 1500 Byte - IP Header 20 Byte - UDPHeader 8 Byte = 1472 Byte
	*/

	UDPChunkSize = 1472
)

type GELFHandler struct {
	server 	string
	port   	int
	conn    *net.UDPConn
	logProperty map[string]interface{}
}

func NewGELFHandler(server string, port int) *GELFHandler {
	baseProperty := map[string]interface{}{"version": "1.1"}
	return &GELFHandler{server: server, port: port, logProperty: baseProperty}
}

func (g *GELFHandler) name() string {
	return "gelf"
}

func (g *GELFHandler) setLevel(level LogLevel) {
	g.AddProperty("level", level)
}

func (g *GELFHandler) AddProperty(key string, value interface{}) {
	g.logProperty[key] = value
}

func (g *GELFHandler) write(msg string) {
	logTime := time.Now().Format(timeFormat)
	g.logProperty["time"] = logTime
	g.logProperty["short_message"] = msg
	jsonMsg := g.toJson()
	if g.conn == nil {
		g.connect()
	}
	if _, err := g.conn.Write(jsonMsg); err != nil {
		fmt.Println("Send message to server error ", err)
	}
}

func (g *GELFHandler) toJson() []byte {
	bytes, err := json.Marshal(g.logProperty)
	if err != nil {
		fmt.Println("Parse JSON error ", err)
		return bytes
	}
	return bytes
}

func (g *GELFHandler) connect() {
	if g.conn == nil {
		addr := fmt.Sprintf("%v:%v", g.server, g.port)
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			fmt.Println("Resolve failed ", err)
			return
		}
		if conn, err := net.DialUDP("udp", nil, udpAddr); err != nil {
			fmt.Println("Connect failed ", err)
			return
		} else {
			g.conn = conn
		}
	}
}