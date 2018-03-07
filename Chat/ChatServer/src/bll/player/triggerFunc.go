package player

import (
	"fmt"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

var (
	// 玩家注册、反注册触发方法
	playerRegisterTriggerFuncMap   = make(map[string]func(*Player))
	playerUnregisterTriggerFuncMap = make(map[string]func(*Player))
)

// 增加玩家注册触发方法
func AddPlayerRegisterTriggerFunc(funcName string, f func(*Player)) {
	if _, exists := playerRegisterTriggerFuncMap[funcName]; exists {
		panic(fmt.Errorf("%s already exists, please choose another name", funcName))
	}
	playerRegisterTriggerFuncMap[funcName] = f

	logUtil.DebugLog("AddPlayerRegisterTriggerFunc funcName:%s，当前共有%d个注册", funcName, len(playerRegisterTriggerFuncMap))
}

// 增加玩家反注册触发方法
func AddPlayerUnregisterTriggerFunc(funcName string, f func(*Player)) {
	if _, exists := playerUnregisterTriggerFuncMap[funcName]; exists {
		panic(fmt.Errorf("%s already exists, please choose another name", funcName))
	}
	playerUnregisterTriggerFuncMap[funcName] = f

	logUtil.DebugLog("AddPlayerUnregisterTriggerFunc funcName:%s，当前共有%d个注册", funcName, len(playerUnregisterTriggerFuncMap))
}

// 触发玩家注册方法
// playerObj:玩家对象
func triggerPlayerRegisterFunc(playerObj *Player) {
	for name, f := range playerRegisterTriggerFuncMap {
		debugUtil.Printf("开始调用玩家注册方法%s\n", name)
		f(playerObj)
		debugUtil.Printf("调用玩家注册方法%s结束\n", name)
	}
}

// 触发玩家反注册方法
// playerObj:玩家对象
func triggerPlayerUnregisterFunc(playerObj *Player) {
	for name, f := range playerUnregisterTriggerFuncMap {
		debugUtil.Printf("开始调用玩家反注册方法%s\n", name)
		f(playerObj)
		debugUtil.Printf("调用玩家反注册方法%s结束\n", name)
	}
}
