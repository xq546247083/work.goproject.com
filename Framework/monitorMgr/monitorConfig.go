package monitorMgr

import (
	"encoding/json"
)

// 监控配置对象
type MonitorConfig struct {
	// 监控使用的服务器IP
	ServerIp string

	// 监控使用的服务器名称
	ServerName string

	// 监控的时间间隔（单位：分钟）
	Interval int
}

func (this *MonitorConfig) String() string {
	bytes, _ := json.Marshal(this)
	return string(bytes)
}

func NewMonitorConfig(serverIp, serverName string, interval int) *MonitorConfig {
	return &MonitorConfig{
		ServerIp:   serverIp,
		ServerName: serverName,
		Interval:   interval,
	}
}
