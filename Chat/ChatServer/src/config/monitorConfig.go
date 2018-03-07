package config

import (
	"work.goproject.com/Framework/monitorMgr"
	"work.goproject.com/goutil/configUtil"
	"work.goproject.com/goutil/debugUtil"
)

var (
	monitorConfig *monitorMgr.MonitorConfig
)

func initMonitorConfig(config *configUtil.XmlConfig) error {
	tempConfig := new(monitorMgr.MonitorConfig)
	err := config.Unmarshal("root/MonitorConfig", tempConfig)
	if err != nil {
		return err
	}

	monitorConfig = tempConfig
	debugUtil.Printf("monitorConfig:%v\n", monitorConfig)
	return nil
}

// GetMonitorConfig 获监测配置
func GetMonitorConfig() *monitorMgr.MonitorConfig {
	return monitorConfig
}
