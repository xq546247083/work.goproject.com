package server_tcp

import (
	"fmt"
	"net"
	"sync"

	"work.goproject.com/Chat/ChatServer/src/clientMgr"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/logUtil"
)

// 启动服务器
func Start(wg *sync.WaitGroup, address string) {
	defer func() {
		wg.Done()
	}()

	// 处理goroutine数量
	goroutineName := "server_tcp.Start"
	goroutineMgr.Monitor(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	msg := fmt.Sprintf("server_tcp begins to listen on:%s...", address)
	fmt.Println(msg)
	logUtil.InfoLog(msg)

	// 监听指定的端口
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("server_tcp listen Error: %s", err))
	}

	for {
		// 阻塞直至新连接到来
		conn, err := listener.Accept()
		if err != nil {
			logUtil.ErrorLog("server_tcp accept error: %s", err)
			continue
		}

		// 创建客户端对象
		clientObj := newClient(conn)
		clientObj.start()
		clientMgr.RegisterClient(clientObj)
	}
}
