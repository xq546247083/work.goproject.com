package clientMgr

import (
	"work.goproject.com/Chat/ChatServer/src/rpcClient"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

var (
	// 查找player方法
	getPlayerHandler func(string, bool) (*Player, bool, error)
)

func init() {
	rpcClient.RegisterGetClientCountHandler(getClientCount)
}

func RegisterGetPlayerHandler(handler func(string, bool) (*Player, bool, error)) {
	getPlayerHandler = handler
}
