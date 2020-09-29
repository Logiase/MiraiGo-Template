package main

import (
	"os"
	"os/signal"

	"MiraiGo-Template/bot"

	// Modules On
	_ "MiraiGo-Template/modules/logging"
)

func main() {
	bot.StartService()
	bot.Login()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}
