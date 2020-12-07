package reminder

import (
	"sync"

	"github.com/Logiase/MiraiGo-Template/config"
	"gopkg.in/yaml.v2"

	//"github.com/6tail/lunar-go/calendar"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
)

func init() {
	bot.RegisterModule(instance)
}

var instance = &reminder{}
var logger = utils.GetModuleLogger("jueyanyingyu.reminder")
var tem map[string]string

type reminder struct {
}

func (r *reminder) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "jueyanyingyu.reminder",
		Instance: instance,
	}
}

func (r *reminder) Init() {
	path := config.GlobalConfig.GetString("jueyanyingyu.reminder.path")

	if path == "" {
		path = "./reminder.yaml"
	}

	bytes := utils.ReadFile(path)
	err := yaml.Unmarshal(bytes, &tem)
	if err != nil {
		logger.WithError(err).Errorf("unable to read config file in %s", path)
	}
}

func (r *reminder) PostInit() {
}

func (r *reminder) Serve(b *bot.Bot) {

}

func (r *reminder) Start(bot *bot.Bot) {
}

func (r *reminder) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}
