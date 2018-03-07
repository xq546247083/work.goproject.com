package server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/ipMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
	"work.goproject.com/goutil/zlibUtil"
)

// 请求上下文对象
type Context struct {
	// 请求对象
	*http.Request

	// 应答写对象
	http.ResponseWriter

	// 请求数据
	requestBytes []byte

	// 数据是否已经解析数据
	ifDataParsed bool

	// Form的数据是否已经解析
	ifFormParsed bool
}

// 检查IP是否合法
func (this *Context) checkIP() *ResultStatus {
	if debugUtil.IsDebug() == false && ipMgr.IsIpValid(webUtil.GetRequestIP(this.Request)) == false {
		return InvalidIP
	}

	return Success
}

func (this *Context) GetFormValue(key string) (value string, exists bool) {
	defer func() {
		this.ifFormParsed = true
	}()

	if !this.ifFormParsed {
		this.Request.ParseForm()
	}

	values := this.Form[key]
	if values != nil && len(values) > 0 {
		value = values[0]
		exists = true
		return
	}

	return
}

// 转换内容
func (this *Context) parseContent() error {
	defer func() {
		this.Body.Close()
		this.ifDataParsed = true
	}()

	data, err := ioutil.ReadAll(this.Body)
	if err != nil {
		logUtil.ErrorLog("url:%s,读取数据出错，错误信息为：%s", this.RequestURI, err)
		return err
	}

	this.requestBytes = data

	return nil
}

// 获取请求字节数据
// 返回值:
// []byte:请求字节数组
// error:错误信息
func (this *Context) GetRequestBytes(isCompressed bool) (result []byte, exists bool, err error) {
	if this.ifDataParsed == false {
		this.parseContent()
	}

	data := this.requestBytes
	if data == nil || len(data) <= 0 {
		return
	} else {
		exists = true
	}

	if isCompressed {
		result, err = zlibUtil.Decompress(data)
		if err != nil {
			logUtil.ErrorLog("解压缩请求数据失败:%s", err)
			return
		}
	} else {
		result = data
	}

	return
}

// 获取请求字符串数据
// 返回值:
// 请求字符串数据
func (this *Context) GetRequestString(isCompressed bool) (result string, exists bool, err error) {
	var data []byte
	data, exists, err = this.GetRequestBytes(isCompressed)
	if err != nil || !exists {
		return
	}

	result = string(data)
	exists = true

	return
}

// 反序列化
// moduleName:模块名称
// obj:反序列化结果数据
// isCompressed:数据是否已经被压缩
// 返回值:
// 错误对象
func (this *Context) Unmarshal(moduleName string, obj interface{}, isCompressed bool) (exists bool, err error) {
	var data []byte
	data, exists, err = this.GetRequestBytes(isCompressed)
	if err != nil || !exists {
		return
	}

	// 反序列化
	if err = json.Unmarshal(data, &obj); err != nil {
		logUtil.ErrorLog("Module:%s, 反序列化%s出错，错误信息为：%s", moduleName, string(data), err)
		return
	}

	return
}

// 输出字符串
func (this *Context) WriteString(result string) {
	this.ResponseWriter.Write([]byte(result))
}

// 输出json数据
func (this *Context) WriteJson(result interface{}) {
	if bytes, err := json.Marshal(result); err == nil {
		this.ResponseWriter.Write(bytes)
	}
}

// 跳转到其它页面
func (this *Context) RedirectTo(url string) {
	this.ResponseWriter.Header().Set("Location", url)
	this.ResponseWriter.WriteHeader(301)
}

// 新建API上下文对象
// request:请求对象
// responseWriter:应答写对象
// 返回值:
// *Context:上下文
func newContext(request *http.Request, responseWriter http.ResponseWriter) *Context {
	return &Context{
		Request:        request,
		ResponseWriter: responseWriter,
	}
}
