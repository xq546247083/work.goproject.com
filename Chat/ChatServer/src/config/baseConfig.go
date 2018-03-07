package config

import (
	"encoding/json"
	"fmt"

	"work.goproject.com/goutil/configUtil"
	"work.goproject.com/goutil/debugUtil"
)

// 基础配置对象
type BaseConfig struct {
	// ChatCenter监听地址
	ChatCenterAddress string

	// 通信协议tcp/websocket
	Protocol string

	// 内网地址
	PrivateIP string

	// 公网地址
	PublicIP string

	// 为ChatServer提价服务的端口
	ChatServerPort string

	// 为GameServer提供服务的端口
	GameServerPort string

	// 为GameServer提供服务的Web端口
	GameServerWebPort string

	// 是否压缩返回给客户端的数据
	IfCompressData bool

	// GoPs监控程序监听地址
	GopsAddr string
}

func (this *BaseConfig) GetPrivateChatServerAddress() string {
	return fmt.Sprintf("%s:%s", this.PrivateIP, this.ChatServerPort)
}

func (this *BaseConfig) GetPrivateGameServerAddress() string {
	return fmt.Sprintf("%s:%s", this.PrivateIP, this.GameServerPort)
}

func (this *BaseConfig) GetPrivateGameServerWebAddress() string {
	return fmt.Sprintf("%s:%s", this.PrivateIP, this.GameServerWebPort)
}

func (this *BaseConfig) GetPrivateGopsAddress() string {
	return fmt.Sprintf("%s:%s", this.PrivateIP, this.GopsAddr)
}

func (this *BaseConfig) GetPublicChatServerAddress() string {
	return fmt.Sprintf("%s:%s", this.PublicIP, this.ChatServerPort)
}

func (this *BaseConfig) GetPublicGameServerAddress() string {
	return fmt.Sprintf("%s:%s", this.PublicIP, this.GameServerPort)
}

func (this *BaseConfig) GetPublicGameServerWebAddress() string {
	return fmt.Sprintf("http://%s:%s/API/player/login", this.PublicIP, this.GameServerWebPort)
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

	if baseConfig.Protocol != "tcp" && baseConfig.Protocol != "websocket" {
		panic("Protocol Error, it should be either tcp or websocket")
	}

	return nil
}

// GetBaseConfig 获取服务器基础配置
func GetBaseConfig() *BaseConfig {
	return baseConfig
}
