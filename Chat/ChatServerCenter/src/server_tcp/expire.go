package server_tcp

import (
	"time"

	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/logUtil"
)

func init() {
	go func() {
		// 处理goroutine数量
		goroutineName := "server_tcp.expire"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		for {
			// 先等待5分钟，以便服务器启动
			time.Sleep(5 * time.Minute)

			// 获取客户端连接列表
			clientList := getClientList()

			// 记录日志
			logUtil.DebugLog("当前客户端数量为：%d，准备清理过期的客户端", len(clientList))

			for _, clientObj := range clientList {
				if clientObj.expired() {
					logUtil.DebugLog("客户端超时被断开，对应的信息为：%s", clientObj.String())
					clientObj.quit()
				} else {
					logUtil.DebugLog("客户端的信息为：%s", clientObj.String())
				}
			}
		}
	}()
}
