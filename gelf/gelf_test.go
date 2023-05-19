package gelf

import (
	"testing"
)

func TestNewGELFHandler(t *testing.T) {
	server := "128.0.255.10"
	port := 12201
	// connect := fmt.Sprintf("%v:%v", server, port)
	// fmt.Println(connect == "128.0.21.56:12201")
	gelf := NewGELFHandler(server, port)
	gelf.AddProperty("source", "cheng-pc")
	gelf.write("INFO","test msg 1")
}
