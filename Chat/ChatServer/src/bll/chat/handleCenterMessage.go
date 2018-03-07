package chat

import (
	"work.goproject.com/Chat/ChatServer/src/bll/player"
	"work.goproject.com/Chat/ChatServer/src/rpcClient"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/logUtil"
)

func init() {
	rpcClient.RegisterCenterMessageHandler(handleCenterMessage)
}

// 处理消息（来自于ChatServerCenter的消息）
func handleCenterMessage(forwardObj *ForwardObject) {
	// 处理goroutine数量
	goroutineName := "handleCenterMessage"
	goroutineMgr.MonitorZero(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	// 根据TransferType来选择数据
	switch forwardObj.ForwardType {
	case Con_ChatMessage:
		handleChatMessage(forwardObj.ForwardChatMessage)
	case Con_UpdatePlayerInfo:
		handleUpdatePlayerInfo(forwardObj.ForwardPlayer)
	case Con_Reload:
		handleReload()
	default:
		logUtil.ErrorLog("从Center收到了未定义的类型%s", forwardObj.ForwardType)
	}
}

func handleChatMessage(chatMessageObj *ForwardChatMessage) {
	// 设置为已经转发，避免再次转发
	chatMessageObj.IsTransfered = true
	chatMessageChannel <- chatMessageObj
}

func handleUpdatePlayerInfo(forwardPlayerObj *ForwardPlayer) {
	player.UpdatePlayerInfo(forwardPlayerObj)
}

func handleReload() {
	errList := reloadMgr.Reload()
	if errList != nil && len(errList) > 0 {
		for _, err := range errList {
			logUtil.ErrorLog("Reload Err:%s", err)
		}
	}
}
