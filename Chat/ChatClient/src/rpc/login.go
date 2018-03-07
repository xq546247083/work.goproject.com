package rpc

import (
	"work.goproject.com/Chat/ChatClient/src/config"
)

var (
	// 存储登陆成功信息的通道
	loginSucceedCh = make(chan int)
)

func login() {
	baseConfig := config.GetBaseConfig()

	requestMap := make(map[string]interface{})
	requestMap["MethodName"] = "Login"
	requestMap["Parameters"] = []interface{}{baseConfig.PlayerId, baseConfig.Token, baseConfig.PartnerId, baseConfig.ServerId}

	// 先登陆服务器
	request(requestMap)
}
