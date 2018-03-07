package src

import (
	"sync"
)

// 服务器组玩家对象
type ServerGroupPlayer struct {
	// 服务器组Id
	serverGroupId int32

	// 玩家集合
	playerMap map[string]*Player

	// 锁对象
	mutex sync.RWMutex
}

func (this *ServerGroupPlayer) GetPlayerList() (playerList []*Player) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	for _, playerObj := range this.playerMap {
		playerList = append(playerList, playerObj)
	}

	return
}

func (this *ServerGroupPlayer) GetPropertyPlayerList(property string, value string) (playerList []*Player) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	for _, playerObj := range this.playerMap {
		if propertyValue, err := playerObj.GetProperty(property); err == nil && propertyValue == value {
			playerList = append(playerList, playerObj)
		}
	}

	return
}

func (this *ServerGroupPlayer) AddPlayer(playerObj *Player) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.playerMap[playerObj.Id] = playerObj
}

func (this *ServerGroupPlayer) DeletePlayer(playerObj *Player) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.playerMap, playerObj.Id)
}

func NewServerGroupPlayer(_serverGroupId int32) *ServerGroupPlayer {
	return &ServerGroupPlayer{
		serverGroupId: _serverGroupId,
		playerMap:     make(map[string]*Player, 1024),
	}
}
