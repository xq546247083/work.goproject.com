package webServer

import (
	"fmt"
	"net/http"
	"sync"

	"work.goproject.com/goutil/logUtil"
)

// Http服务器对象
type HttpServer struct {
	addr string
	IWebServer
	server http.Server
}

func (this *HttpServer) SetAddr(addr string) {
	this.addr = addr
	this.server.Addr = addr
}

// 启动HttpServer
func (this *HttpServer) Start(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	// 开启监听
	msg := fmt.Sprintf("http server begins to listen on: %s...", this.addr)
	fmt.Println(msg)
	logUtil.InfoLog(msg)

	if err := this.server.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("http server ListenAndServe Error:%v", err))
	}
}

// 创建新的HttpServer
// isCheckIP:该属性已丢弃，可以任意赋值
func NewHttpServer(addr string, isCheckIP bool) *HttpServer {
	webServerObj := newWebServer(isCheckIP)

	return &HttpServer{
		addr:       addr,
		IWebServer: webServerObj,
		server: http.Server{
			Addr:    addr,
			Handler: webServerObj,
		},
	}
}
