package src

import (
	"time"
)

// 服务器响应对象
type ServerResponseData struct {
	// Id
	Id int

	// 聊天频道
	Channel string

	// 聊天消息
	Message string

	// 语音信息
	Voice string

	// 发送人
	FromPlayer   string
	FromPlayerId string

	// 接收人
	ToPlayer   string `json:"ToPlayer,omitempty"`
	ToPlayerId string `json:"ToPlayerId,omitempty"`

	// 创建的时间戳
	TimeStamp int64
}

// 创建新的服务器响应对象
func NewServerResponseData(id int, channel, message, voice string, from, to *Player) *ServerResponseData {
	var fromPlayer, fromPlayerId, toPlayer, toPlayerId string
	if from != nil {
		fromPlayerId = from.Id
		fromPlayer = from.String()
	}
	if to != nil {
		toPlayerId = to.Id
		toPlayer = to.String()
	}

	return &ServerResponseData{
		Id:           id,
		Channel:      channel,
		Message:      message,
		Voice:        voice,
		FromPlayer:   fromPlayer,
		FromPlayerId: fromPlayerId,
		ToPlayer:     toPlayer,
		ToPlayerId:   toPlayerId,
		TimeStamp:    time.Now().Unix(),
	}
}

// 从其它类型转化为服务器响应对象
func ConvertToServerResponseData(id int, channel, message, voice, fromPlayer, fromPlayerId, toPlayer, toPlayerId string, timeStamp int64) *ServerResponseData {
	return &ServerResponseData{
		Id:           id,
		Channel:      channel,
		Message:      message,
		Voice:        voice,
		FromPlayer:   fromPlayer,
		FromPlayerId: fromPlayerId,
		ToPlayer:     toPlayer,
		ToPlayerId:   toPlayerId,
		TimeStamp:    timeStamp,
	}
}
