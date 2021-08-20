package bot

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/tuotoo/qrcode"
	asc2art "github.com/yinghau76/go-ascii-art"
	"image"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/sirupsen/logrus"
)

var reloginLock = new(sync.Mutex)

const sessionToken = "session.token"

// Bot 全局 Bot
type Bot struct {
	*client.QQClient

	start    bool
	isQRCode bool
}

func (bot *Bot) saveToken() {
	_ = ioutil.WriteFile(sessionToken, bot.GenToken(), 0o677)
}
func (bot *Bot) clearToken() {
	os.Remove(sessionToken)
}

func (bot *Bot) getToken() ([]byte, error) {
	return ioutil.ReadFile(sessionToken)
}

// ReLogin 掉线时可以尝试使用会话缓存重新登陆，只允许在OnDisconnected中调用
func (bot *Bot) ReLogin(e *client.ClientDisconnectedEvent) error {
	reloginLock.Lock()
	defer reloginLock.Unlock()
	if bot.Online {
		return nil
	}
	logger.Warnf("Bot已离线: %v", e.Message)
	logger.Warnf("尝试重连...")
	token, err := bot.getToken()
	if err == nil {
		err = bot.TokenLogin(token)
		if err == nil {
			bot.saveToken()
			return nil
		}
	}
	logger.Warnf("快速重连失败: %v", err)
	if bot.isQRCode {
		logger.Errorf("快速重连失败, 扫码登录无法恢复会话.")
		return errors.New("qrcode login relogin failed")
	}
	logger.Warnf("快速重连失败, 尝试普通登录. 这可能是因为其他端强行T下线导致的.")
	time.Sleep(time.Second)

	resp, err := bot.Login()
	if err != nil {
		logger.Errorf("登录时发生致命错误: %v", err)
		return err
	}
	err = login(resp)
	if err == nil {
		bot.saveToken()
	}
	return err
}

// Instance Bot 实例
var Instance *Bot

var logger = logrus.WithField("bot", "internal")

// Init 快速初始化
// 使用 config.GlobalConfig 初始化账号
// 使用 ./device.json 初始化设备信息
func Init() {
	account := config.GlobalConfig.GetInt64("bot.account")
	password := config.GlobalConfig.GetString("bot.password")

	InitBot(account, password)

	deviceJson := utils.ReadFile("./device.json")
	if deviceJson == nil {
		logger.Fatal("can not read ./device.json")
	}
	err := client.SystemDeviceInfo.ReadJson(deviceJson)
	if err != nil {
		logger.WithError(err).Fatal("read device.json error")
	}
}

// InitBot 使用 account password 进行初始化账号
func InitBot(account int64, password string) {
	if account == 0 {
		Instance = &Bot{
			QQClient: client.NewClientEmpty(),
			isQRCode: true,
		}
	} else {
		Instance = &Bot{
			QQClient: client.NewClient(account, password),
		}
	}
}

// UseDevice 使用 device 进行初始化设备信息
func UseDevice(device []byte) error {
	return client.SystemDeviceInfo.ReadJson(device)
}

// GenRandomDevice 生成随机设备信息
func GenRandomDevice() {
	client.GenRandomDevice()
	b, _ := utils.FileExist("./device.json")
	if b {
		logger.Warn("device.json exists, will not write device to file")
		return
	}
	err := ioutil.WriteFile("device.json", client.SystemDeviceInfo.ToJson(), os.FileMode(0755))
	if err != nil {
		logger.WithError(err).Errorf("unable to write device.json")
	}
}

func qrcodeLogin() error {
	rsp, err := Instance.FetchQRCode()
	if err != nil {
		return err
	}
	fi, err := qrcode.Decode(bytes.NewReader(rsp.ImageData))
	if err != nil {
		return err
	}
	_ = ioutil.WriteFile("qrcode.png", rsp.ImageData, 0o644)
	defer func() { _ = os.Remove("qrcode.png") }()
	logger.Infof("请使用手机QQ扫描二维码 (qrcode.png) : ")
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
			logger.Info("扫码被用户取消.")
			os.Exit(1)
		case client.QRCodeTimeout:
			logger.Info("二维码过期")
			os.Exit(1)
		case client.QRCodeWaitingForConfirm:
			logger.Infof("扫码成功, 请在手机端确认登录.")
		case client.QRCodeConfirmed:
			res, err := Instance.QRCodeLogin(s.LoginInfo)
			if err != nil {
				return err
			}
			return login(res)
		case client.QRCodeImageFetch, client.QRCodeWaitingForScan:
			// ignore
		}
	}
}

