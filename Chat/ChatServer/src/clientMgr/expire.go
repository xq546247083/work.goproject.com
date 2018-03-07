package clientMgr

import (
	"time"

	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/logUtil"
)

func init() {
	// 清理过期的客户端
	go clearExpiredClient()
}

// 清理过期的客户端
func clearExpiredClient() {
	// 处理goroutine数量
	goroutineName := "clientMgr.clearExpiredClient"
	goroutineMgr.Monitor(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	for {
		// 休眠指定的时间（单位：秒）(放在此处是因为程序刚启动时并没有过期的客户端，所以先不用占用资源；并且此时LogPath尚未设置，如果直接执行后面的代码会出现panic异常)
		time.Sleep(5 * time.Minute)

		// 获取过期的客户端列表
		expiredClientList := getExpiredClientList()
		expiredClientCount := len(expiredClientList)
		beforeClientCount := getClientCount()

		// 客户端断开
		if expiredClientCount > 0 {
			for _, item := range expiredClientList {
				item.Quit()
			}
		}

		// 记录日志
		logUtil.DebugLog("清理前的客户端数量为：%d，本次清理不活跃的客户端数量为：%d", beforeClientCount, expiredClientCount)
	}
}
