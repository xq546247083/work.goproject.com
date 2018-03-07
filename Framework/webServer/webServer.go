package webServer

import (
	"fmt"
	"net/http"
	"time"

	"work.goproject.com/Framework/ipMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

var (
	defaultPageMap = make(map[string]bool)
)

func init() {
	defaultPageMap["/"] = true
	defaultPageMap["/favicon.ico"] = true
}

// web服务对象
type WebServer struct {
	// 所有对外提供的处理器集合
	handlerMap map[string]*handler

	// 需要应用的Header集合
	headerMap map[string]string

	// 指定的HTTP请求方法，如果不指定则任意方法均有效
	specifiedMethod string

	// 当HTTP请求方法无效时，调用的处理器
	invalidMethodHandler func(*Context)

	// 默认页的处理器（默认页指的是/或/favicon.ico，通常是被负载均衡器调用或在浏览器中直接访问，而不会被程序调用）
	defaultPageHandler func(*Context)

	// 未找到指定的回调时的处理器，通常是因为该url是可变的，而不是固定的
	notFoundPageHandler func(*Context)

	// 是否需要验证IP
	isCheckIP bool

	// 在Debug模式下是否需要验证IP地址，默认情况下不验证
	ifCheckIPWhenDebug bool

	// 当IP无效时调用的处理器
	ipInvalidHandler func(*Context)

	// 当检测到参数无效时调用的处理器
	paramInvalidHandler func(*Context)

	// 处理请求数据的处理器（例如压缩、解密等）
	requestDataHandler func(*Context, []byte) ([]byte, error)

	// 处理响应数据的处理器（例如压缩、加密等）
	responseDataHandler func(*Context, []byte) ([]byte, error)

	// 请求执行时间的处理器
	executeTimeHandler func(*Context)
}

// 注册API
// path：注册的访问路径
// callback：回调方法
// configObj：Handler配置对象
func (this *WebServer) RegisterHandler(path string, handlerFuncObj handlerFunc, configObj *HandlerConfig) {
	// 判断是否已经注册过，避免命名重复
	if _, exists := this.handlerMap[path]; exists {
		panic(fmt.Sprintf("%s has been registered, please try a new path", path))
	}

	this.handlerMap[path] = newHandler(path, handlerFuncObj, configObj)
}

// 获取请求方法
// path:方法名称
// 返回值:
// 请求方法
// 是否存在
func (this *WebServer) getHandler(path string) (handlerObj *handler, exists bool) {
	handlerObj, exists = this.handlerMap[path]
	return
}

// 设定Http Header信息
func (this *WebServer) SetHeader(header map[string]string) {
	this.headerMap = header
}

// 处理Http Header信息
func (this *WebServer) handleHeader(context *Context) {
	if this.headerMap != nil && len(this.headerMap) > 0 {
		for k, v := range this.headerMap {
			context.responseWriter.Header().Set(k, v)
		}
	}
}

// 设定HTTP请求方法
func (this *WebServer) SetMethod(method string) {
	this.specifiedMethod = method
}

// 设定当HTTP请求方法无效时，调用的处理器
func (this *WebServer) SetInvalidMethodHandler(handler func(*Context)) {
	this.invalidMethodHandler = handler
}

// 处理HTTP请求方法
// 返回值
// isTerminate:是否终止本次请求
func (this *WebServer) handleMethod(context *Context) (isTerminate bool) {
	if this.specifiedMethod != "" {
		if context.request.Method != this.specifiedMethod {
			if this.invalidMethodHandler != nil {
				this.invalidMethodHandler(context)
			} else {
				http.Error(context.responseWriter, fmt.Sprintf("Expected %s Method", this.specifiedMethod), 406)
			}
			isTerminate = true
		}
	}

	return
}

// 设定默认页的处理器
func (this *WebServer) SetDefaultPageHandler(handler func(*Context)) {
	this.defaultPageHandler = handler
}

// 处理默认页
// 返回值
// isTerminate:是否终止本次请求
func (this *WebServer) handleDefaultPage(context *Context) (isTerminate bool) {
	// 首页进行特别处理
	if _, exists := defaultPageMap[context.GetRequestPath()]; exists {
		if this.defaultPageHandler != nil {
			this.defaultPageHandler(context)
		} else {
			context.WriteString("Welcome to home page.")
		}
		isTerminate = true
	}

	return
}

// 设定未找到指定的回调时的处理器
func (this *WebServer) SetNotFoundPageHandler(handler func(*Context)) {
	this.notFoundPageHandler = handler
}

// 设定在Debug模式下是否需要验证IP地址
func (this *WebServer) SetIfCheckIPWhenDebug(value bool) {
	this.ifCheckIPWhenDebug = value
}

// 设定当IP无效时调用的处理器
func (this *WebServer) SetIPInvalidHandler(handler func(*Context)) {
	this.ipInvalidHandler = handler
}

// 验证IP
// 返回值
// isTerminate:是否终止本次请求
func (this *WebServer) handleIP(context *Context) (isTerminate bool) {
	if this.isCheckIP == false {
		return
	}

	if debugUtil.IsDebug() == true && this.ifCheckIPWhenDebug == false {
		return
	}

	if ipMgr.IsIpValid(context.GetRequestIP()) {
		return
	}

	if this.ipInvalidHandler != nil {
		this.ipInvalidHandler(context)
	} else {
		http.Error(context.responseWriter, "你的IP不允许访问，请联系管理员", 401)
	}

	isTerminate = true
	return
}

// 设定当检测到参数无效时调用的处理器
func (this *WebServer) SetParamInvalidHandler(handler func(*Context)) {
	this.paramInvalidHandler = handler
}

// 设定处理请求数据的处理器（例如压缩、解密等）
func (this *WebServer) SetRequestDataHandler(handler func(*Context, []byte) ([]byte, error)) {
	this.requestDataHandler = handler
}

// 设定处理响应数据的处理器（例如压缩、加密等）
func (this *WebServer) SetResponseDataHandler(handler func(*Context, []byte) ([]byte, error)) {
	this.responseDataHandler = handler
}

// 设定请求执行时间的处理器
func (this *WebServer) SetExecuteTimeHandler(handler func(*Context)) {
	this.executeTimeHandler = handler
}

// 处理请求执行时间逻辑
func (this *WebServer) handleExecuteTime(context *Context) {
	if this.executeTimeHandler != nil {
		this.executeTimeHandler(context)
	}
}

// http应答处理
// responseWriter:应答对象
// request:请求对象
func (this *WebServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	defer func() {
		if data := recover(); data != nil {
			logUtil.LogUnknownError(data)
		}
	}()

	// 构造上下文对象
	context := newContext(request, responseWriter, this.requestDataHandler, this.responseDataHandler)

	// 处理Http Header信息
	this.handleHeader(context)

	// 处理默认页
	if this.handleDefaultPage(context) {
		return
	}

	// 处理返回值信息
	defer func() {
		// 处理请求执行时间逻辑
		this.handleExecuteTime(context)
	}()

	// 验证IP（全局）
	if this.handleIP(context) {
		return
	}

	// 处理请求方法
	if this.handleMethod(context) {
		return
	}

	// 根据路径选择不同的处理方法
	var handlerObj *handler
	var exists bool
	if handlerObj, exists = this.getHandler(context.GetRequestPath()); !exists {
		if this.notFoundPageHandler != nil {
			this.notFoundPageHandler(context)
		} else {
			http.Error(context.responseWriter, "访问的页面不存在", 404)
		}
		return
	}

	// Check IP（局部）
	if handlerObj.checkIP(context, this.ifCheckIPWhenDebug) == false {
		if this.ipInvalidHandler != nil {
			this.ipInvalidHandler(context)
		} else {
			http.Error(context.responseWriter, "你的IP不允许访问，请联系管理员", 401)
		}
		return
	}

	// Check Param
	if handlerObj.checkParam(context, this.specifiedMethod) == false {
		if this.paramInvalidHandler != nil {
			this.paramInvalidHandler(context)
		} else {
			http.Error(context.responseWriter, "你的参数不正确，请检查", 500)
		}
		return
	}

	// Call Function
	handlerObj.funcObj(context)

	context.EndTime = time.Now()
}

func newWebServer(isCheckIP bool) *WebServer {
	return &WebServer{
		handlerMap: make(map[string]*handler),
		headerMap:  make(map[string]string),
		isCheckIP:  isCheckIP,
	}
}
