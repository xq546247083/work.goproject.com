package channelMgr

import (
	"encoding/json"
	"fmt"

	"work.goproject.com/Chat/ChatServer/src/bll/dbConfig"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/debugUtil"
)

var (
	channelMap = make(map[string]IChannel)
)

const (
	con_Private_Channel = "Private"
)

func register(channelObj IChannel) {
	// 判断聊天频道是否开启
	if dbConfig.IsChannelExists(channelObj.Channel()) == false {
		return
	}

	// 判断是不已经存在，避免重复注册
	if _, exists := channelMap[channelObj.Channel()]; exists {
		panic(fmt.Sprintf("%s has already existed, please choose another name.", channelObj.Channel()))
	}

	// 初始化数据
	var err error
	if err = channelObj.InitConfig(); err != nil {
		panic(err)
	}
	if err = channelObj.InitBaseChannel(); err != nil {
		panic(err)
	}
	if err = channelObj.InitHistoryMgr(); err != nil {
		panic(err)
	}
	if err = channelObj.InitHistory(); err != nil {
		panic(err)
	}

	// 注册重新加载的方法
	reloadMgr.RegisterReloadFunc(fmt.Sprintf("channelMgr.%s.reload", channelObj.Channel()), channelObj.ReloadConfig)

	// 将对象注册到集合中
	channelMap[channelObj.Channel()] = channelObj
	fmt.Printf("there are %d channel registered.\n", len(channelMap))
}

func initConfig(channelObj IChannel) (err error) {
	configObj := new(config)

	configStr, err := dbConfig.GetStringConfig(channelObj.ConfigName())
	if err != nil {
		err = fmt.Errorf("Couldn't find config whose name should be %s, please check.", channelObj.ConfigName())
		return
	}

	err = json.Unmarshal([]byte(configStr), configObj)
	if err != nil {
		err = fmt.Errorf("Unmarshal config err whose name is: %s, content: %s, err is:%s, please check.", channelObj.ConfigName(), configStr, err)
		return
	}

	channelObj.SetConfig(configObj)
	debugUtil.Printf("%s:%v\n", channelObj.ConfigName(), configObj)
	return
}

func GetChannel(channel string) (channelObj IChannel, exists bool) {
	channelObj, exists = channelMap[channel]
	return
}

func GetPrivateChannel() (channelObj IChannel, exists bool) {
	channelObj, exists = channelMap[con_Private_Channel]
	return
}
