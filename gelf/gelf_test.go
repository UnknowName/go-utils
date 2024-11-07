package gelf

import (
	"fmt"
	"testing"
)

func TestNewGELFHandler(t *testing.T) {
	server := "172.18.61.78"
	port := 12202
	msg := `test msg`
	fmt.Println(len(msg))
	gelf := NewGELFHandler(server, port)
	gelf.AddProperty("source", "cheng-pc5")
	gelf.write(msg)
}
