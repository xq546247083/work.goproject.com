package rpcClient

import (
	"encoding/json"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

// 处理消息
func handleMessage(id int32, message []byte) {
	debugUtil.Printf("id:%d,message:%s\n", id, string(message))

	// 如果id=0表示是服务器主动推送过来的消息，否则是客户端请求后的信息返回
	if id == 0 {
		// 将返回结果反序列化
		forwardObj := new(ForwardObject)
		if err := json.Unmarshal(message, forwardObj); err != nil {
			logUtil.ErrorLog("反序列化%s出错，错误信息为：%s", string(message), err)
			return
		}

		handleActiveMess(forwardObj)
	} else {
		// 将返回结果反序列化
		responseObj := new(CenterResponseObject)
		if err := json.Unmarshal(message, responseObj); err != nil {
			logUtil.ErrorLog("反序列化%s出错，错误信息为：%s", string(message), err)
			return
		}

		handlePassiveMess(id, responseObj)
	}
}

// 处理由服务器主动推送过来的消息
func handleActiveMess(forwardObj *ForwardObject) {
	// 由HandleCenterMessage方法来进行处理
	if centerMessageHandler == nil {
		panic("handleCenterMessage尚未被赋值")
	}

	go centerMessageHandler(forwardObj)
}

// 处理由客户端发送给服务器,再由服务器反馈的消息
func handlePassiveMess(id int32, responseObj *CenterResponseObject) {
	callbackFunc, exists := getCallbackFunc(id)
	if !exists {
		debugUtil.Println("receive response is invalid data, id is :", id)
		return
	}

	defer func() {
		deleteCallbackFunc(id)
	}()

	// 返回成功，则调用指定的回调方法；否则表示一些提示、警告、或者版本、资源更新等信息；否则表示其它信息的返回
	if responseObj.Code == Success.Code {
		if callbackFunc == nil {
			return
		}

		// 调用对应的回调函数
		callbackFunc(responseObj.Data)
	} else {
		// 处理特殊的返回值
		switch responseObj.Code {
		default:
			logUtil.ErrorLog("receive response from server failed，the error info is：", responseObj.Message)
		}
	}
}

func updateOnlineCount() {
	// 参数赋值
	params := make([]interface{}, 2, 2)
	if getClientCountHandler == nil {
		params[0] = 0
	} else {
		params[0] = getClientCountHandler()
	}
	if getPlayerCountHandler == nil {
		params[1] = 0
	} else {
		params[1] = getPlayerCountHandler()
	}

	// 发送OnlineCount消息
	debugUtil.Println("Send updateOnlineCount")

	// 发送请求
	request(Con_ChatServerCenter_Method_UpdateOnlineCount, params, nil)
}
