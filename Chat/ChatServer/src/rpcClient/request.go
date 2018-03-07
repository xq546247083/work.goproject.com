package rpcClient

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

var (
	// 请求Id:每个请求都会带上一个唯一Id，以便在接收到服务器的返回数据时能够区分出来自于不同的请求
	requestId int32 = 0

	// 回调方法集合，及其锁对象
	callbackFuncMap = make(map[int32]func(interface{}))
	requestMutex    sync.Mutex
)

// 注册回调方法
// id:自增Id
// callback:回调方法
func registerCallbackFunc(id int32, callbackFunc func(interface{})) {
	requestMutex.Lock()
	defer requestMutex.Unlock()

	callbackFuncMap[id] = callbackFunc
}

// 获取回调方法
// id:自增Id
// 返回值：
// 回调方法
func getCallbackFunc(id int32) (callbackFunc func(interface{}), exists bool) {
	requestMutex.Lock()
	defer requestMutex.Unlock()

	callbackFunc, exists = callbackFuncMap[id]

	return
}

// 删除回调方法
// id:自增Id
func deleteCallbackFunc(id int32) {
	requestMutex.Lock()
	defer requestMutex.Unlock()

	delete(callbackFuncMap, id)
}

// 向服务端发送请求
// methodName：传输类型
// parameters：调用的方法参数
// function：请求对应的回调方法
func request(methodName string, parameters []interface{}, function func(interface{})) {
	if methodName != Con_ChatServerCenter_Method_Login && clientObj.isLogin == false {
		msg := "客户端尚未登陆，请先登陆"
		logUtil.ErrorLog(msg)
		debugUtil.Println(msg)
		return
	}

	getIncrementId := func() int32 {
		atomic.AddInt32(&requestId, 1)
		return requestId
	}

	// 注册回调方法
	id := getIncrementId()
	requestObj := NewCenterRequestObject(methodName, parameters)
	message, _ := json.Marshal(requestObj)
	if err := clientObj.sendByteMessage(id, message); err != nil {
		clientObj.quit()
	}

	registerCallbackFunc(id, function)
	debugUtil.Printf("request id is:%d, methodName:%v\n", id, methodName)
}
