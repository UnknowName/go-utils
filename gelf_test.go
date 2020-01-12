package goutils

import (
	"fmt"
	"testing"
)

func TestNewGELFHandler(t *testing.T) {
	server := "128.0.255.10"
	port := 12201
	connect := fmt.Sprintf("%v:%v", server, port)
	fmt.Println(connect == "128.0.255.10:12201")
	gelf := NewGELFHandler(server, port)
	gelf.AddProperty("from", "cheng-pc")
	gelf.write("test msg 1")
}
