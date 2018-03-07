package player

import (
	"sync"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
)

var (
	// 玩家集合
	playerMap   = make(map[string]*Player, 1024)
	playerMutex sync.RWMutex

	// 区服玩家集合
	serverGroupPlayerMap   = make(map[int32]*ServerGroupPlayer)
	serverGroupPlayerMutex sync.RWMutex
)

func GetServerGroupList() (list []*ServerGroupPlayer) {
	serverGroupPlayerMutex.RLock()
	defer serverGroupPlayerMutex.RUnlock()

	for _, item := range serverGroupPlayerMap {
		list = append(list, item)
	}

	return
}

// 获取服务器组对应的玩家列表对象
// serverGroupId:服务器组Id
// 返回值:
// 服务器组对应的玩家列表对象
// 是否存在
func GetServerGroupPlayer(serverGroupId int32) *ServerGroupPlayer {
	serverGroupPlayerMutex.Lock()
	defer serverGroupPlayerMutex.Unlock()

	if serverGroupPlayerObj, exists := serverGroupPlayerMap[serverGroupId]; !exists {
		serverGroupPlayerObj = NewServerGroupPlayer(serverGroupId)
		serverGroupPlayerMap[serverGroupId] = serverGroupPlayerObj
		return serverGroupPlayerObj
	} else {
		return serverGroupPlayerObj
	}
}

// 注册玩家对象到缓存中
// playerObj：玩家对象
func RegisterPlayer(playerObj *Player) {
	registerToPlayerMap := func(playerObj *Player) {
		playerMutex.Lock()
		defer playerMutex.Unlock()
		playerMap[playerObj.Id] = playerObj
	}

	registerToServerGroup := func(playerObj *Player) {
		serverGroupPlayerObj := GetServerGroupPlayer(playerObj.ServerGroupId())
		serverGroupPlayerObj.AddPlayer(playerObj)
	}

	// 添加到玩家集合中
	registerToPlayerMap(playerObj)

	// 添加到区服玩家集合中
	registerToServerGroup(playerObj)

	// 触发玩家注册方法
	triggerPlayerRegisterFunc(playerObj)
}

// 从缓存中取消玩家注册
// playerObj：玩家对象
func UnRegisterPlayer(playerObj *Player) {
	unRegisterToPlayerMap := func(playerObj *Player) {
		playerMutex.Lock()
		defer playerMutex.Unlock()
		delete(playerMap, playerObj.Id)
	}

	unRegisterToServerGroup := func(playerObj *Player) {
		serverGroupPlayerObj := GetServerGroupPlayer(playerObj.ServerGroupId())
		serverGroupPlayerObj.DeletePlayer(playerObj)
	}

	// 从玩家集合中删除
	unRegisterToPlayerMap(playerObj)

	// 从区服玩家集合中删除
	unRegisterToServerGroup(playerObj)

	// 触发玩家反注册方法
	triggerPlayerUnregisterFunc(playerObj)
}

// 根据Id获取玩家对象（先从缓存中取，取不到再从数据库中去取）
// id：玩家Id
// isLoadFromDB：是否要从数据库中获取数据
// 返回值：
// 玩家对象
// 是否存在该玩家
// 错误对象
func GetPlayer(id string, isLoadFromDB bool) (playerObj *Player, exists bool, err error) {
	if id == "" {
		return
	}

	getPlayerFromCache := func(_id string) (_playerObj *Player, _exists bool) {
		playerMutex.RLock()
		defer playerMutex.RUnlock()

		_playerObj, _exists = playerMap[id]
		return
	}

	getPlayerFromDB := func(_id string) (_playerObj *Player, _exists bool, _err error) {
		_playerObj, _exists, _err = get(id)
		return
	}

	if playerObj, exists = getPlayerFromCache(id); !exists && isLoadFromDB {
		playerObj, exists, err = getPlayerFromDB(id)
	}

	debugUtil.Printf("GetPlayer Id:%s, exists:%v, err:%v\n", id, exists, err)

	return
}

// 获取玩家数量
// 返回值：
// 玩家数量
func getPlayerCount() int {
	playerMutex.RLock()
	defer playerMutex.RUnlock()

	return len(playerMap)
}

// 获取所有的玩家列表
// 返回值：
// 所有的玩家列表
func GetAllPlayerList() (list []*Player) {
	playerMutex.RLock()
	defer playerMutex.RUnlock()

	for _, item := range playerMap {
		list = append(list, item)
	}

	return
}
