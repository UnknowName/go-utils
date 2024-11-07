package gelf

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel uint8

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	timeFormat = "2006/1/2 15:04:05.000"
)

var log *Log

func init() {
	once := sync.Once{}
	if log == nil {
		once.Do(func() {
			log = new(Log)
			log.handlers = make(map[string]LogHandler)
		})
	}
}

type LogHandler interface {
	write(msg string)
	name() string
}

type Log struct {
	level    LogLevel
	handlers map[string]LogHandler
}

func NewLog() *Log {
	if log == nil {
		log = new(Log)
		log.handlers = make(map[string]LogHandler)
	}
	return log
}

func (l *Log) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Log) AddHandlers(handlers []LogHandler) {
	for _, handler := range handlers {
		if _, ok := l.handlers[handler.name()]; !ok {
			l.handlers[handler.name()] = handler
		}
	}
}

func (l *Log) Debug(msg string) {
	if l.level <= DEBUG {
		for _, handler := range l.handlers {
			handler.write(msg)
		}
	}
}

func (l *Log) Info(msg string) {
	if l.level <= INFO {
		for _, handler := range l.handlers {
			handler.write(msg)
		}
	}
}

func (l *Log) Warn(msg string) {
	if l.level <= WARN {
		for _, handler := range l.handlers {
			handler.write(msg)
		}
	}
}

func (l *Log) Error(msg string) {
	if l.level <= ERROR {
		for _, handler := range l.handlers {
			handler.write(msg)
		}
	}
}

type ConsoleHandler struct {
}

func NewConsoleHandler() *ConsoleHandler {
	return &ConsoleHandler{}
}

func (c *ConsoleHandler) name() string {
	return "console"
}

func (c *ConsoleHandler) write(msg string) {
	_formatMsg := c.formatMsg(msg)
	_, _ = os.Stdout.Write([]byte(_formatMsg))
}

func (c *ConsoleHandler) formatMsg(msg string) string {
	logTime := time.Now().Format(timeFormat)
	_, _file, line, _ := runtime.Caller(2)
	_vars := strings.Split(_file, "/")
	file := _vars[len(_vars)-1]
	_msg := fmt.Sprintf("%v [%v:%v] %v\n", logTime, file, line, msg)
	return _msg
}
