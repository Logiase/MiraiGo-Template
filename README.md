# MiraiGo-Template

A template for MiraiGo

> v2 版本正在重写，请提出 *你的想法* 或 *你对当前设计的不满*  

[![Go Report Card](https://goreportcard.com/badge/github.com/Logiase/MiraiGo-Template)](https://goreportcard.com/report/github.com/Logiase/MiraiGo-Template)

基于 [MiraiGo](https://github.com/Mrs4s/MiraiGo) 的多模块组合设计

包装了基础功能,同时设计了一个~~良好~~的项目结构

## 不了解go?

golang 极速入门

[点我看书](https://github.com/justjavac/free-programming-books-zh_CN#go)

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

## 快速入门

你可以克隆本项目, 或者将本项目作为依赖.

在开始之前, 你需要首先生成设备文件.

新建文件 `tools_test.go` , 内容如下:

```go
package main_test

import (
	"testing"

	"github.com/Logiase/MiraiGo-Template/bot"
)

func TestGenDevice(t *testing.T) {
	bot.GenRandomDevice()
}
```

然后运行 `TestGenDevice` 来生成一份设备文件

### 克隆

如果你克隆本项目, 请首先更新项目依赖, 同步到协议库最新版本, 否则可能出现某些意外的bug ( 或产生新的bug )

```go
go get -u
```

### 将 [MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template) 作为go module使用

可参考当前 [app.go](./app.go) 将其引入

使用这种方法可以引入其他小伙伴编写的第三方module

## 内置 Module

 - internal.logging
 将收到的消息按照格式输出至 os.stdout

## 第三方 Module

欢迎PR

 - [logiase.autoreply](https://github.com/Logiase/MiraiGo-module-autoreply)
 按照收到的消息进行回复
 
## 进阶内容 

### Docker 支持

参照 [Dockerfile](./Dockerfile)

## 引入的第三方 go module

 - [MiraiGo](https://github.com/Mrs4s/MiraiGo)
    核心协议库
 - [viper](https://github.com/spf13/viper)
    用于解析配置文件，同时可监听配置文件的修改
 - [logrus](github.com/sirupsen/logrus)
    功能丰富的Logger
 - [asciiart](github.com/yinghau76/go-ascii-art)
    用于在console显示图形验证码
