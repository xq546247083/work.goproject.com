package reloadMgr

import (
	"fmt"

	"work.goproject.com/goutil/logUtil"
)

var (
	reloadFuncMap = make(map[string]func() error)
)

// RegisterReloadFunc ...注册Reload方法
// funcName:方法名称
// reloadFunc：reload方法
func RegisterReloadFunc(funcName string, reloadFunc func() error) {
	if _, exists := reloadFuncMap[funcName]; exists {
		panic(fmt.Sprintf("%s已经存在，请重新取名", funcName))
	}

	reloadFuncMap[funcName] = reloadFunc
	logUtil.NormalLog(fmt.Sprintf("RegisterReloadFunc funcName:%s，当前共有%d个注册", funcName, len(reloadFuncMap)), logUtil.Info)
}

// Reload ...重新加载
// 返回值：
// 错误列表
func Reload() (errList []error) {
	for funcName, reloadFunc := range reloadFuncMap {
		if err := reloadFunc(); err == nil {
			logUtil.NormalLog(fmt.Sprintf("Call ReloadFunc:%s Success.", funcName), logUtil.Info)
		} else {
			logUtil.NormalLog(fmt.Sprintf("Call ReloadFunc:%s Fail, Error:%s", funcName, err), logUtil.Error)
			errList = append(errList, err)
		}
	}

	return
}
