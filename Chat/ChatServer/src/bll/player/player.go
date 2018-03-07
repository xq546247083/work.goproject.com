package player

import (
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
)

// 更新来自于GS的玩家信息
func UpdateFromGS(gsPlayerObj *Player) {
	debugUtil.Printf("Receive data from gs. %v\n", gsPlayerObj)

	playerObj, exists, err := GetPlayer(gsPlayerObj.Id, true)
	if err != nil {
		return
	}

	if !exists {
		playerObj = NewPlayerFromGS(gsPlayerObj)
		insert(playerObj)
		RegisterPlayer(playerObj)
	} else {
		if playerObj.IsInfoChanged(gsPlayerObj) == false {
			return
		}

		playerObj.UpdateInfoFromGS(gsPlayerObj)
		update(playerObj)
	}
}

// 更新玩家的禁言状态
// forwardPlayerObj：转发的玩家对象
func UpdatePlayerInfo(forwardPlayerObj *ForwardPlayer) {
	playerObj, exists, err := GetPlayer(forwardPlayerObj.Id, false)
	if err != nil || !exists {
		return
	}

	playerObj.SilentEndTime = forwardPlayerObj.SilentEndTime
	update(playerObj)
}
