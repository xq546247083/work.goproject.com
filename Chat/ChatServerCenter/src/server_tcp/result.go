package server_tcp

import (
	"encoding/json"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

// 推送信息给客户端
// clientList:客户端列表
// forwardObj:传输对象
func push(clientList []*client, forwardObj *ForwardObject) {
	var message []byte
	var err error
	if message, err = json.Marshal(forwardObj); err != nil {
		logUtil.ErrorLog("server_tcp.push序列化出错:%s", err)
		return
	}

	debugUtil.Println("push message:", string(message))
	for _, clientObj := range clientList {
		debugUtil.Printf("client:%s\n", clientObj.String())

		if forwardObj.ForwardType == Con_ChatMessage {
			// 判断是否需要排除该客户端，以避免重复处理
			if clientObj.chatServer != nil && clientObj.chatServer.chatServerAddress == forwardObj.ForwardChatMessage.ExcludeChatServer {
				continue
			}
		}

		debugUtil.Printf("push client:%s\n", clientObj.String())
		clientObj.appendSendData(newSendDataItem(0, message))
	}
}

// 发送响应结果
// clientObj:客户端对象
// requestObj:请求对象（如果为nil则代表是服务端主动推送信息，否则为客户端请求信息）
// responseObject:响应对象（不能为指针类型，否则在registerFunction时判断类型会出错）
func responseResult(clientObj *client, requestObj *CenterRequestObject, responseObj *CenterResponseObject) {
	message, err := json.Marshal(responseObj)
	if err != nil {
		logUtil.ErrorLog("server_tcp.responseResult序列化出错:%s", err)
		return
	}

	var sendDataItemObj *sendDataItem
	if requestObj == nil {
		sendDataItemObj = newSendDataItem(0, message)
	} else {
		sendDataItemObj = newSendDataItem(requestObj.Id, message)
	}

	clientObj.appendSendData(sendDataItemObj)
}
