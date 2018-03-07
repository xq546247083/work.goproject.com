package src

// ChatServerCenter转发给ChatServer的对象
type ForwardObject struct {
	// 转发的消息类型
	ForwardType ForwardType

	// 转发的消息
	ForwardChatMessage *ForwardChatMessage

	// 转发的玩家对象
	ForwardPlayer *ForwardPlayer
}

func NewForwardObject_Player(forwardPlayerObj *ForwardPlayer) *ForwardObject {
	return &ForwardObject{
		ForwardType:   Con_UpdatePlayerInfo,
		ForwardPlayer: forwardPlayerObj,
	}
}

func NewForwardObject_ChatMessage(forwardChatMessageObj *ForwardChatMessage) *ForwardObject {
	return &ForwardObject{
		ForwardType:        Con_ChatMessage,
		ForwardChatMessage: forwardChatMessageObj,
	}
}

func NewForwardObject_Reload() *ForwardObject {
	return &ForwardObject{
		ForwardType: Con_Reload,
	}
}
