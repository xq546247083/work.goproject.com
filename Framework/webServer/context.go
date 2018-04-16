package webServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/netUtil"
	"work.goproject.com/goutil/typeUtil"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// 请求上下文对象
type Context struct {
	// 请求对象
	request *http.Request

	// 响应对象
	responseWriter http.ResponseWriter

	// 数据是否已经解析数据
	ifBodyParsed bool

	// 请求数据
	bodyContent []byte

	// Form的数据是否已经解析
	ifFormParsed bool

	// MultipleForm的数据是否已经解析
	ifMultipartFormParsed bool

	// 请求的开始时间
	StartTime time.Time

	// 请求的结束时间
	EndTime time.Time

	// 处理请求数据的处理器（例如压缩、解密等）
	requestDataHandler func(*Context, []byte) ([]byte, error)

	// 处理响应数据的处理器（例如压缩、加密等）
	responseDataHandler func(*Context, []byte) ([]byte, error)
}

// 获取请求对象
func (this *Context) GetRequest() *http.Request {
	return this.request
}

// 获取响应对象
func (this *Context) GetResponseWriter() http.ResponseWriter {
	return this.responseWriter
}

// 获取请求路径（不带参数）
func (this *Context) GetRequestPath() string {
	return this.request.URL.Path
}

// 获取请求的客户端的IP地址
func (this *Context) GetRequestIP() string {
	return netUtil.GetHttpAddr(this.request).Host
}

// 获取请求执行的秒数
func (this *Context) GetExecuteSeconds() int64 {
	return this.EndTime.Unix() - this.StartTime.Unix()
}

// 格式化context对象
func (this *Context) String() string {
	var bodyContent string
	if bytes, exist, err := this.GetRequestBytes(); err != nil && exist {
		bodyContent = string(bytes)
	}

	return fmt.Sprintf("IP:%s, URL:%s, FormValue:%v, BodyContent:%s",
		this.GetRequestIP(), this.GetRequestPath(), this.GetFormValueData(), bodyContent)
}

func (this *Context) parseForm() {
	if !this.ifFormParsed {
		this.request.ParseForm()
		this.ifFormParsed = true
	}
}

// 获取请求的参数值（包括GET/POST/PUT/DELETE等所有参数）
func (this *Context) FormValue(key string) (value string) {
	this.parseForm()

	return this.request.FormValue(key)
}

// 获取POST的参数值
func (this *Context) PostFormValue(key string) (value string) {
	this.parseForm()

	return this.request.PostFormValue(key)
}

// 获取所有参数的MapData类型（包括GET/POST/PUT/DELETE等所有参数）
func (this *Context) GetFormValueData() typeUtil.MapData {
	this.parseForm()

	valueMap := make(map[string]interface{})
	for k, v := range this.request.Form {
		valueMap[k] = v[0]
	}

	return typeUtil.MapData(valueMap)
}

// 获取POST参数的MapData类型
func (this *Context) GetPostFormValueData() typeUtil.MapData {
	this.parseForm()

	valueMap := make(map[string]interface{})
	for k, v := range this.request.PostForm {
		valueMap[k] = v[0]
	}

	return typeUtil.MapData(valueMap)
}

func (this *Context) parseMultipartForm() {
	if !this.ifMultipartFormParsed {
		this.request.ParseMultipartForm(defaultMaxMemory)
		this.ifMultipartFormParsed = true
	}
}

// 获取MultipartForm的MapData类型
func (this *Context) GetMultipartFormValueData() typeUtil.MapData {
	this.parseMultipartForm()

	valueMap := make(map[string]interface{})
	if this.request.MultipartForm != nil {
		for k, v := range this.request.MultipartForm.Value {
			valueMap[k] = v[0]
		}
	}

	return typeUtil.MapData(valueMap)
}

func (this *Context) parseBodyContent() (err error) {
	if this.ifBodyParsed {
		return
	}

	defer func() {
		this.request.Body.Close()
		this.ifBodyParsed = true
	}()

	this.bodyContent, err = ioutil.ReadAll(this.request.Body)
	if err != nil {
		logUtil.ErrorLog("url:%s,read body failed. Err：%s", this.GetRequestPath(), err)
		return
	}

	return
}

// 获取请求字节数据
// 返回值:
// []byte:请求字节数组
// exist:是否存在数据
// error:错误信息
func (this *Context) GetRequestBytes() (result []byte, exist bool, err error) {
	if err = this.parseBodyContent(); err != nil {
		return
	}

	result = this.bodyContent
	if result == nil || len(result) == 0 {
		return
	}

	// handle request data
	if this.requestDataHandler != nil {
		if result, err = this.requestDataHandler(this, result); err != nil {
			return
		}
	}

	exist = true

	return
}

// 获取请求字符串数据
// 返回值:
// result:请求字符串数据
// exist:是否存在数据
// error:错误信息
func (this *Context) GetRequestString() (result string, exist bool, err error) {
	var data []byte
	if data, exist, err = this.GetRequestBytes(); err != nil || !exist {
		return
	}

	result = string(data)
	exist = true

	return
}

// 反序列化为对象（JSON）
// obj:反序列化结果数据
// isCompressed:数据是否已经被压缩
// 返回值:
// 错误对象
func (this *Context) Unmarshal(obj interface{}) (exist bool, err error) {
	var data []byte
	if data, exist, err = this.GetRequestBytes(); err != nil || !exist {
		return
	}

	// Unmarshal
	if err = json.Unmarshal(data, &obj); err != nil {
		logUtil.ErrorLog("Unmarshal %s failed. Err:：%s", string(data), err)
		return
	}

	exist = true

	return
}

func (this *Context) writeBytes(data []byte) {
	if this.responseDataHandler != nil {
		var err error
		if data, err = this.responseDataHandler(this, data); err != nil {
			return
		}
	}

	this.responseWriter.Write(data)
}

// 输出字符串给客户端
func (this *Context) WriteString(result string) {
	this.writeBytes([]byte(result))
}

// 输出json数据给客户端
func (this *Context) WriteJson(result interface{}) {
	data, err := json.Marshal(result)
	if err != nil {
		logUtil.ErrorLog("Marshal %v failed. Err:：%s", result, err)
		return
	}

	this.writeBytes(data)
}

// 重定向到其它页面
func (this *Context) RedirectTo(url string) {
	this.responseWriter.Header().Set("Location", url)
	this.responseWriter.WriteHeader(301)
}

// 新建API上下文对象
// request:请求对象
// responseWriter:应答写对象
// 返回值:
// *Context:上下文
func newContext(request *http.Request, responseWriter http.ResponseWriter,
	requestDataHandler func(*Context, []byte) ([]byte, error),
	responseDataHandler func(*Context, []byte) ([]byte, error)) *Context {
	return &Context{
		request:             request,
		responseWriter:      responseWriter,
		StartTime:           time.Now(),
		EndTime:             time.Now(),
		requestDataHandler:  requestDataHandler,
		responseDataHandler: responseDataHandler,
	}
}
