package config

import (
	"encoding/json"

	"work.goproject.com/goutil/configUtil"
	"work.goproject.com/goutil/debugUtil"
)

// 证书配置
type CrtConfig struct {
	// 钥匙
	Key string

	// 证书
	Crt string
}

func (this *CrtConfig) String() string {
	bytes, _ := json.Marshal(this)
	return string(bytes)
}

var (
	crtConfig *CrtConfig
)

func initCrtConfig(config *configUtil.XmlConfig) error {
	// 如果配置的websockets聊天通信，那么加载证书配置
	baseConfigTemp := GetBaseConfig()
	if baseConfigTemp.Protocol != "websockets"{
		return nil
	}

	tempConfig := new(CrtConfig)
	err := config.Unmarshal("root/CrtConfig", tempConfig)
	if err != nil {
		return err
	}

	crtConfig = tempConfig
	debugUtil.Printf("CrtConfig:%v\n", crtConfig)

	if crtConfig.Crt == "" || crtConfig.Key == "" {
		panic("证书配置错误，不能为空")
	}

	return nil
}

// GetCrtConfig 获取服务器证书配置
func GetCrtConfig() *CrtConfig {
	return crtConfig
}
