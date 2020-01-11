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
    consoleHandler := &utils.ConsoleHandler{}
    log.AddHandler(consoleHandler)
}


func main() {
    log.Debug("msg from log")
}

```