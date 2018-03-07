/*
	此包提供通用的HTTP/HTTPS服务器功能；
	使用方法如下：
	1、初始化一个HttpServer/HttpsServer
	server := NewHttpServer(addr string, isCheckIP bool)
	或
	server := NewHttpsServer(addr, certFileName, keyFileName string, isCheckIP bool)
	其中参数说明如下：
	addr:服务器监听地址
	isCheckIP:是否需要验证客户端IP（此设置针对所有的请求，对于每个请求也可以单独设置；以此设置为优先）
	certFileName:证书文件的路径
	keyFileName:密钥文件的路径

	2、设置server的属性
	在服务器运行过程中，会有一些默认的行为，如果想要改变默认行为，可以通过调用以下的方法进行设置

	1）、设定Http Header信息
	SetHeader(header map[string]string)
	默认的Http Header为空；如果设定了Http Header，则在每次请求时，都会给ResponseWriter加上该header属性
	2）、设定HTTP请求方法
	SetMethod(method string)
	默认情况下，对请求的方法没有任何限制；如果设定了Method，则只允许该Method访问
	3）、设定当HTTP请求方法无效时，调用的处理器
	SetInvalidMethodHandler(handler func(*Context))
	默认情况下，如果请求方法与所设置的方法不一致时，会返回406错误；如果设置了此属性，则会调用此属性进行处理
	4）、设定默认页的处理器
	SetDefaultPageHandler(handler func(*Context))
	默认页指的是/, /favicon.ico这两个页面；默认情况下，仅仅是输出Welcome to home page.；如果需要针对做一些处理，可以设置此属性
	5）、设定未找到指定的回调时的处理器
	SetNotFoundPageHandler(handler func(*Context))
	当服务器找不到对应的地址的Handler时，会返回404错误。如果需要处理一些非固定的地址时，可以使用此方法；比如回调地址中包含AppId，所以导致地址可变
	6）、设定在Debug模式下是否需要验证IP地址
	SetIfCheckIPWhenDebug(value bool)
	默认情况下，在DEBUG模式时不验证IP；可以通过此属性改变此此行为
	7）、设定当IP无效时调用的处理器
	SetIPInvalidHandler(handler func(*Context))
	当需要验证IP并且IP无效时，默认情况下会返回401错误；如果设定了此属性，则可以改变该行为
	8）、设定当检测到参数无效时调用的处理器
	SetParamInvalidHandler(handler func(*Context))
	当检测到参数无效时，默认情况下会返回500错误；如果设置了此属性，则可以改变该行为
	9）、设定处理请求数据的处理器（例如压缩、解密等）
	SetRequestDataHandler(handler func(*Context, []byte) ([]byte, error))
	如果设定此属性，则在处理接收到的请求数据时，会调用此属性
	10）、设定处理响应数据的处理器（例如压缩、加密等）
	SetResponseDataHandler(handler func(*Context, []byte) ([]byte, error))
	如果设定此属性，则在处理返回给客户端的数据时，会调用此属性
	11）、设定请求执行时间的处理器
	SetExecuteTimeHandler(handler func(*Context))
	如果在请求结束后想要处理调用的时间，则需要设置此属性；例如请求时间过长则记录日志等

	3、注册handler
	server.RegisterHandler(path string, handlerFuncObj handlerFunc, configObj *HandlerConfig)
	参数如下：
	// path：注册的访问路径
	// callback：回调方法
	// configObj：Handler配置对象
	例如：server.RegisterHandler("/get/notice", getNoticeConfig, &webServer.HandlerConfig{IsCheckIP: false, ParamNameList: []string{"appid"}})

	4、启动对应的服务器
	server.Start(wg *sync.WaitGroup)

	5、context中提供了很多实用的方法
	1）、GetRequestPath() string：获取请求路径（该路径不带参数）
	2）、GetRequestIP() string：获取请求的客户端的IP地址
	3）、GetExecuteSeconds() int64：获取请求执行的秒数。当然也可以通过获取StartTime, EndTime属性自己进行更高精度的处理
	4）、String() string：将context里面的内容进行格式化，主要用于记录日志
	5）、FormValue(key string) string：获取请求的参数值（包括GET/POST/PUT/DELETE等所有参数）
	6）、PostFormValue(key string) string：获取POST的参数值
	7）、GetFormValueData() typeUtil.MapData：获取所有参数的MapData类型（包括GET/POST/PUT/DELETE等所有参数）
	8）、GetPostFormValueData() typeUtil.MapData：获取POST参数的MapData类型
	9）、GetMultipartFormValueData() typeUtil.MapData：获取MultipartForm的MapData类型
	10）、GetRequestBytes() (result []byte, exists bool, err error)：获取请求字节数据
	11）、GetRequestString() (result string, exists bool, err error)：获取请求字符串数据
	12）、Unmarshal(obj interface{}) (exists bool, err error)：反序列化为对象（JSON）
	13）、WriteString(result string)：输出字符串给客户端
	14）、WriteJson(result interface{})：输出json数据给客户端
	15）、RedirectTo(url string)：重定向到其它页面
	如果以上的方法不能满足需求，则可以调用以下的方法来获取原始的Request/ResponseWriter对象进行处理
	16）、GetRequest() *http.Request：获取请求对象
	17）、GetResponseWriter() http.ResponseWriter：获取响应对象
*/
package webServer
