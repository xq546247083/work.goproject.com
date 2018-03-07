package src

import (
	"time"
)

// 在ChatServer和ChatServerCenter之间传输的玩家对象
type ForwardPlayer struct {
	// 玩家Id
	Id string

	// 禁言的结束时间
	SilentEndTime time.Time
}

func NewForwardPlayer(id string, silentEndTime time.Time) *ForwardPlayer {
	return &ForwardPlayer{
		Id:            id,
		SilentEndTime: silentEndTime,
	}
}
