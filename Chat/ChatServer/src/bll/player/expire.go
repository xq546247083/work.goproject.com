package player

import (
	"time"

	"work.goproject.com/Chat/ChatServer/src/clientMgr"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/logUtil"
)

func init() {
	go clearExpirePlayer()
}

func clearExpirePlayer() {
	// 处理goroutine数量
	goroutineName := "player.clearExpirePlayer"
	goroutineMgr.Monitor(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	for {
		// 时间是清理客户端过期时间的2倍，是为了让玩家在内存中可以多存在一段时间，
		// 以便减轻数据库的压力
		time.Sleep(10 * time.Minute)

		beforePlayerCount := getPlayerCount()
		expirePlayerCount := 0

		allPlayerList := GetAllPlayerList()
		for _, item := range allPlayerList {
			if item.ClientId == 0 {
				UnRegisterPlayer(item)
				expirePlayerCount += 1
			} else if _, exists := clientMgr.GetClient(item.ClientId); !exists {
				UnRegisterPlayer(item)
				expirePlayerCount += 1
			}
		}

		// 记录日志
		logUtil.DebugLog("清理前的玩家数量为：%d，本次清理不活跃的玩家数量为：%d", beforePlayerCount, expirePlayerCount)
	}
}
