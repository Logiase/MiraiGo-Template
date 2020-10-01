package autoreply

import (
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"gopkg.in/yaml.v2"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &ar{}
var logger = utils.GetModuleLogger("logiase.autoreply")
var tem map[string]string

type ar struct {
}

func (a *ar) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "logiase.autoreply",
		Instance: instance,
	}
}

func (a *ar) Init() {
	bytes := utils.ReadFile("./autoreply.yaml")
	err := yaml.Unmarshal(bytes, &tem)
	if err != nil {
		logger.WithError(err).Errorf("unable to read autoreply.yaml")
	}
}

func (a *ar) PostInit() {
}

func (a *ar) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		out := autoreply(msg.ToString())
		if out == "" {
			return
		}
		m := message.NewSendingMessage().Append(message.NewText(out))
		c.SendGroupMessage(msg.GroupCode, m)
	})

	b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
		out := autoreply(msg.ToString())
		if out == "" {
			return
		}
		m := message.NewSendingMessage().Append(message.NewText(out))
		c.SendPrivateMessage(msg.Sender.Uin, m)
	})
}

func (a *ar) Start(bot *bot.Bot) {
}

func (a *ar) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
}

func autoreply(in string) string {
	out, ok := tem[in]
	if !ok {
		return ""
	}
	return out
}
