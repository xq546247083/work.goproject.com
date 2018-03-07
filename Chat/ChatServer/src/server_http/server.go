package server_http

import (
	"encoding/json"
	"net/http"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/webUtil"
)

// http服务对象
type httpServer struct{}

// http应答处理
// response:应答对象
// request:请求对象
func (this *httpServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	context := newContext(request, response)
	responseObj := NewServerResponseObject()

	defer func() {
		if data := recover(); data != nil {
			logUtil.LogUnknownError(data)
			return
		}

		// 特殊路径进行特别处理
		if request.URL.Path == "/" || request.URL.Path == "/favicon.ico" {
			return
		}

		// 获取输入参数的字符串形式
		parameter := ""
		if len(request.Form) > 0 {
			parameter_byte, _ := json.Marshal(request.Form)
			parameter = string(parameter_byte)
		}

		// 记录日志
		if debugUtil.IsDebug() || responseObj.ResultStatus != Success{
			result, _ := json.Marshal(responseObj)
			logUtil.DebugLog("%s-->IP:%s;Request:%v;Response:%s;", request.URL.Path, webUtil.GetRequestIP(request), parameter, string(result))
		}
	}()

	// 特殊路径进行特别处理
	if request.URL.Path == "/" || request.URL.Path == "/favicon.ico" {
		context.WriteString("ok")
		return
	}

	// 验证IP
	if rs := context.checkIP(); rs != Success {
		context.WriteJson(responseObj.SetResultStatus(rs))
		return
	}

	var handler Handler
	var exists bool
	if handler, exists = getHandler(request.URL.Path); !exists {
		logUtil.ErrorLog("访问的页面不存在，RequestAddr: %s  request.URL.Path: %s", request.RemoteAddr, request.URL.Path)
		http.Error(response, "访问的页面不存在", 404)
		return
	}

	// 调用处理方法，并返回结果
	responseObj = handler(context)
	context.WriteJson(responseObj)
}
