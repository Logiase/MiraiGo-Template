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

	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
)

type Bot struct {
	*client.QQClient

	start bool
}

// Instance Bot 实例
var Instance *Bot

var logger = logrus.WithField("bot", "internal")

// Init 快速初始化
// 使用 config.GlobalConfig 初始化账号
// 使用 ./device.json 初始化设备信息
func Init() {
	Instance = &Bot{
		client.NewClient(
			config.GlobalConfig.GetInt64("bot.account"),
			config.GlobalConfig.GetString("bot.password"),
		),
		false,
	}
	err := client.SystemDeviceInfo.ReadJson(utils.ReadFile("./device.json"))
	if err != nil {
		logger.WithError(err).Panic("device.json error")
	}
}

// InitBot 使用 account password 进行初始化账号
func InitBot(account int64, password string) {
	Instance = &Bot{
		client.NewClient(account, password),
		false,
	}
}

// UseDevice 使用 device 进行初始化设备信息
func UseDevice(device []byte) error {
	return client.SystemDeviceInfo.ReadJson(device)
}

// Login 登录
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

// RefreshList 刷新联系人
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

// StartService 启动服务
// 根据 Module 生命周期 此过程应在Login前调用
// 请勿重复调用
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

// Stop 停止所有服务
// 调用此函数并不会使Bot离线
func Stop() {
	logger.Warn("stopping ...")
	wg := sync.WaitGroup{}
	for _, mi := range modules {
		wg.Add(1)
		mi.Instance.Stop(Instance, &wg)
	}
	wg.Wait()
	logger.Info("stopped")
	modules = make(map[string]ModuleInfo)
}
