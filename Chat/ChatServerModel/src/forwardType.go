package src

// ChatServerCenter转发的消息类型
type ForwardType string

const (
	// 聊天消息
	Con_ChatMessage ForwardType = "ChatMessage"

	// 更新玩家信息
	Con_UpdatePlayerInfo ForwardType = "UpdatePlayerInfo"

	// 重新加载配置
	Con_Reload ForwardType = "Reload"
)
