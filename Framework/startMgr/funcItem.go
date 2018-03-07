package startMgr

import (
	"time"

	"work.goproject.com/goutil/runtimeUtil"

	"work.goproject.com/goutil/logUtil"
)

// 模块类型
// 具体有哪些模块类型需要对应的程序自己定义
type ModuleType int

// 方法项定义
type FuncItem struct {
	// 方法名称
	Name string

	// 模块类型
	ModuleType ModuleType

	// 方法定义
	Definition func() error
}

// 调用方法
// 返回值：
// 错误对象
func (this *FuncItem) Call() error {
	err := this.Definition()

	return err
}

// 调用方法
// operateName:操作名称
// 返回值：
// 错误对象
func (this *FuncItem) Call2(operateName string) error {
	startTime := time.Now()
	defer func() {
		logUtil.InfoLog("%s %s 执行耗时:%v 执行后内存总占用：%vMB", this.Name, operateName, time.Since(startTime).String(), runtimeUtil.GetMemSize()/1024.0/1024.0)
	}()

	err := this.Definition()

	return err
}

// 新建函数项
// name:模块名
// moduleType:对应的模块类型
// definition:目标处理函数
func NewFuncItem(name string, moduleType ModuleType, definition func() error) *FuncItem {
	return &FuncItem{
		Name:       name,
		ModuleType: moduleType,
		Definition: definition,
	}
}
