package server_tcp

import (
	"fmt"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/logUtil"
)

// 处理器对象
type handler struct {
	// 方法名称
	name string

	// 方法定义
	handlerFunc func(*client, []interface{}) *CenterResponseObject

	// 方法参数数量
	paramCount int
}

// 检测参数数量
// parameters：参数数组
func (this *handler) checkParamCount(parameters []interface{}) *ResultStatus {
	if this.paramCount > 0 && len(parameters) == 0 {
		return ParamIsEmpty
	}

	if this.paramCount != len(parameters) {
		return ParamNotMatch
	}

	return Success
}

// 创建新的请求方法对象
// _name：方法名称
// _handlerFunc：方法定义
// _paramCount：方法参数数量
func newHandler(_name string, _handlerFunc func(*client, []interface{}) *CenterResponseObject, _paramCount int) *handler {
	return &handler{
		name:        _name,
		handlerFunc: _handlerFunc,
		paramCount:  _paramCount,
	}
}

var (
	// 所有对外提供的处理器集合
	handlerMap = make(map[string]*handler)
)

// 注册方法
// name:方法名称
// handlerFunc:方法定义
// paramCount:方法参数数量
func registerHandler(name string, handlerFunc func(*client, []interface{}) *CenterResponseObject, paramCount int) {
	if _, exists := handlerMap[name]; exists {
		panic(fmt.Errorf("存在重复的注册：%s", name))
	}

	handlerMap[name] = newHandler(name, handlerFunc, paramCount)
	logUtil.DebugLog("server_tcp.RegisterHandler:%s", name)
}

func getHandler(name string) (handlerObj *handler, exists bool) {
	handlerObj, exists = handlerMap[name]
	return
}
