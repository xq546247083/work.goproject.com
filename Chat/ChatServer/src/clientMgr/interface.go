package clientMgr

import (
	. "work.goproject.com/Chat/ChatServerModel/src"
)

// 客户端连接接口
type IClient interface {
	// 获取客户端对象的唯一标识
	GetId() int32

	// 获取玩家Id
	GetPlayerId() string

	// 玩家登录
	PlayerLogin(string)

	// 追加发送数据
	AppendSendData(*ServerResponseObject)

	// 客户端连接超时
	Expired() bool

	// 客户端连接对象退出登录
	Quit()
}
