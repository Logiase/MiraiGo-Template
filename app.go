package main

import (
	"os"
	"os/signal"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"

	_ "github.com/Logiase/MiraiGo-Template/modules/logging"
)

func init() {
	config.Init()
}

func main() {
	bot.Init()

	bot.StartService()
	bot.Login()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}
