package channelMgr

import (
	"time"

	. "work.goproject.com/Chat/ChatServerModel/src"
)

// 分组历史对象
// 需要实现IGroupHistory接口
type CountryHistory struct {
	// Id
	Id int `gorm:"column:Id;primary_key"`

	// 分组消息的唯一标识
	Identifier string `gorm:"column:Identifier"`

	// 聊天渠道
	Channel string `gorm:"column:Channel"`

	// 聊天消息
	Message string `gorm:"column:Message"`

	// 语音信息
	Voice string `gorm:"column:Voice"`

	// 说话者
	FromPlayer   string `gorm:"column:FromPlayer"`
	FromPlayerId string `gorm:"column:FromPlayerId"`

	// 创建时间
	Crtime time.Time `gorm:"column:Crtime"`
}

func (this *CountryHistory) GetId() int {
	return this.Id
}

func (this *CountryHistory) GetIdentifier() string {
	return this.Identifier
}

func (this *CountryHistory) SetIdentifier(id int, identifier string) {
	this.Id = id
	this.Identifier = identifier
	this.Crtime = time.Now()
}

func (this *CountryHistory) SetChannel(value string) {
	this.Channel = value
}

func (this *CountryHistory) SetMessage(message, voice string) {
	this.Message = message
	this.Voice = voice
}

func (this *CountryHistory) SetFromPlayer(fromPlayer, fromPlayerId string) {
	this.FromPlayer = fromPlayer
	this.FromPlayerId = fromPlayerId
}

// 转化为ResponseData，以便返回给客户端
func (this *CountryHistory) ToServerResponseData() *ServerResponseData {
	return ConvertToServerResponseData(
		this.Id,
		this.Channel,
		this.Message,
		this.Voice,
		this.FromPlayer,
		this.FromPlayerId,
		"",
		"",
		this.Crtime.Unix(),
	)
}

func newCountryHistory() IGroupHistory {
	return &CountryHistory{}
}
