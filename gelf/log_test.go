package gelf

import (
	"testing"
)

func TestNewLog(t *testing.T) {
	log := NewLog()
	log.SetLevel(DEBUG)
	console := NewConsoleHandler()
	gelf := NewGELFHandler("128.0.255.10", 12201)
	gelf.AddProperty("app", "test-log")
	log.AddHandlers(console)
	log.AddHandlers(console, gelf)
	log.Info("test from msg on gelf")
	log.Info("test from msg on gelf2")
}
