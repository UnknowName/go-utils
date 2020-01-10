package goutils

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int8

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	timeFormat = "2006/1/2 15:04:05.000"
)


type LogHandler interface {
	write(msg string)
}


type Log struct {
	Level LogLevel
	Handlers []LogHandler
}

var log *Log

func init() {
	if log == nil {
		log = new(Log)
	}
}

func NewLog() *Log {
	if log == nil {
		log = new(Log)
	}
	return log
}

func (l *Log) SetLevel(level LogLevel) {
	l.Level = level
}

func (l *Log) AddHandler(h LogHandler) {
	exist := false
	if l.Handlers == nil {
		l.Handlers = make([]LogHandler, 0)
	} else {
		for _, handler := range l.Handlers {
			if handler == h {
				exist = true
				break
			}
		}
	}
	if !exist {
		l.Handlers = append(l.Handlers, h)
	}
}

func (l *Log) Debug(msg string) {
	logTime := time.Now().Format(timeFormat)
	if l.Level <= DEBUG {
		for _, handler := range l.Handlers {
			_formatLog := fmt.Sprintf("%v INFO %v\n", logTime, msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Info(msg string) {
	logTime := time.Now().Format(timeFormat)
	if l.Level <= INFO {
		for _, handler := range l.Handlers {
			_formatLog := fmt.Sprintf("%v INFO %v\n", logTime, msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Warn(msg string) {
	if l.Level <= WARN {
		for _, handler := range l.Handlers {
			handler.write("WARN: " + msg)
		}
	}
}

func (l *Log) Error(msg string) {
	if l.Level <= ERROR {
		for _, handler := range l.Handlers {
			handler.write("ERROR: " + msg)
		}
	}
}


type ConsoleHandler struct {

}

func (c *ConsoleHandler) write(msg string) {
	_, _ = os.Stdout.Write([]byte(msg))
}