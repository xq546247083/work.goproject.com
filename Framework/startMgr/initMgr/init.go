package initMgr

import (
	"fmt"
	"sync"
	"time"

	. "work.goproject.com/Framework/startMgr"
	"work.goproject.com/goutil/logUtil"
)

var (
	funcMap     = make(map[string]*FuncItem)
	mutex       sync.Mutex
	operateName = "Init"
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

// 调用一个方法
// name:方法名称（唯一标识）
// 返回值：
// 错误对象
func CallOne(name string) (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if item, exists := funcMap[name]; !exists {
		panic(fmt.Sprintf("%s不存在", name))
	} else {
		// 调用方法
		err = item.Call2(operateName)
	}

	return
}

// 调用任意数量的方法
// nameList:任意数量的方法名称
// 返回值：
// 错误列表
func CallAny(nameList ...string) (errList []error) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, name := range nameList {
		if item, exists := funcMap[name]; !exists {
			panic(fmt.Sprintf("%s不存在", name))
		} else {
			// 调用方法
			if err := item.Call2(operateName); err != nil {
				errList = append(errList, err)
			}
		}
	}

	return
}

// 调用所有方法
// 返回值：
// 错误列表
func CallAll() (errList []error) {
	mutex.Lock()
	defer mutex.Unlock()

	startTime := time.Now()
	defer func() {
		logUtil.InfoLog("%s 执行总时间:%s", operateName, time.Since(startTime))
	}()

	for _, item := range funcMap {
		// 调用方法
		if err := item.Call2(operateName); err != nil {
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

	startTime := time.Now()
	defer func() {
		logUtil.InfoLog("%v %s 执行总时间:%s", moduleType, operateName, time.Since(startTime))
	}()

	for _, item := range funcMap {

		if item.ModuleType != moduleType {
			continue
		}

		// 调用方法
		if err := item.Call2(operateName); err != nil {
			errList = append(errList, err)
		}
	}

	return
}
