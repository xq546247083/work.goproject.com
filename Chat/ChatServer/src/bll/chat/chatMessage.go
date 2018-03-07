package chat

import (
	"time"

	"work.goproject.com/Chat/ChatServer/src/bll/channelMgr"
	"work.goproject.com/Chat/ChatServer/src/rpcClient"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/debugUtil"
)

var (
	chatMessageChannel = make(chan *ForwardChatMessage, 1024)
)

func init() {
	// 处理数据
	go func() {
		// 处理goroutine数量
		goroutineName := "chat.handleChatMessage"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		for {
			select {
			case chatMessageObj := <-chatMessageChannel:
				go handleEachChatMessage(chatMessageObj)
			default:
				// 如果channel中没有数据，则休眠5毫秒
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
}

func handleEachChatMessage(chatMessageObj *ForwardChatMessage) {
	// 处理goroutine数量
	goroutineName := "chat.handleEachChatMessage"
	goroutineMgr.MonitorZero(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	debugUtil.Printf("chatMessageObj:%v\n", chatMessageObj)

	// 获取聊天频道对象
	channelObj, _ := channelMgr.GetChannel(chatMessageObj.Channel)

	// 保存历史消息
	channelObj.SaveHistory(chatMessageObj)

	// 获取目标玩家列表
	playerList := channelObj.GetPlayerList(chatMessageObj)
	debugUtil.Printf("player count:%d\n", len(playerList))

	// 构造返回值
	responseObj := NewServerResponseObject()
	responseObj.SetMethodName("SendMessage")
	responseObj.SetData(channelObj.NewServerResponseData(chatMessageObj))

	// 向玩家推送数据
	SendToPlayer(playerList, responseObj)

	// 如果是跨服聊天，且尚未向ChatServerCenter转发；则向ChatServerCenter转发数据
	if channelObj.IsCanCrossServer() && chatMessageObj.IsTransfered == false {
		rpcClient.ForwardChatMessageChannel <- chatMessageObj
	}
}
