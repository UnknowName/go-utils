package gelf

import (
	"testing"
)

func TestNewLog(t *testing.T) {
	log := NewLog()
	log.SetLevel(DEBUG)
	console := NewConsoleHandler()
	gelf := NewGELFHandler("128.0.255.10", 12201)
	gelf.AddProperty("source", "test")
	log.AddHandlers(console)
	log.AddHandlers(console, gelf)
	msg := `test msg`
	log.Info("test from msg on gelf")
	log.Info(msg)
}
