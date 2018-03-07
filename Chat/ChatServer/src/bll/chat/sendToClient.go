package chat

import (
	"time"

	"work.goproject.com/Chat/ChatServer/src/clientMgr"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
)

// 发送在另一台设备登陆的信息
// clientObj：客户端对象
func sendLoginAnotherDeviceMsg(clientObj clientMgr.IClient) {
	responseObj := NewServerResponseObject()
	responseObj.SetMethodName("Login")
	responseObj.SetResultStatus(LoginOnAnotherDevice)

	// 先发送消息，然后再断开连接
	SendToClient(clientObj, responseObj)

	// 启动独立goroutine来发断开连接
	go func() {
		time.Sleep(2 * time.Second)

		// 客户端断开
		clientObj.Quit()
	}()
}

// 发送数据给客户端
// player：玩家对象
// responseObj：Socket服务器的返回对象
func SendToClient(clientObj clientMgr.IClient, responseObj *ServerResponseObject) {
	clientMgr.ResponseResult(clientObj, responseObj)
}

// 发送数据给玩家
// playerList：玩家列表
// responseObj：Socket服务器的返回对象
func SendToPlayer(playerList []*Player, responseObj *ServerResponseObject) {
	for _, item := range playerList {
		debugUtil.Printf("item.ClientId:%d\n", item.ClientId)
		if item.ClientId > 0 {
			clientObj, exists := clientMgr.GetClient(item.ClientId)
			debugUtil.Printf("item.ClientId:%d,exists:%v\n", item.ClientId, exists)
			if exists {
				SendToClient(clientObj, responseObj)
			}
		}
	}
}
