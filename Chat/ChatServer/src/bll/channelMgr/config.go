package channelMgr

// 聊天频道的配置信息
type config struct {
	// 开放等级
	OpenLv int

	// 是否可以跨服(不针对跨服频道、世界频道；因为它们具有内在的是不可跨服的属性)
	IsCrossServer bool

	// 消息发送间隔（单位：秒）
	MessageInterval int

	// 最大消息长度
	MaxMessageLength int

	// 是否可以包含禁止词汇
	IfCanContainForbidContent bool

	// 最大历史消息数量(如果是私聊，则表示保存天数)
	MaxHistoryCount int

	// 对应的玩家属性名称
	PlayerProperty string
}

type sendMessageConfig struct {
	// 消息内容
	Message string

	// 语音
	Voice string

	// 私聊目标玩家Id
	ToPlayerId string
}

func NewSendMessageConfig(message, voice, toPlayerId string) *sendMessageConfig {
	return &sendMessageConfig{
		Message:    message,
		Voice:      voice,
		ToPlayerId: toPlayerId,
	}
}

// 获取历史的配置
type getHistoryConfig struct {
	// 消息Id
	MessageId int

	// 消息数量
	Count int

	// 目标玩家Id
	TargetPlayerId string
}

func NewGetHistoryConfig(messageId, count int, targetPlayerId string) *getHistoryConfig {
	return &getHistoryConfig{
		MessageId:      messageId,
		Count:          count,
		TargetPlayerId: targetPlayerId,
	}
}

type deleteHistoryConfig struct {
	// 目标玩家Id
	TargetPlayerId string
}

func NewDeleteHistoryConfig(targetPlayerId string) *deleteHistoryConfig {
	return &deleteHistoryConfig{
		TargetPlayerId: targetPlayerId,
	}
}
