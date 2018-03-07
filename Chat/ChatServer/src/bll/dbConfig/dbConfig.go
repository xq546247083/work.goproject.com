package dbConfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"work.goproject.com/Chat/ChatServer/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/managecenterMgr"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/typeUtil"
	"work.goproject.com/goutil/stringUtil"
)

var (
	configData typeUtil.MapData

	// ManageCenter配置对象
	ManageCenterConfig *managecenterMgr.ManageCenterConfig

	// 聊天频道配置
	ChannelConfig = make(map[string]bool, 8)

	// 忽略内容配置
	OmitContentConfig = make([]string, 0, 8)
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

	var channelConfigStr string
	if channelConfigStr, err = configData.String("ChannelConfig"); err != nil {
		return
	} else {
		for _, v := range strings.Split(channelConfigStr, ",") {
			ChannelConfig[v] = true
		}
		debugUtil.Printf("ChannelConfig:%v\n", ChannelConfig)
	}

	var omitContentConfigStr string
	if omitContentConfigStr, err = configData.String("OmitContentConfig"); err != nil {
		return
	} else {
		for _, v := range strings.Split(omitContentConfigStr, "||") {
			if(!stringUtil.IsEmpty(v)){
				OmitContentConfig = append(OmitContentConfig, v)
			}
		}
		debugUtil.Printf("OmitContentConfig:%v\n", OmitContentConfig)
	}

	return
}

func GetStringConfig(key string) (value string, err error) {
	value, err = configData.String(key)
	return
}

func IsChannelExists(channel string) bool {
	_, exists := ChannelConfig[channel]
	return exists
}

func IsOmit(str string) bool {
	for _, item := range OmitContentConfig {
		if strings.HasPrefix(str, item) {
			return true
		}
	}

	return false
}
