package server_http

import (
	"fmt"
	"net/http"
	"sync"

	"work.goproject.com/goutil/logUtil"
)

// 启动Web服务器
// wg:WaitGroup对象
// address:服务器地址
func Start(wg *sync.WaitGroup, address string) {
	defer func() {
		wg.Done()
	}()

	// 开启服务
	serverInstance := http.Server{
		Addr:    address,
		Handler: new(httpServer),
	}

	// 开启监听
	msg := fmt.Sprintf("server_http begins to listen on: %s...", address)
	fmt.Println(msg)
	logUtil.InfoLog(msg)

	if err := serverInstance.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("server_http ListenAndServe Error:%v", err))
	}
}
