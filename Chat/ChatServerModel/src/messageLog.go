package src

import (
	"time"
)

type MessageLog struct {
	// 主键
	Id int `gorm:"column:Id;primary_key"`

	// 玩家Id
	PlayerId string `gorm:"column:PlayerId"`

	// 玩家名称
	PlayerName string `gorm:"column:PlayerName"`

	// 合作商Id
	PartnerId int32 `gorm:"column:PartnerId"`

	// 服务器Id
	ServerId int32 `gorm:"column:ServerId"`

	// 服务器组Id
	ServerGroupId int32 `gorm:"column:ServerGroupId"`

	// 聊天内容
	Message string `gorm:"column:Message"`

	// 语音信息
	Voice string `gorm:"column:Voice"`

	// 聊天频道
	Channel string `gorm:"column:Channel"`

	// 如果是私聊，表示目标玩家
	ToPlayerId string `gorm:"column:ToPlayerId"`

	// 发送时间
	Crtime time.Time `gorm:"column:Crtime"`
}

func (this *MessageLog) TableName() string {
	return "log_message"
}

func NewMessageLog(playerId, playerName string, partnerId, serverId, serverGroupId int32, message, voice, channel, toPlayerId string) *MessageLog {
	return &MessageLog{
		PlayerId:      playerId,
		PlayerName:    playerName,
		PartnerId:     partnerId,
		ServerId:      serverId,
		ServerGroupId: serverGroupId,
		Message:       message,
		Voice:         voice,
		Channel:       channel,
		ToPlayerId:    toPlayerId,
		Crtime:        time.Now(),
	}
}

type MaxIdAndMinId struct {
	MaxId int64 `gorm:"column:MaxId"`
	MinId int64 `gorm:"column:MinId"`
}
