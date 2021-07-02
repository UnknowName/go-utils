# gelf

## Examples

```go
package main

import "goutils/gelf"

var log *gelf.Log

func init() {
	if log == nil {
		log = gelf.NewLog()
		level := gelf.DEBUG
		log.SetLevel(level)
		consoleHandler := gelf.NewConsoleHandler()
		gelfHandler := gelf.NewGELFHandler("128.0.255.10", 12201)
		gelfHandler.AddProperty("env", "test")
		log.AddHandlers(consoleHandler, gelfHandler)
	}
}

func main() {
	log.Info("test msg from gelf log")
}

```