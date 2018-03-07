package server_tcp

import (
	"encoding/json"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

func init() {
	registerHandler(Con_ChatServerCenter_Method_Login, login, 3)
	registerHandler(Con_ChatServerCenter_Method_Forward, forward, 1)
	registerHandler(Con_ChatServerCenter_Method_UpdateOnlineCount, updateOnlineCount, 2)
}

func handleRequest(clientObj *client, id int32, request []byte) {
	responseObj := NewCenterResponseObject()

	// 提取请求内容
	requestObj := new(CenterRequestObject)
	if err := json.Unmarshal(request, requestObj); err != nil {
		logUtil.ErrorLog("server_tcp反序列化%s出错，错误信息为：%s", string(request), err)
		return
	}

	// 对requestObj的属性Id赋值
	requestObj.Id = id

	// 查找方法
	if handlerObj, exists := getHandler(requestObj.MethodName); !exists {
		responseResult(clientObj, requestObj, responseObj.SetResultStatus(MethodNotDefined))
		return
	} else {
		// 检测参数数量
		if rs := handlerObj.checkParamCount(requestObj.Parameters); rs != Success {
			responseResult(clientObj, requestObj, responseObj.SetResultStatus(rs))
			return
		}

		// 调用方法
		responseObj = handlerObj.handlerFunc(clientObj, requestObj.Parameters)

		// 输出结果
		responseResult(clientObj, requestObj, responseObj)
	}
}

func login(clientObj *client, parameters []interface{}) *CenterResponseObject {
	responseObj := NewCenterResponseObject()

	debugUtil.Printf("login param:%s,%s,%s\n", parameters[0], parameters[1], parameters[2])

	var chatServerAddress string
	var gameServerAddress string
	var gameServerWebAddress string
	var ok bool

	if chatServerAddress, ok = parameters[0].(string); !ok {
		return responseObj.SetResultStatus(ParamTypeError)
	}
	if gameServerAddress, ok = parameters[1].(string); !ok {
		return responseObj.SetResultStatus(ParamTypeError)
	}
	if gameServerWebAddress, ok = parameters[2].(string); !ok {
		return responseObj.SetResultStatus(ParamTypeError)
	}

	clientObj.login(chatServerAddress, gameServerAddress, gameServerWebAddress)

	return responseObj
}

func forward(clientObj *client, parameters []interface{}) *CenterResponseObject {
	responseObj := NewCenterResponseObject()

	if msg, ok := parameters[0].(string); !ok {
		return responseObj.SetResultStatus(ParamTypeError)
	} else {
		// 解析数据
		chatMessageObj := new(ForwardChatMessage)
		if err := json.Unmarshal([]byte(msg), chatMessageObj); err != nil {
			logUtil.ErrorLog("反序列化%s为ChatMessageObject出错，错误信息为：%s", msg, err)
			return responseObj.SetResultStatus(DataError)
		}

		// 排除当前服务器，以免重复处理
		if clientObj.chatServer != nil {
			// 此字段的赋值与result.push方法中的判断必须一致，否则会无效;
			chatMessageObj.ExcludeChatServer = clientObj.chatServer.chatServerAddress
		}

		// 然后转发给所有服务器
		ForwardObjectChannel <- NewForwardObject_ChatMessage(chatMessageObj)
	}

	return responseObj
}

func updateOnlineCount(clientObj *client, parameters []interface{}) *CenterResponseObject {
	responseObj := NewCenterResponseObject()
	clientCount := 0
	playerCount := 0

	if clientCount_float64, ok := parameters[0].(float64); !ok {
		return responseObj.SetResultStatus(ParamTypeError)
	} else {
		clientCount = int(clientCount_float64)
	}

	if playerCount_float64, ok := parameters[1].(float64); !ok {
		return responseObj.SetResultStatus(ParamTypeError)
	} else {
		playerCount = int(playerCount_float64)
	}

	debugUtil.Printf("updateOnlineCount:clientCount:%d, playerCount:%d\n", clientCount, playerCount)

	if clientObj.chatServer != nil {
		clientObj.updateOnlineCount(clientCount, playerCount)
	}

	return responseObj
}
