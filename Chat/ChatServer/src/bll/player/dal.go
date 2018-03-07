package player

import (
	"work.goproject.com/Chat/ChatServer/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func get(id string) (playerObj *Player, exists bool, err error) {
	playerObj = NewEmptyPlayer()

	result := dal.GetDB().Where("Id = ?", id).First(playerObj)
	if result.RecordNotFound() {
		return
	}

	if err = result.Error; err != nil {
		dal.WriteLog("player.get", err)
		return
	}

	exists = true

	// 初始化玩家其它信息
	playerObj.Init()

	return
}

func insert(playerObj *Player) (err error) {
	result := dal.GetDB().Create(playerObj)
	if err = result.Error; err != nil {
		dal.WriteLog("player.insert", err)
		return
	}

	return
}

func update(playerObj *Player) (err error) {
	// 避免并发更新
	playerObj.Mutex.Lock()
	defer playerObj.Mutex.Unlock()

	result := dal.GetDB().Save(playerObj)
	if err = result.Error; err != nil {
		dal.WriteLog("player.update", err)
		return
	}

	return
}
