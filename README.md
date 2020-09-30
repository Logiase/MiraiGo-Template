# MiraiGo-Template
A template for MiraiGo

[![Go Report Card](https://goreportcard.com/badge/github.com/Logiase/MiraiGo-Template)](https://goreportcard.com/report/github.com/Logiase/MiraiGo-Template)

多模块组合设计

## 基础

账号配置[application.yaml](./application.yaml)
```yaml
bot:
  # 账号
  account: 1234567
  # 密码
  password: example
```

## Module 配置

module参考[log.go](./modules/logging/log.go)

```go
package mymodule

import (
    "aaa"
    "bbb"
    "MiraiGo-Template/bot"
)

var instance *Logging

func init() {
	instance = &Logging{}
	bot.RegisterModule(instance)
}

type Logging struct {
}

// ...
```

编写自己的Module后在[app.go](./app.go)中启用Module 

```go
package main

import (
    // ...
    
    _ "modules/mymodule"
)

// ...
```