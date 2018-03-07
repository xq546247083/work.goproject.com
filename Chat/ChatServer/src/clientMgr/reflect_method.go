package clientMgr

import (
	"reflect"
)

// 反射的方法和输入、输出参数类型组合类型
type methodAndInOutTypes struct {
	// 反射出来的对应方法对象
	Method reflect.Value

	// 反射出来的方法的输入参数的类型集合
	InTypes []reflect.Type

	// 反射出来的方法的输出参数的类型集合
	OutTypes []reflect.Type
}

func newmethodAndInOutTypes(_method reflect.Value, _inTypes []reflect.Type, _outTypes []reflect.Type) *methodAndInOutTypes {
	return &methodAndInOutTypes{
		Method:   _method,
		InTypes:  _inTypes,
		OutTypes: _outTypes,
	}
}
