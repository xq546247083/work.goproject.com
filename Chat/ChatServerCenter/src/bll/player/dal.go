package player

import (
	"work.goproject.com/Chat/ChatServerCenter/src/dal"
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
	result := dal.GetDB().Save(playerObj)
	if err = result.Error; err != nil {
		dal.WriteLog("player.update", err)
		return
	}

	return
}
