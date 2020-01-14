# Go utils

## Examples

- `log.go`
```go
package main

import (
    utils "github.com/unknowname/goutils"
)


var log *utils.Log

func init() {
    log = utils.NewLog()
    level := utils.DEBUG
    log.SetLevel(level)
    consoleHandler := utils.NewConsoleHandler()
    gelf := utils.NewGELFHandler("128.0.255.10", 12201)
    log.AddHandlers(consoleHandler, gelf)
}


func main() {
    log.Debug("msg from log")
}

```