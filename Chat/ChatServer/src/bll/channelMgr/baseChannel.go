package channelMgr

import (
	"work.goproject.com/Chat/ChatServer/src/bll/dbConfig"
	"work.goproject.com/Chat/ChatServer/src/bll/word"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/stringUtil"
)

type baseChannel struct {
	*config
	channel string
}

func newBaseChannel(config *config, channel string) *baseChannel {
	return &baseChannel{
		config:  config,
		channel: channel,
	}
}

func (this *baseChannel) UpdateConfig(configObj *config) {
	this.config = configObj
}

func (this *baseChannel) IsOpen(playerObj *Player) bool {
	return playerObj.Lv >= this.config.OpenLv
}

func (this *baseChannel) IsCanCrossServer() bool {
	return this.config.IsCrossServer
}

func (this *baseChannel) IsSpeakFast(playerObj *Player) bool {
	return playerObj.IsSpeakFast(this.channel, this.config.MessageInterval)
}

func (this *baseChannel) Speak(playerObj *Player) {
	playerObj.Speak(this.channel)
}

func (this *baseChannel) ValidateSendMessage(playerObj *Player, chatMessageObj *ForwardChatMessage, configObj *sendMessageConfig) *ResultStatus {
	if dbConfig.IsOmit(configObj.Message) {
		return Success
	}

	if word.IfContainsForbidWords(configObj.Message) {
		return ContainForbiddenWord
	}

	return Success
}

func (this *baseChannel) ValidateHistoryParam(playerObj *Player) *ResultStatus {
	return Success
}

func (this *baseChannel) HandleMessage(chatMessageObj *ForwardChatMessage) {
	if dbConfig.IsOmit(chatMessageObj.Message) {
		return
	}

	// 处理敏感词汇
	chatMessageObj.Message = word.HandleSensitiveWords(chatMessageObj.Message)

	// 处理消息长度
	chatMessageObj.Message = this.HandleMessageLength(chatMessageObj.Message)
}

func (this *baseChannel) HandleMessageLength(message string) string {
	if len(message) > this.config.MaxMessageLength {
		return stringUtil.Substring(message, 0, this.config.MaxMessageLength)
	}

	return message
}

func (this *baseChannel) DeleteHistory(*Player, *deleteHistoryConfig) {
	// 基类没有任何行为，只是为了实现接口方法
	panic("DeleteHistory should not be called.Just to implement interface")
}
