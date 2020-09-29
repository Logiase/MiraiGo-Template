package bot

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/sirupsen/logrus"

	asc2art "github.com/yinghau76/go-ascii-art"

	"MiraiGo-Template/config"
	"MiraiGo-Template/utils"
)

type Bot struct {
	*client.QQClient

	start bool
}

var Instance *Bot

var logger = logrus.WithField("bot", "internal")

func init() {
	Instance = &Bot{
		client.NewClient(
			config.GlobalConfig.GetInt64("bot.account"),
			config.GlobalConfig.GetString("bot.password"),
		),
		false,
	}
	client.SystemDeviceInfo.ReadJson(utils.ReadFile("./device.json"))
}

func Login() {
	resp, err := Instance.Login()
	console := bufio.NewReader(os.Stdin)
	for {
		if err != nil {
			logger.WithError(err).Fatal("unable to login")
		}

		if !resp.Success {
			switch resp.Error {
			case client.NeedCaptcha:
				img, _, _ := image.Decode(bytes.NewReader(resp.CaptchaImage))
				fmt.Println(asc2art.New("image", img).Art)
				logger.Warn("captcha: ")
				text, _ := console.ReadString('\n')
				resp, err = Instance.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), resp.CaptchaSign)
				continue
			case client.UnsafeDeviceError:
				logger.Warnf("device lock -> %v", resp.VerifyUrl)
				return
			case client.OtherLoginError, client.UnknownLoginError:
				logger.Fatalf("login failed: %v", resp.ErrorMessage)
			}
		}
		break
	}
	logger.Info("bot login: %s", Instance.Nickname)
}

func RefreshList() {
	logger.Info("start reload friends list")
	err := Instance.ReloadFriendList()
	if err != nil {
		logger.WithError(err).Error("unable to load friends list")
	}
	logger.Infof("load %d friends", len(Instance.FriendList))
	logger.Info("start reload groups list")
	err = Instance.ReloadGroupList()
	if err != nil {
		logger.WithError(err).Error("unable to load groups list")
	}
	logger.Infof("load %d groups", len(Instance.GroupList))
}

func StartService() {
	if Instance.start {
		return
	}

	Instance.start = true

	logger.Infof("initializing modules ...")
	for _, mi := range modules {
		mi.Instance.Init()
	}
	logger.Info("all modules initialized")

	logger.Info("registering modules serve functions ...")
	for _, mi := range modules {
		mi.Instance.Serve(Instance)
	}
	logger.Info("all modules serve functions registered")

	logger.Info("starting modules tasks ...")
	for _, mi := range modules {
		go mi.Instance.Start(Instance)
	}
	logger.Info("tasks running")
}

func Stop() {
	logger.Warn("stopping ...")
	wg := sync.WaitGroup{}
	for _, mi := range modules {
		wg.Add(1)
		mi.Instance.Stop(Instance, &wg)
	}
	wg.Wait()
	logger.Info("stopped")
}
