package server_tcp

import (
	"time"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/goroutineMgr"
)

var (
	// 转发对象的通道
	ForwardObjectChannel = make(chan *ForwardObject, 1024)
)

func init() {
	go func() {
		// 处理goroutine数量
		goroutineName := "rpcServer.forward"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		for {
			select {
			case forwardObj := <-ForwardObjectChannel:
				push(getClientList(), forwardObj)
			default:
				// 如果channel中没有数据，则休眠5毫秒
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
}
