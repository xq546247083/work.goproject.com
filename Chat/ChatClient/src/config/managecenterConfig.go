package config

import (
	"work.goproject.com/Framework/managecenterMgr"
	"work.goproject.com/goutil/configUtil"
	"work.goproject.com/goutil/debugUtil"
)

var (
	manageCenterConfig *managecenterMgr.ManageCenterConfig
)

func initManageCenterConfig(config *configUtil.XmlConfig) error {
	tmpConfig := new(managecenterMgr.ManageCenterConfig)
	err := config.Unmarshal("root/ManageCenterConfig", tmpConfig)
	if err != nil {
		return err
	}

	manageCenterConfig = tmpConfig
	debugUtil.Println("ManageCenterConfig:", manageCenterConfig)

	return nil
}

func GetManageCenterConfig() *managecenterMgr.ManageCenterConfig {
	return manageCenterConfig
}
