package player

import (
	"time"

	"work.goproject.com/Chat/ChatServerCenter/src/server_tcp"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

// 根据Id获取玩家对象
// id：玩家Id
// 返回值：
// 玩家对象
// 是否存在该玩家
// 错误对象
func getPlayer(id string) (playerObj *Player, exists bool, err error) {
	if id == "" {
		return
	}

	playerObj, exists, err = get(id)
	return
}

// 更新玩家的禁言信息
// playerObj：玩家对象
// silentEndTime：禁言结束时间
func updateSilentInfo(playerObj *Player, silentEndTime time.Time) (err error) {
	playerObj.SilentEndTime = silentEndTime

	err = update(playerObj)
	if err != nil {
		return
	}

	// 推送给ChatServer
	server_tcp.ForwardObjectChannel <- NewForwardObject_Player(NewForwardPlayer(playerObj.Id, playerObj.SilentEndTime))
	return
}