// Login 登录
func Login() {
	Instance.AllowSlider = true
	if ok, _ := utils.FileExist(sessionToken); ok {
		token, err := Instance.getToken()
		if err != nil {
			goto NormalLogin
		}
		if Instance.Uin != 0 {
			r := binary.NewReader(token)
			sessionUin := r.ReadInt64()
			if sessionUin != Instance.Uin {
				logger.Warnf("QQ号(%v)与会话缓存内的QQ号(%v)不符，将清除会话缓存", Instance.Uin, sessionUin)
				Instance.clearToken()
				goto NormalLogin
			}
		}
		if err = Instance.TokenLogin(token); err != nil {
			Instance.clearToken()
			Instance.Disconnect()
			Instance.Release()
			Instance.QQClient = client.NewClientEmpty()
			logger.Warnf("恢复会话失败: %v , 尝试使用正常流程登录.", err)
			time.Sleep(time.Second)
		} else {
			Instance.saveToken()
			logger.Debug("恢复会话成功")
			return
		}
	}

NormalLogin:
	if Instance.Uin == 0 {
		logger.Info("未指定账号密码，请扫码登陆")
		err := qrcodeLogin()
		if err != nil {
			logger.Fatal("login failed: %v", err)
		} else {
			logger.Infof("bot login: %s", Instance.Nickname)
		}
	} else {
		logger.Info("使用帐号密码登陆")
		resp, err := Instance.Login()
		if err != nil {
			logger.Fatalf("login failed: %v", err)
		}

		err = login(resp)

		if err != nil {
			logger.Fatal("login failed: %v", err)
		} else {
			logger.Infof("bot login: %s", Instance.Nickname)
		}
	}
	Instance.saveToken()
}

var console = bufio.NewReader(os.Stdin)

var readLine = func() (str string) {
	str, _ = console.ReadString('\n')
	str = strings.TrimSpace(str)
	return
}

var readLineTimeout = func(t time.Duration, defaultV string) (str string) {
	r := make(chan string)
	go func() {
		select {
		case r <- readLine():
		case <-time.After(t):
		}
	}()
	str = defaultV
	select {
	case str = <-r:
	case <-time.After(t):
	}
	return
}

func login(resp *client.LoginResponse) error {
	var err error

	for {
		if err != nil {
			return err
		}
		if resp.Success {
			return nil
		}

		var text string
		switch resp.Error {
		case client.SliderNeededError:
			// code below copyright by https://github.com/Mrs4s/go-cqhttp
			logger.Warn("登录需要滑条验证码, 请使用手机QQ扫描二维码以继续登录.")
			Instance.Disconnect()
			Instance.Release()
			Instance.QQClient = client.NewClientEmpty()
			return qrcodeLogin()
		case client.NeedCaptcha:
			logger.Warn("登录需要验证码.")
			img, _, _ := image.Decode(bytes.NewReader(resp.CaptchaImage))
			fmt.Println(asc2art.New("image", img).Art)
			logger.Warn("请输入验证码 (captcha.jpg)： (Enter 提交)")
			text = readLine()
			resp, err = Instance.SubmitCaptcha(text, resp.CaptchaSign)
			continue
		case client.SMSNeededError:
			logger.Warnf("账号已开启设备锁, 按 Enter 向手机 %v 发送短信验证码.", resp.SMSPhone)
			readLine()
			if !Instance.RequestSMS() {
				logger.Warnf("发送验证码失败，可能是请求过于频繁.")
				return errors.New("sms send error")
			}
			logger.Warn("请输入短信验证码： (Enter 提交)")
			text = readLine()
			resp, err = Instance.SubmitSMS(text)
			continue
		case client.SMSOrVerifyNeededError:
			logger.Warn("账号已开启设备锁，请选择验证方式:")
			logger.Warnf("1. 向手机 %v 发送短信验证码", resp.SMSPhone)
			logger.Warn("2. 使用手机QQ扫码验证.")
			logger.Warn("请输入(1 - 2) (将在10秒后自动选择2)：")
			text = readLineTimeout(time.Second*10, "2")
			if strings.Contains(text, "1") {
				if !Instance.RequestSMS() {
					logger.Warnf("发送验证码失败，可能是请求过于频繁.")
					return errors.New("sms send error")
				}
				logger.Warn("请输入短信验证码： (Enter 提交)")
				text = readLine()
				resp, err = Instance.SubmitSMS(text)
				continue
			}
			fallthrough
		case client.UnsafeDeviceError:
			logger.Warnf("账号已开启设备锁，请前往 -> %v <- 验证后重启Bot.", resp.VerifyUrl)
			logger.Infof("按 Enter 或等待 5s 后继续....")
			readLineTimeout(time.Second*5, "")
			os.Exit(0)
		case client.OtherLoginError, client.UnknownLoginError, client.TooManySMSRequestError:
			msg := resp.ErrorMessage
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
