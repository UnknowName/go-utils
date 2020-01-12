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
	Level LogLevel
	Handlers []LogHandler
}

func NewLog() *Log {
	if log == nil {
		log = new(Log)
	}
	return log
}

func (l *Log) formatMsg(msg string) string {
	var level string
	switch l.Level {
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
	if l.Level <= DEBUG {
		for _, handler := range l.Handlers {
			_formatLog := l.formatMsg(msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Info(msg string) {
	if l.Level <= INFO {
		for _, handler := range l.Handlers {
			_formatLog := l.formatMsg(msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Warn(msg string) {
	if l.Level <= WARN {
		for _, handler := range l.Handlers {
			_formatLog := l.formatMsg(msg)
			handler.write(_formatLog)
		}
	}
}

func (l *Log) Error(msg string) {
	if l.Level <= ERROR {
		for _, handler := range l.Handlers {
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