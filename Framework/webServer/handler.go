package webServer

import (
	"work.goproject.com/Framework/ipMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/typeUtil"
)

type handlerFunc func(*Context)

// 请求处理配置对象
type HandlerConfig struct {
	// 是否需要验证IP
	IsCheckIP bool

	// 方法参数名称集合
	ParamNameList []string
}

// 请求处理对象
type handler struct {
	// 注册的访问路径
	path string

	// 方法定义
	funcObj handlerFunc

	// 请求处理配置对象
	configObj *HandlerConfig
}

// 检查IP是否合法
func (this *handler) checkIP(context *Context, ifCheckIPWhenDebug bool) bool {
	if this.configObj == nil {
		return true
	}

	if this.configObj.IsCheckIP == false {
		return true
	}

	if debugUtil.IsDebug() == true && ifCheckIPWhenDebug == false {
		return true
	}

	if ipMgr.IsIpValid(context.GetRequestIP()) {
		return true
	}

	return false
}

// 检测参数
func (this *handler) checkParam(context *Context, methodName string) bool {
	if this.configObj == nil {
		return true
	}

	for _, name := range this.configObj.ParamNameList {
		var formValueData typeUtil.MapData
		if methodName == "POST" {
			formValueData = context.GetPostFormValueData()
		} else {
			formValueData = context.GetFormValueData()
		}

		if _, exists := formValueData[name]; !exists {
			return false
		}
	}

	return true
}

// 检测POST参数
func (this *handler) checkPostParam(context *Context) bool {
	if this.configObj == nil {
		return true
	}

	for _, name := range this.configObj.ParamNameList {
		if _, exists := context.GetPostFormValueData()[name]; !exists {
			return false
		}
	}

	return true
}

// 创建新的请求方法对象
func newHandler(path string, funcObj handlerFunc, configObj *HandlerConfig) *handler {
	return &handler{
		path:      path,
		funcObj:   funcObj,
		configObj: configObj,
	}
}
