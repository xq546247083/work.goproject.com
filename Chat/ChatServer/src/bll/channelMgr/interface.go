package channelMgr

import (
	. "work.goproject.com/Chat/ChatServerModel/src"
)

// 聊天频道接口
type IChannel interface {
	// 频道名称
	Channel() string

	// 配置名称
	ConfigName() string

	// 初始化配置
	InitConfig() error

	// 设置配置对象
	SetConfig(*config)

	// 重新加载配置
	ReloadConfig() error

	// 初始化基类
	InitBaseChannel() error

	// 初始化消息历史管理器
	InitHistoryMgr() error

	// 初始化历史
	InitHistory() error

	// 是否开启
	IsOpen(*Player) bool

	// 是否跨服聊天
	IsCanCrossServer() bool

	// 是不说话太快
	IsSpeakFast(*Player) bool

	// 说话
	Speak(*Player)

	// 验证发送消息参数
	ValidateSendMessage(*Player, *ForwardChatMessage, *sendMessageConfig) *ResultStatus

	// 验证历史消息参数
	ValidateHistoryParam(*Player) *ResultStatus

	// 处理消息(敏感词汇、消息长度等)
	HandleMessage(*ForwardChatMessage)

	// 获取玩家列表
	GetPlayerList(*ForwardChatMessage) []*Player

	// 创建返回给客户端的数据对象
	NewServerResponseData(*ForwardChatMessage) *ServerResponseData

	// 获取历史记录信息
	GetHistoryInfo(*Player) interface{}

	// 获取指定的历史信息
	GetHistory(*Player, *getHistoryConfig) (historyList []*ServerResponseData)

	// 保存历史消息
	SaveHistory(*ForwardChatMessage)

	// 删除历史数据
	DeleteHistory(*Player, *deleteHistoryConfig)
}
