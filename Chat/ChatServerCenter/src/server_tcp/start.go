package server_tcp

import (
	"fmt"
	"net"
	"sync"

	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/Framework/ipMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/netUtil"
)

// 启动服务器
// wg：WaitGroup
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

		// 验证ip是否有效
		ip := netUtil.GetRemoteIp(conn)
		if debugUtil.IsDebug() == false && ipMgr.IsIpValid(ip) == false {
			conn.Close()
			logUtil.ErrorLog("RemoteIP:%s Not Allowed for server_tcp", ip)
			continue
		}

		// 创建客户端对象
		clientObj := newClient(conn)
		clientObj.start()
		registerClient(clientObj)
	}
}
