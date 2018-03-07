package rpcClient

import (
	"encoding/json"
	"time"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/goroutineMgr"
)

var (
	ForwardChatMessageChannel = make(chan *ForwardChatMessage, 1024*100)
)

func init() {
	// 处理数据
	go func() {
		// 处理goroutine数量
		goroutineName := "rpcClient.handleData"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		for {
			select {
			case chatMessageObj := <-ForwardChatMessageChannel:
				forward(chatMessageObj)
			default:
				// 如果channel中没有数据，则休眠5毫秒
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
}

// 转发聊天消息
func forward(chatMessageObj *ForwardChatMessage) {
	message, _ := json.Marshal(chatMessageObj)
	params := make([]interface{}, 1, 1)
	params[0] = string(message)

	//发送请求
	request(Con_ChatServerCenter_Method_Forward, params, nil)
}
