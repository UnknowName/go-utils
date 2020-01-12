package goutils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int8
var log *Log

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


func init() {
	if log == nil {
		log = new(Log)
	}
}


type Log struct {
	level LogLevel
	handlers []LogHandler
}

func NewLog() *Log {
	if log == nil {
		log = new(Log)
	}
	return log
}

func (l *Log) formatMsg(msg string) string {
	var level string
	switch l.level {
	case 0:
		level = "DEBUG"
	case 1:
		level = "INFO"
	case 2:
		level = "WARN"
	case 3:
		level = "ERROR"
	default:
		level = "DEBUG"
	}
	logTime := time.Now().Format(timeFormat)
	_, _file, line, _ := runtime.Caller(2)
	_vars := strings.Split(_file, "/")
	file := _vars[len(_vars) - 1]
	_msg := fmt.Sprintf("%v [%v:%v] %v %v\n", logTime, file, line, level, msg)
	return _msg
}

func (l *Log) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Log) AddHandlers(hs...LogHandler) {
	if l.handlers == nil {
		l.handlers = make([]LogHandler, 0)
	}
	existHandler := len(l.handlers)
	if existHandler == 0 {
		l.handlers = append(l.handlers, hs...)
	} else {
		for _, newHandler := range hs {
			exist := false
			for _, existHandler := range l.handlers {
				if newHandler == existHandler {
					exist = true
					break
				}
				if !exist {
					l.handlers = append(l.handlers, newHandler)
				}
			}
		}
	}
}

func (l *Log) Debug(msg string) {
	if l.level <= DEBUG {
		for _, handler := range l.handlers {
			_formatLog := l.formatMsg(msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Info(msg string) {
	if l.level <= INFO {
		for _, handler := range l.handlers {
			_formatLog := l.formatMsg(msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Warn(msg string) {
	if l.level <= WARN {
		for _, handler := range l.handlers {
			_formatLog := l.formatMsg(msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Error(msg string) {
	if l.level <= ERROR {
		for _, handler := range l.handlers {
			_formatLog := l.formatMsg(msg)
			handler.write(_formatLog)
		}
	}
}


type ConsoleHandler struct {

}

func (c *ConsoleHandler) write(msg string) {
	_, _ = os.Stdout.Write([]byte(msg))
}