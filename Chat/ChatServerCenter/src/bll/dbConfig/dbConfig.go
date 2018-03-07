package dbConfig

import (
	"encoding/json"
	"fmt"

	"work.goproject.com/Chat/ChatServerCenter/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/managecenterMgr"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/typeUtil"
)

var (
	configData typeUtil.MapData

	// 其它配置对象
	otherConfigObj *OtherConfig

	// ManageCenter配置对象
	ManageCenterConfig *managecenterMgr.ManageCenterConfig
)

func init() {
	if err := reload(); err != nil {
		panic(fmt.Errorf("初始化Config失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadMgr.RegisterReloadFunc("dbConfig.reload", reload)
}

func reload() (err error) {
	var configList []*DBConfig
	if err = dal.GetAll(&configList); err != nil {
		return
	}

	// 转换为MapData类型
	configMap := make(map[string]interface{})
	for _, item := range configList {
		configMap[item.ConfigKey] = item.ConfigValue
	}
	configData = typeUtil.NewMapData(configMap)

	// 解析其它配置
	var otherConfigStr string
	if otherConfigStr, err = configData.String("OtherConfig"); err != nil {
		return
	} else {
		if err = json.Unmarshal([]byte(otherConfigStr), &otherConfigObj); err != nil {
			return
		}
		debugUtil.Printf("otherConfigObj:%v\n", otherConfigObj)
	}

	// 解析ManageCenterConfig
	var manageCenterConfigStr string
	if manageCenterConfigStr, err = configData.String("ManageCenterConfig"); err != nil {
		return
	} else {
		if err = json.Unmarshal([]byte(manageCenterConfigStr), &ManageCenterConfig); err != nil {
			return
		}
		debugUtil.Printf("ManageCenterConfig:%v\n", ManageCenterConfig)
	}

	return
}

func GetStringConfig(key string) (value string, err error) {
	value, err = configData.String(key)
	return
}

func GetOtherConfig() *OtherConfig {
	return otherConfigObj
}
