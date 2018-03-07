package src

// 在ChatServer和ChatServerCenter之间传输的聊天消息对象
type ForwardChatMessage struct {
	// 消息Id
	Id int

	// 聊天消息
	Message string

	// 语音信息
	Voice string

	// 服务器组Id
	ServerGroupId int32

	// 聊天频道
	Channel string

	// 源玩家对象
	Player *Player

	// 聊天频道对应的目标属性(例如UnionId, ShimenId, Nation等)
	TargetProperty interface{}

	// 排除指定的聊天服务器
	ExcludeChatServer string

	// 是否被传输过？如果已经被传输过，则不用再次被传输
	IsTransfered bool
}

func NewForwardChatMessage(message, voice string, serverGroupId int32, channel string, player *Player) *ForwardChatMessage {
	return &ForwardChatMessage{
		Message:       message,
		Voice:         voice,
		ServerGroupId: serverGroupId,
		Channel:       channel,
		Player:        player,
	}
}
