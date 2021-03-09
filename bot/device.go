package bot

import (
	"io/ioutil"
	"os"

	"github.com/Logiase/MiraiGo-Template/v2/logging"
	"github.com/Logiase/MiraiGo-Template/v2/utils"
	"github.com/Mrs4s/MiraiGo/client"
)

// UseDevice 使用 device 进行初始化设备信息
func UseDevice(device []byte) error {
	return client.SystemDeviceInfo.ReadJson(device)
}

// GenRandomDevice 生成随机设备信息
func GenRandomDevice() {
	client.GenRandomDevice()
	b, _ := utils.FileExist("./device.json")
	if b {
		logging.InternalLogger.Warnf("device.json exists, will not write device to file")
	}
	err := ioutil.WriteFile("device.json", client.SystemDeviceInfo.ToJson(), os.FileMode(0755))
	if err != nil {
		logging.InternalLogger.Errorf("unable to write device.json: %v", err)
	}
}
