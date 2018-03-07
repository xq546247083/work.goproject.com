package onlineLog

import (
	"time"
)

// 在线日志对象
type OnlineLog struct {
	// 在线时间
	OnlineTime time.Time `gorm:"column:OnlineTime;primary_key"`

	// 序号
	Sid int `gorm:"column:Sid;primary_key"`

	// 聊天服务器地址
	ServerAddress string `gorm:"column:ServerAddress"`

	// 客户端数量
	ClientCount int `gorm:"column:ClientCount"`

	// 玩家数量
	PlayerCount int `gorm:"column:PlayerCount"`

	// 所有服务器的总数量
	TotalCount int `gorm:"column:TotalCount"`
}

func (this *OnlineLog) TableName() string {
	return "log_online"
}

func (this *OnlineLog) SetSid(value int) {
	this.Sid = value
}

func (this *OnlineLog) SetTotalCount(value int) {
	this.TotalCount = value
}

func NewOnlineLog(serverAddress string, clientCount, playerCount int) *OnlineLog {
	return &OnlineLog{
		OnlineTime:    time.Now(),
		ServerAddress: serverAddress,
		ClientCount:   clientCount,
		PlayerCount:   playerCount,
	}
}
