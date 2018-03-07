package rpcClient

import (
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/logUtil"
)

var (
	// 处理中心发来的数据
	centerMessageHandler func(*ForwardObject)

	// 查找getPlayerCount方法
	getPlayerCountHandler func() int

	// 获取客户端数量的方法
	getClientCountHandler func() int
)

func RegisterCenterMessageHandler(handler func(*ForwardObject)) {
	centerMessageHandler = handler
	logUtil.DebugLog("rpcClient.RegisterCenterMessageHandler")
}

func RegisterGetPlayerCountHandler(handler func() int) {
	getPlayerCountHandler = handler
	logUtil.DebugLog("rpcClient.RegisterGetPlayerCountHandler")
}

func RegisterGetClientCountHandler(handler func() int) {
	getClientCountHandler = handler
	logUtil.DebugLog("rpcClient.RegisterGetClientCountHandler")
}
