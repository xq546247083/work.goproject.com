package webServer

import (
	"fmt"
	"net/http"
	"sync"

	"work.goproject.com/goutil/logUtil"
)

// Https服务器对象
type HttpsServer struct {
	addr         string
	certFileName string
	keyFileName  string
	IWebServer
	server http.Server
}

func (this *HttpsServer) SetAddr(addr string) {
	this.addr = addr
	this.server.Addr = addr
}

// 启动HttpsServer
func (this *HttpsServer) Start(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	// 开启监听
	msg := fmt.Sprintf("http server begins to listen on: %s...", this.addr)
	fmt.Println(msg)
	logUtil.InfoLog(msg)

	if err := this.server.ListenAndServeTLS(this.certFileName, this.keyFileName); err != nil {
		panic(fmt.Sprintf("https server ListenAndServeTLS Error:%v", err))
	}
}

// 创建新的HttpsServer
// isCheckIP:该属性已丢弃，可以任意赋值
func NewHttpsServer(addr, certFileName, keyFileName string, isCheckIP bool) *HttpsServer {
	webServerObj := newWebServer(isCheckIP)

	return &HttpsServer{
		addr:         addr,
		certFileName: certFileName,
		keyFileName:  keyFileName,
		IWebServer:   webServerObj,
		server: http.Server{
			Addr:    addr,
			Handler: webServerObj,
		},
	}
}
