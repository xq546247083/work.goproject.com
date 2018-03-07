package config

import (
	"fmt"

	"work.goproject.com/Framework/configMgr"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/configUtil"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

var (
	configManager = configMgr.NewConfigManager()
)

func init() {
	// 优先加基础配置
	configManager.RegisterInitFunc(initBaseConfig)
	configManager.RegisterInitFunc(initDBConfig)
	configManager.RegisterInitFunc(initMonitorConfig)
}

func init() {
	// 设置日志文件的存储目录
	logUtil.SetLogPath("LOG")

	if err := reload(); err != nil {
		panic(fmt.Errorf("初始化配置文件失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadMgr.RegisterReloadFunc("config.reload", reload)
}

func reload() error {
	// 读取配置文件内容
	configObj := configUtil.NewXmlConfig()
	err := configObj.LoadFromFile("config.xml")
	if err != nil {
		return err
	}

	debug, err := configObj.Bool("root/DEBUG", "")
	if err != nil {
		return err
	}

	// 设置debugUtil的状态
	debugUtil.SetDebug(debug)

	// 调用所有已经注册的配置初始化方法
	if err := configManager.Init(configObj); err != nil {
		return err
	}

	return nil
}
