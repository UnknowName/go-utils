package goutils

import (
	"testing"
)

func TestNewLog(t *testing.T) {
	log := NewLog()
	log.SetLevel(DEBUG)
	console := &ConsoleHandler{}
	gelf := NewGELFHandler("128.0.255.10", 12201)
	log.AddHandlers(console)
	log.AddHandlers(gelf, console)
	log.Info("test from msg on glef")
}
