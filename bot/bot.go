package bot

import (
	"bufio"
	"bytes"
	"fmt"
	_ "image/png"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/pkg/errors"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/tuotoo/qrcode"

	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/sirupsen/logrus"
)

// Bot 全局 Bot
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
	deviceJSONContent := utils.ReadFile("./device.json")
	InitWithDeviceJSONContent(deviceJSONContent)
}

type InitOption struct {
	Account           int64
	Password          string
	DeviceJSONContent []byte //cannot be nil if using option init
}

func InitWithOption(option InitOption) error {
	Instance = &Bot{
		QQClient: client.NewClient(
			option.Account,
			option.Password,
		),
		start: false,
	}

	device := new(client.DeviceInfo)
	err := device.ReadJson(option.DeviceJSONContent)
	if err != nil {
		return errors.Errorf("failed to apply device.json with err:%s", err)
	}
	Instance.UseDevice(device)
	return nil
}

func InitWithDeviceJSONContent(deviceJSONContent []byte) {
	var account = config.GlobalConfig.GetInt64("bot.account")
	var password = config.GlobalConfig.GetString("bot.password")
	err := InitWithOption(InitOption{
		Account:           account,
		Password:          password,
		DeviceJSONContent: deviceJSONContent,
	})
	if err != nil {
		panic(err)
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
	deviceInfo := new(client.DeviceInfo)
	err := deviceInfo.ReadJson(device)
	if err != nil {
		return err
	}
	Instance.UseDevice(deviceInfo)
	return nil
}

// GenRandomDevice 生成随机设备信息
func GenRandomDevice() {
	device := client.GenRandomDevice()
	b, _ := utils.FileExist("./device.json")
	if b {
		logger.Warn("device.json exists, will not write device to file")
		return
	}
	err := os.WriteFile("device.json", device.ToJson(), os.FileMode(0644))
	if err != nil {
		logger.WithError(err).Errorf("unable to write device.json")
	}
}

// SaveToken 会话缓存
func SaveToken() {
	AccountToken := Instance.GenToken()
	_ = os.WriteFile("session.token", AccountToken, 0o644)
}

type LoginMethod string

const (
	LoginMethodToken  = "token"
	LoginMethodQRCode = "qrcode"
	LoginMethodCommon = "common"
)

// Login 登录
func Login() error {
	var tokenData []byte = nil
	// 存在token缓存的情况快速恢复会话
	if exist, _ := utils.FileExist("./session.token"); exist {
		logger.Infof("检测到会话缓存, 尝试快速恢复登录")
		token, err := os.ReadFile("./session.token")
		if err != nil {
			return fmt.Errorf("failed to read token from file with err: %w", err)
		}
		tokenData = token
	}
	fmt.Println(Instance.Uin)
	var loginMethodValue = config.GlobalConfig.GetString("bot.login-method")
	return LoginWithOption(LoginOption{
		LoginMethod:              LoginMethod(loginMethodValue),
		Token:                    tokenData,
		UseTokenWhenUnmatchedUin: true,
	})
}

type LoginOption struct {
	LoginMethod              LoginMethod
	Token                    []byte //if not nil, try with most priority
	UseTokenWhenUnmatchedUin bool
}

func LoginWithOption(option LoginOption) error {
	if option.Token != nil {
		err := func() error {
			logger.Infof("检测到会话缓存, 尝试快速恢复登录")
			var token = option.Token
			r := binary.NewReader(token)
			cu := r.ReadInt64()
			if Instance.Uin != 0 {
				if cu != Instance.Uin && !option.UseTokenWhenUnmatchedUin {
					return fmt.Errorf("配置文件内的QQ号 (%v) 与会话缓存内的QQ号 (%v) 不相同", Instance.Uin, cu)
				}
			}
			if err := Instance.TokenLogin(token); err != nil {
				time.Sleep(time.Second)
				Instance.Disconnect()
				return errors.Errorf("恢复会话失败(%s)", err)
			} else {
				logger.Infof("快速恢复登录成功")
				return nil
			}
		}()
		if err != nil {
			logger.WithError(err).Warn("failed restore session by token, fallback to common or qrcode")
		} else {
			return nil
		}
	}
	switch option.LoginMethod {
	case LoginMethodCommon:
		return CommonLogin()
	case LoginMethodQRCode:
		return QrcodeLogin()
	default:
		return errors.New("unknown login method")
	}
}

// CommonLogin 普通账号密码登录
func CommonLogin() error {
	res, err := Instance.Login()
	if err != nil {
		return err
	}
	return loginResponseProcessor(res)
}

// QrcodeLogin 扫码登陆
func QrcodeLogin() error {
	rsp, err := Instance.FetchQRCode()
	if err != nil {
		return err
	}
	fi, err := qrcode.Decode(bytes.NewReader(rsp.ImageData))
	if err != nil {
		return err
	}
	_ = os.WriteFile("qrcode.png", rsp.ImageData, 0o644)
	defer func() { _ = os.Remove("qrcode.png") }()
	if Instance.Uin != 0 {
		logger.Infof("请使用账号 %v 登录手机QQ扫描二维码 (qrcode.png) : ", Instance.Uin)
	} else {
		logger.Infof("请使用手机QQ扫描二维码 (qrcode.png) : ")
	}
	time.Sleep(time.Second)
	qrcodeTerminal.New().Get(fi.Content).Print()
	s, err := Instance.QueryQRCodeStatus(rsp.Sig)
	if err != nil {
		return err
	}
	prevState := s.State
	for {
		time.Sleep(time.Second)
		s, _ = Instance.QueryQRCodeStatus(rsp.Sig)
		if s == nil {
			continue
		}
		if prevState == s.State {
			continue
		}
		prevState = s.State
		switch s.State {
		case client.QRCodeCanceled:
			logger.Fatalf("扫码被用户取消.")
		case client.QRCodeTimeout:
			logger.Fatalf("二维码过期")
		case client.QRCodeWaitingForConfirm:
			logger.Infof("扫码成功, 请在手机端确认登录.")
		case client.QRCodeConfirmed:
			res, err := Instance.QRCodeLogin(s.LoginInfo)
			if err != nil {
				return err
			}
			return loginResponseProcessor(res)
		case client.QRCodeImageFetch, client.QRCodeWaitingForScan:
			// ignore
		}
	}
}

// ErrSMSRequestError SMS请求出错
var ErrSMSRequestError = errors.New("sms request error")

var console = bufio.NewReader(os.Stdin)

func readLine() (str string) {
	str, _ = console.ReadString('\n')
	str = strings.TrimSpace(str)
	return
}

func readLineTimeout(t time.Duration, de string) (str string) {
	r := make(chan string)
	go func() {
		select {
		case r <- readLine():
		case <-time.After(t):
		}
	}()
	str = de
	select {
	case str = <-r:
	case <-time.After(t):
	}
	return
}

// loginResponseProcessor 登录结果处理
func loginResponseProcessor(res *client.LoginResponse) error {
	var err error
	for {
		if err != nil {
			return err
		}
		if res.Success {
			return nil
		}
		var text string
		switch res.Error {
		case client.SliderNeededError:
			logger.Warnf("登录需要滑条验证码, 请使用手机QQ扫描二维码以继续登录.")
			Instance.Disconnect()
			Instance.Release()
			Instance.QQClient = client.NewClientEmpty()
			return QrcodeLogin()
		case client.NeedCaptcha:
			logger.Warnf("登录需要验证码.")
			_ = os.WriteFile("captcha.jpg", res.CaptchaImage, 0o644)
			logger.Warnf("请输入验证码 (captcha.jpg)： (Enter 提交)")
			text = readLine()
			_ = os.Remove("captcha.jpg")
			res, err = Instance.SubmitCaptcha(text, res.CaptchaSign)
			continue
		case client.SMSNeededError:
			logger.Warnf("账号已开启设备锁, 按 Enter 向手机 %v 发送短信验证码.", res.SMSPhone)
			readLine()
			if !Instance.RequestSMS() {
				logger.Warnf("发送验证码失败，可能是请求过于频繁.")
				return errors.WithStack(ErrSMSRequestError)
			}
			logger.Warn("请输入短信验证码： (Enter 提交)")
			text = readLine()
			res, err = Instance.SubmitSMS(text)
			continue
		case client.SMSOrVerifyNeededError:
			logger.Warnf("账号已开启设备锁，请选择验证方式:")
			logger.Warnf("1. 向手机 %v 发送短信验证码", res.SMSPhone)
			logger.Warnf("2. 使用手机QQ扫码验证.")
			logger.Warn("请输入(1 - 2) (将在10秒后自动选择2)：")
			text = readLineTimeout(time.Second*10, "2")
			if strings.Contains(text, "1") {
				if !Instance.RequestSMS() {
					logger.Warnf("发送验证码失败，可能是请求过于频繁.")
					return errors.WithStack(ErrSMSRequestError)
				}
				logger.Warn("请输入短信验证码： (Enter 提交)")
				text = readLine()
				res, err = Instance.SubmitSMS(text)
				continue
			}
			fallthrough
		case client.UnsafeDeviceError:
			logger.Warnf("账号已开启设备锁，请前往 -> %v <- 验证后重启Bot.", res.VerifyUrl)
			logger.Infof("按 Enter 或等待 5s 后继续....")
			readLineTimeout(time.Second*5, "")
			os.Exit(0)
		case client.OtherLoginError, client.UnknownLoginError, client.TooManySMSRequestError:
			msg := res.ErrorMessage
			if strings.Contains(msg, "版本") {
				msg = "密码错误或账号被冻结"
			}
			if strings.Contains(msg, "冻结") {
				logger.Fatalf("账号被冻结")
			}
			logger.Warnf("登录失败: %v", msg)
			logger.Infof("按 Enter 或等待 5s 后继续....")
			readLineTimeout(time.Second*5, "")
			os.Exit(0)
		}
	}
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
	for _, mi := range modules {
		mi.Instance.PostInit()
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
