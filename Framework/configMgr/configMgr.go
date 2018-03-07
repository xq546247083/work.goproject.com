package configMgr

import (
	"work.goproject.com/goutil/configUtil"
)

// 配置管理对象
type ConfigManager struct {
	// 初始化方法列表
	initFuncList []func(*configUtil.XmlConfig) error
}

// 注册初始化方法
func (this *ConfigManager) RegisterInitFunc(initFunc func(*configUtil.XmlConfig) error) {
	this.initFuncList = append(this.initFuncList, initFunc)
}

// 初始化
func (this *ConfigManager) Init(configObj *configUtil.XmlConfig) error {
	for _, initFunc := range this.initFuncList {
		if err := initFunc(configObj); err != nil {
			return err
		}
	}

	return nil
}

// 创建配置管理对象
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		initFuncList: make([]func(*configUtil.XmlConfig) error, 0, 8),
	}
}
