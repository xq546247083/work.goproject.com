package webServer

import (
	"work.goproject.com/Framework/ipMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/typeUtil"
)

// 参数传输类型
type ParamTransferType string

const (
	// 以Form表单形式来传递参数
	Con_Form ParamTransferType = "Form"
	// 以二进制流形式来传递参数
	Con_Stream ParamTransferType = "Stream"
)

type handlerFunc func(*Context)

// 请求处理配置对象
type HandlerConfig struct {
	// 是否需要验证IP
	IsCheckIP bool

	// 参数类型
	ParamTransferType ParamTransferType

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
func (this *handler) checkParam(context *Context, methodName string, paramCheckHandler func(paramMap typeUtil.MapData, paramNameList []string) ([]string, bool)) (missParamList []string, valid bool) {
	valid = true

	if this.configObj == nil {
		return
	}

	// 获取方法的参数集合
	var dataMap typeUtil.MapData
	switch this.configObj.ParamTransferType {
	case Con_Stream:
		data := new(map[string]interface{})
		if exist, err := context.Unmarshal(data); err != nil || !exist {
			valid = false
			return
		}
		dataMap = typeUtil.NewMapData(*data)
	default:
		if methodName == "POST" {
			dataMap = context.GetPostFormValueData()
		} else {
			dataMap = context.GetFormValueData()
		}
	}

	// 定义默认的参数验证方法
	defaultParamCheckHandler := func(paramMap typeUtil.MapData, paramNameList []string) (_missParamList []string, _valid bool) {
		_valid = true

		// 遍历判断每一个参数是否存在；为了搜集所有的参数，所以不会提前返回
		for _, name := range paramNameList {
			if _, exist := paramMap[name]; !exist {
				_missParamList = append(_missParamList, name)
				_valid = false
			}
		}

		return
	}

	// 如果没有指定验证参数的方法，就使用默认方法
	if paramCheckHandler == nil {
		missParamList, valid = defaultParamCheckHandler(dataMap, this.configObj.ParamNameList)
	} else {
		missParamList, valid = paramCheckHandler(dataMap, this.configObj.ParamNameList)
	}

	return
}

// 检测POST参数
func (this *handler) checkPostParam(context *Context) bool {
	if this.configObj == nil {
		return true
	}

	for _, name := range this.configObj.ParamNameList {
		if _, exist := context.GetPostFormValueData()[name]; !exist {
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
