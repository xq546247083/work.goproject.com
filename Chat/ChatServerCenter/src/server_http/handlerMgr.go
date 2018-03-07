package server_http

import (
	"fmt"

	. "work.goproject.com/Chat/ChatServerModel/src"
)

type Handler func(*Context) *CenterResponseObject

var (
	// 处理函数集合
	handlerMap map[string]Handler
)

func init() {
	handlerMap = make(map[string]Handler, 8)
}

// 详细的注册一个WebAPI处理函数
// pattern:路由地址
// handler:处理函数
// paramInfo:参数列表
func RegisterHandler(pattern string, handler Handler) {
	if _, exist := handlerMap[pattern]; exist {
		panic(fmt.Errorf("存在重复的webapi注册：%s", pattern))
	}

	// 添加处理对象
	handlerMap[pattern] = handler
}

// 获取处理函数
func getHandler(pattern string) (handler Handler, exists bool) {
	handler, exists = handlerMap[pattern]
	return
}
