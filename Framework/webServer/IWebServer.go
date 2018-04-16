package webServer

import (
	"net/http"

	"work.goproject.com/goutil/typeUtil"
)

// webserver接口
type IWebServer interface {
	// 注册API
	// path：注册的访问路径
	// callback：回调方法
	// configObj：Handler配置对象
	RegisterHandler(path string, handlerFuncObj handlerFunc, configObj *HandlerConfig)

	// 设定Http Header信息
	SetHeader(header map[string]string)

	// 设定HTTP请求方法
	SetMethod(method string)

	// 设定当HTTP请求方法无效时，调用的处理器
	SetInvalidMethodHandler(handler func(*Context))

	// 设定默认页的处理器
	SetDefaultPageHandler(handler func(*Context))

	// 设定未找到指定的回调时的处理器
	SetNotFoundPageHandler(handler func(*Context))

	// 设定在Debug模式下是否需要验证IP地址
	SetIfCheckIPWhenDebug(value bool)

	// 设定当IP无效时调用的处理器
	SetIPInvalidHandler(handler func(*Context))

	// 设定用于检测参数是否有效的处理器
	SetParamCheckHandler(handler func(typeUtil.MapData, []string) ([]string, bool))

	// 设定当检测到参数无效时调用的处理器
	SetParamInvalidHandler(handler func(*Context, []string))

	// 设定处理请求数据的处理器（例如压缩、解密等）
	SetRequestDataHandler(handler func(*Context, []byte) ([]byte, error))

	// 设定处理响应数据的处理器（例如压缩、加密等）
	SetResponseDataHandler(handler func(*Context, []byte) ([]byte, error))

	// 设定请求执行时间的处理器
	SetExecuteTimeHandler(handler func(*Context))

	// http应答处理
	// responseWriter:应答对象
	// request:请求对象
	ServeHTTP(responseWriter http.ResponseWriter, request *http.Request)
}
