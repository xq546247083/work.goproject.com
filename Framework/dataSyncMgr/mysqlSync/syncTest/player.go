package main

import (
	"sync"
)

type player struct {
	// 玩家id
	Id string `gorm:"column:Id;primary_key"`

	// 玩家名称
	Name string `gorm:"column:Name"`
}

func (this *player) resetName(name string) {
	this.Name = name
}

func (this *player) tableName() string {
	return "player"
}

func newPlayer(id, name string) *player {
	return &player{
		Id:   id,
		Name: name,
	}
}

type playerMgr struct {
	playerMap map[string]*player

	mutex sync.Mutex
}

func (this *playerMgr) insert(obj *player) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.playerMap[obj.Id] = obj
}

func (this *playerMgr) delete(obj *player) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.playerMap, obj.Id)
}

func (this *playerMgr) randomSelect() *player {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, obj := range this.playerMap {
		return obj
	}
	return nil
}

func newPlayerMgr() *playerMgr {
	return &playerMgr{
		playerMap: make(map[string]*player),
	}
}
