# Go utils

## Examples

- `utils.log`
```go
package main

import (
    utils "github/unknowname/goutils"
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