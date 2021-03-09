package bot

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/Logiase/MiraiGo-Template/v2/logging"
	"github.com/Mrs4s/MiraiGo/client"
	asc2art "github.com/yinghau76/go-ascii-art"
)

// Bot client.QQClient 包装
type Bot struct {
	*client.QQClient

	start bool

	logger logging.ILogger
}

/******************************************
 **           Global Instance            **
 ******************************************/

// Instance 全局Bot实例
// 单 Bot 设计
var Instance *Bot

// InitBot 初始化全局默认Bot
func InitBot(b *Bot) {
	Instance = b
}

// LoginP 使用Instance进行登录
func LoginP() {
	Instance.LoginP()
}

// RefreshList 使用Instance刷新联系人
func RefreshList() {
	Instance.RefreshList()
}

// StartService 启动服务
//
// 请勿重复调用
func StartService() {
	Instance.StartService()
}

/******************************************
 **                 Bot                  **
 ******************************************/

// NewBot 创建新Bot实例
// 使用默认设置
func NewBot(account int64, password string) *Bot {
	return &Bot{
		client.NewClient(account, password),
		false,
		logging.NewLogger(strconv.FormatInt(account, 10)),
	}
}

// NewBotWithLogger 创建新Bot实例
func NewBotWithLogger(account int64, password string, logger logging.ILogger) *Bot {
	return &Bot{
		client.NewClient(account, password),
		false,
		logger,
	}
}

// NewBotMD5 创建新Bot实例
func NewBotMD5(account int64, passwordMD5 [16]byte) *Bot {
	return &Bot{
		client.NewClientMd5(account, passwordMD5),
		false,
		logging.NewLogger(strconv.FormatInt(account, 10)),
	}
}

// LoginP 登录
// (因为函数名冲突)
func (b *Bot) LoginP() {
	resp, err := b.Login()
	console := bufio.NewReader(os.Stdin)

	for {
		if err != nil {
			b.logger.Fatalf("unable to login: %v", err)
		}

		var text string
		if !resp.Success {
			switch resp.Error {

			case client.NeedCaptcha:
				img, _, _ := image.Decode(bytes.NewReader(resp.CaptchaImage))
				fmt.Println(asc2art.New("image", img).Art)
				fmt.Print("please input captcha: ")
				text, _ := console.ReadString('\n')
				resp, err = b.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), resp.CaptchaSign)
				continue

			case client.UnsafeDeviceError:
				fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
				os.Exit(4)

			case client.SMSNeededError:
				fmt.Println("device lock enabled, Need SMS Code")
				fmt.Printf("Send SMS to %s ? (yes)", resp.SMSPhone)
				t, _ := console.ReadString('\n')
				t = strings.TrimSpace(t)
				if t != "yes" {
					os.Exit(2)
				}
				if !b.RequestSMS() {
					b.logger.Fatalf("unable to request SMS Code")
				}
				fmt.Printf("please input SMS Code: ")
				text, _ = console.ReadString('\n')
				resp, err = b.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
				continue

			case client.TooManySMSRequestError:
				b.logger.Fatalf("too many SMS request, please try later.\n")

			case client.SMSOrVerifyNeededError:
				fmt.Println("device lock enabled, choose way to verify:")
				fmt.Println("1. Send SMS Code to ", resp.SMSPhone)
				fmt.Println("2. Scan QR Code")
				fmt.Print("input (1,2):")
				text, _ = console.ReadString('\n')
				text = strings.TrimSpace(text)
				switch text {
				case "1":
					if !b.RequestSMS() {
						fmt.Println("unable to request SMS Code")
						os.Exit(2)
					}
					fmt.Print("please input SMS Code: ")
					text, _ = console.ReadString('\n')
					resp, err = b.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
					continue
				case "2":
					fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
					os.Exit(2)
				default:
					fmt.Println("invalid input")
					os.Exit(2)
				}

			case client.SliderNeededError:
				if client.SystemDeviceInfo.Protocol == client.AndroidPhone {
					fmt.Println("Android Phone Protocol DO NOT SUPPORT Slide verify")
					fmt.Println("please use other protocol")
					os.Exit(2)
				}
				b.AllowSlider = false
				b.Disconnect()
				resp, err = b.Login()
				continue

			case client.OtherLoginError, client.UnknownLoginError:
				b.logger.Fatalf("login failed: %v", resp.ErrorMessage)
			}

		}

		break
	}
	b.logger.Infof("bot login: %s", b.Nickname)
}

// RefreshList 刷新联系人
func (b *Bot) RefreshList() {
	b.logger.Infof("start reload friends list")
	err := Instance.ReloadFriendList()
	if err != nil {
		b.logger.Errorf("unable to load friends list: %v", err)
	}
	b.logger.Infof("load %d friends", len(Instance.FriendList))
	b.logger.Infof("start reload groups list")
	err = Instance.ReloadGroupList()
	if err != nil {
		b.logger.Errorf("unable to load groups list: %v", err)
	}
	b.logger.Infof("load %d groups", len(Instance.GroupList))
}

// StartService 启动服务
//
// 请勿重复调用
func (b *Bot) StartService() {

}

// // StartService 启动服务
// // 根据 Module 生命周期 此过程应在Login前调用
// // 请勿重复调用
// func StartService() {
// 	if Instance.start {
// 		return
// 	}

// 	Instance.start = true

// 	logger.Infof("initializing modules ...")
// 	for _, mi := range modules {
// 		mi.Instance.Init()
// 	}
// 	for _, mi := range modules {
// 		mi.Instance.PostInit()
// 	}
// 	logger.Info("all modules initialized")

// 	logger.Info("registering modules serve functions ...")
// 	for _, mi := range modules {
// 		mi.Instance.Serve(Instance)
// 	}
// 	logger.Info("all modules serve functions registered")

// 	logger.Info("starting modules tasks ...")
// 	for _, mi := range modules {
// 		go mi.Instance.Start(Instance)
// 	}
// 	logger.Info("tasks running")
// }

// // Stop 停止所有服务
// // 调用此函数并不会使Bot离线
// func Stop() {
// 	logger.Warn("stopping ...")
// 	wg := sync.WaitGroup{}
// 	for _, mi := range modules {
// 		wg.Add(1)
// 		mi.Instance.Stop(Instance, &wg)
// 	}
// 	wg.Wait()
// 	logger.Info("stopped")
// 	modules = make(map[string]ModuleInfo)
// }
