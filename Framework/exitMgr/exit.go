package exitMgr

import (
	"fmt"

	"work.goproject.com/goutil/logUtil"
)

var (
	exitFuncMap = make(map[string]func())
)

// RegisterExitFunc ...注册Exit方法
// funcName:方法名称
// exitFunc：exit方法
func RegisterExitFunc(funcName string, exitFunc func()) {
	if _, exists := exitFuncMap[funcName]; exists {
		panic(fmt.Sprintf("%s已经存在，请重新取名", funcName))
	}

	exitFuncMap[funcName] = exitFunc
	logUtil.NormalLog(fmt.Sprintf("RegisterExitFunc funcName:%s，当前共有%d个注册", funcName, len(exitFuncMap)), logUtil.Info)
}

// Exit ...退出程序
// 返回值：
// 无
func Exit() {
	for funcName, exitFunc := range exitFuncMap {
		exitFunc()
		logUtil.NormalLog(fmt.Sprintf("Call ExitFunc:%s Finish.", funcName), logUtil.Info)
	}
}
