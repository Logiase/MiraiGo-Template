# MiraiGo-Template
A template for MiraiGo

[![Go Report Card](https://goreportcard.com/badge/github.com/Logiase/MiraiGo-Template)](https://goreportcard.com/report/github.com/Logiase/MiraiGo-Template)

基于[MiraiGo](https://github.com/Mrs4s/MiraiGo)的多模块组合设计

包装了基础功能,同时设计了一个~~良好~~的项目结构

## 基础配置

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

## 将 [MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template) 作为go module使用

可参考当前 [app.go](./app.go) 将其引入

使用这种方法可以引入其他小伙伴编写的第三方module

### 第三方 Module

暂无 ~~呜呜呜~~

## 引入的第三方 go module

 - [viper](https://github.com/spf13/viper)
    用于解析配置文件，同时可监听配置文件的修改
 - [logrus](github.com/sirupsen/logrus)
    功能丰富的Logger
 - [asciiart](github.com/yinghau76/go-ascii-art)
    用于在console显示图形验证码