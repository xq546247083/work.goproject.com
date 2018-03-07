package printMgr

import (
	"fmt"
	"sync"

	. "work.goproject.com/Framework/startMgr"
)

var (
	funcMap = make(map[string]*FuncItem)
	mutex   sync.Mutex
)

// 注册方法(如果名称重复会panic)
// name:方法名称（唯一标识）
// moduleType:模块类型
// definition:方法定义
func Register(name string, moduleType ModuleType, definition func() error) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := funcMap[name]; exists {
		panic(fmt.Sprintf("%s已经存在，请重新命名", name))
	}

	funcMap[name] = NewFuncItem(name, moduleType, definition)
}

// 调用所有方法
// 返回值：
// 错误列表
func CallAll() (errList []error) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, item := range funcMap {
		// 调用方法
		if err := item.Call(); err != nil {
			errList = append(errList, err)
		}
	}

	return
}

// 按照模块类型进行调用
// moduleType:模块类型
// 返回值：
// errList:错误列表
func CallType(moduleType ModuleType) (errList []error) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, item := range funcMap {
		if item.ModuleType != moduleType {
			continue
		}

		// 调用方法
		if err := item.Call(); err != nil {
			errList = append(errList, err)
		}
	}

	return
}
