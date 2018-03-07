package server_tcp

import (
	"time"

	"work.goproject.com/Chat/ChatServerCenter/src/bll/onlineLog"
	"work.goproject.com/Framework/goroutineMgr"
)

func init() {
	go func() {
		// 处理goroutine数量
		goroutineName := "rpcServer.online"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		for {
			// 因为刚开始时不存在过期，所以先暂停5分钟
			time.Sleep(5 * time.Minute)

			// 获取客户端连接列表
			clientList := getClientList()
			onlineLogList := make([]*onlineLog.OnlineLog, 0, 8)

			for _, clientObj := range clientList {
				// 如果有ChatServer连接，则记录在线日志
				if clientObj.chatServer != nil && clientObj.expired() == false {
					onlineLogList = append(onlineLogList, onlineLog.NewOnlineLog(clientObj.chatServer.chatServerAddress, clientObj.chatServer.clientCount, clientObj.chatServer.playerCount))
				}
			}

			onlineLog.Save(onlineLogList)
		}
	}()
}
