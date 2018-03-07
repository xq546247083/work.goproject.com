package config

import (
	"encoding/json"

	"work.goproject.com/goutil/configUtil"
	"work.goproject.com/goutil/debugUtil"
)

type BaseConfig struct {
	// 为ChatServer提供服务的监听地址
	ChatServerAddress string

	// 为Web接口提供服务的监听地址
	WebServerAddress string

	// GoPs监控程序监听地址
	GopsAddr string
}

func (this *BaseConfig) String() string {
	bytes, _ := json.Marshal(this)
	return string(bytes)
}

var (
	baseConfig *BaseConfig
)

func initBaseConfig(config *configUtil.XmlConfig) error {
	tempConfig := new(BaseConfig)
	err := config.Unmarshal("root/BaseConfig", tempConfig)
	if err != nil {
		return err
	}

	baseConfig = tempConfig
	debugUtil.Printf("baseConfig:%v\n", baseConfig)

	return nil
}

// GetBaseConfig 获取服务器基础配置
func GetBaseConfig() *BaseConfig {
	return baseConfig
}
