package clientMgr

import (
	"encoding/json"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

// 处理客户端请求
// clientObj：对应的客户端对象
// request：请求内容字节数组(json格式)
// 返回值：无
func HandleRequest(clientObj IClient, request []byte) {
	responseObj := NewServerResponseObject()

	defer func() {
		// 如果是客户端数据错误，则将客户端请求数据记录下来
		if responseObj.ResultStatus == ClientDataError {
			logUtil.ErrorLog("client:%s\n请求的数据为：%s, 返回的结果为客户端数据错误", clientObj, string(request))
		}

		// 调用发送消息接口
		ResponseResult(clientObj, responseObj)

		// 记录DEBUG日志
		if debugUtil.IsDebug() {
			result, _ := json.Marshal(responseObj)
			logUtil.DebugLog("client:%s\nRequest:%s\nResponse:%s", clientObj, string(request), string(result))
		}
	}()

	// 解析请求字符串
	requestObj := new(ServerRequestObject)
	if err := json.Unmarshal(request, requestObj); err != nil {
		logUtil.ErrorLog("反序列化出错，错误信息为：%s", err)
		responseObj.SetResultStatus(ClientDataError)
		return
	}

	// 为参数赋值
	responseObj.SetMethodName(requestObj.MethodName)

	// 对参数要特殊处理：将Client、Player特殊处理
	parameters := make([]interface{}, 0)
	if requestObj.MethodName == "Login" {
		parameters = append(parameters, interface{}(clientObj))
		parameters = append(parameters, requestObj.Parameters...)
	} else {
		// 判断玩家是否已经登陆
		if clientObj.GetPlayerId() == "" {
			responseObj.SetResultStatus(NoLogin)
			return
		}

		// 判断是否能找到玩家
		var playerObj *Player
		var exists bool
		var err error
		if getPlayerHandler == nil {
			panic("getPlayerHandler is nil, please set first")
		}
		playerObj, exists, err = getPlayerHandler(clientObj.GetPlayerId(), false)
		if err != nil {
			responseObj.SetResultStatus(DataError)
			return
		}

		if !exists {
			responseObj.SetResultStatus(NoLogin)
			return
		}

		parameters = append(parameters, interface{}(clientObj))
		parameters = append(parameters, interface{}(playerObj))
		parameters = append(parameters, requestObj.Parameters...)
	}

	// 为参数赋值
	requestObj.ModuleName = "Chat"
	requestObj.Parameters = parameters

	// 调用方法
	responseObj = callFunction(requestObj)
}
