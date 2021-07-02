package utils

import (
	"fmt"
	"testing"
)

func TestNewTMessage(t *testing.T) {
	msg := NewTMessage("test", nil, false)
	fmt.Println(msg)
}

func TestTMessage_Encode(t *testing.T) {
	msg := NewTMessage("test", nil, false)
	str := string(msg.Encode())
	fmt.Println(str)
}
