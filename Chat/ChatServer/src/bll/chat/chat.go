package chat

import (
	"work.goproject.com/Chat/ChatServer/src/bll/channelMgr"
	"work.goproject.com/Chat/ChatServer/src/bll/dbConfig"
	"work.goproject.com/Chat/ChatServer/src/bll/messageLog"
	"work.goproject.com/Chat/ChatServer/src/bll/player"
	"work.goproject.com/Chat/ChatServer/src/clientMgr"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/managecenterMgr"
)

func init() {
	clientMgr.RegisterFunction(new(ChatModule))
}

// 聊天模块
type ChatModule int8

// 登陆
// clientObj:客户端对象
// id:玩家Id
// token:令牌信息
// partnerId:合作商Id
// serverId:服务器Id
func (this *ChatModule) Login(clientObj clientMgr.IClient, id, token string, partnerId, serverId int32) ServerResponseObject {
	responseObj := NewServerResponseObject()

	// 定义变量
	var err error
	var exists bool
	var playerObj *Player

	// 判断服务器组是否存在
	if _, _, exists = managecenterMgr.GetServerGroup(partnerId, serverId); !exists {
		return *responseObj.SetResultStatus(ServerGroupNotExist)
	}

	// 寻找玩家
	if playerObj, exists, err = player.GetPlayer(id, true); err != nil {
		return *responseObj.SetResultStatus(DataError)
	} else if !exists {
		return *responseObj.SetResultStatus(PlayerNotExists)
	}

	// 验证token是否正确
	if playerObj.Token != token {
		return *responseObj.SetResultStatus(TokenInvalid)
	}

	// 判断是否重复登陆
	if playerObj.ClientId > 0 {
		if oldClientObj, exists := clientMgr.GetClient(playerObj.ClientId); exists {
			// 如果不是同一个客户端，则先给客户端发送在其他设备登陆信息，然后断开连接
			if clientObj != oldClientObj {
				sendLoginAnotherDeviceMsg(oldClientObj)
			}
		}
	}

	// 更新客户端对象的玩家Id、以及玩家对象对应的客户端Id
	clientObj.PlayerLogin(id)
	playerObj.ClientLogin(clientObj.GetId())

	// 将玩家对象添加到玩家列表中
	player.RegisterPlayer(playerObj)

	return *responseObj
}

// 发送消息
// clientObj:客户端对象
// playerObj:玩家对象
// channelTypeValue:聊天频道
// message:消息内容
// voice:语音信息
// toPlayerId:私聊玩家Id
func (this *ChatModule) SendMessage(clientObj clientMgr.IClient, playerObj *Player, channel, message, voice, toPlayerId string) ServerResponseObject {
	responseObj := NewServerResponseObject()

	// 判断是否被禁言
	if isInSilent, _ := playerObj.IsInSilent(); isInSilent {
		return *responseObj.SetResultStatus(InSilent)
	}

	// 获取聊天频道对象
	channelObj, exists := channelMgr.GetChannel(channel)
	if !exists {
		return *responseObj.SetResultStatus(ChannelNotDefined)
	}

	// 判断玩家等级是否达到系统开放等级
	if channelObj.IsOpen(playerObj) == false {
		return *responseObj.SetResultStatus(LvIsNotEnough)
	}

	// 判断玩家说话是否太快
	if channelObj.IsSpeakFast(playerObj) {
		return *responseObj.SetResultStatus(SendMessageTooFast)
	}

	// 构造用于记录日志和转发的聊天消息对象
	chatMessageObj := NewForwardChatMessage(message, voice, playerObj.ServerGroupId(), channel, playerObj)

	// 验证参数
	sendMessageConfig := channelMgr.NewSendMessageConfig(message, voice, toPlayerId)
	if rs := channelObj.ValidateSendMessage(playerObj, chatMessageObj, sendMessageConfig); rs != Success {
		return *responseObj.SetResultStatus(rs)
	}

	// 处理消息
	channelObj.HandleMessage(chatMessageObj)

	// 更新说话信息
	channelObj.Speak(playerObj)

	// 保存日志
	id, err := messageLog.Save(playerObj.Id, playerObj.Name, playerObj.PartnerId, playerObj.ServerId, playerObj.ServerGroupId(), message, voice, channel, toPlayerId)
	if err != nil {
		return *responseObj.SetResultStatus(DataError)
	} else {
		// 更新消息Id
		chatMessageObj.Id = id
	}

	// 发送到通道中，进行后续处理
	chatMessageChannel <- chatMessageObj

	return *responseObj
}

// 获取历史记录信息
func (this *ChatModule) GetHistoryInfo(clientObj clientMgr.IClient, playerObj *Player) ServerResponseObject {
	responseObj := NewServerResponseObject()

	// 组装返回值
	data := make(map[string]interface{})

	// 遍历所有的聊天频道，并调用对应的方法
	for channel, _ := range dbConfig.ChannelConfig {
		if channelObj, exists := channelMgr.GetChannel(channel); exists {
			data[channel] = channelObj.GetHistoryInfo(playerObj)
		}
	}

	responseObj.SetData(data)

	return *responseObj
}

// 获取历史记录
// clientObj:客户端对象
// playerObj:玩家对象
// channelTypeValue:聊天频道
// messageId:消息Id（获取消息时，要取Id<messageId的数据）
// count:获取的消息数量
// targetPlayerId:目标玩家Id
func (this *ChatModule) GetHistory(clientObj clientMgr.IClient, playerObj *Player, channel string, messageId, count int, targetPlayerId string) ServerResponseObject {
	responseObj := NewServerResponseObject()

	// 获取聊天频道对象
	channelObj, exists := channelMgr.GetChannel(channel)
	if !exists {
		return *responseObj.SetResultStatus(ChannelNotDefined)
	}

	// 验证参数
	if rs := channelObj.ValidateHistoryParam(playerObj); rs != Success {
		return *responseObj.SetResultStatus(rs)
	}

	// 组装返回值
	getHistoryConfig := channelMgr.NewGetHistoryConfig(messageId, count, targetPlayerId)
	responseObj.SetData(channelObj.GetHistory(playerObj, getHistoryConfig))

	return *responseObj
}

// 删除私聊历史记录
// clientObj:客户端对象
// playerObj:玩家对象
// targetPlayerId:目标玩家Id
func (this *ChatModule) DeletePrivateHistory(clientObj clientMgr.IClient, playerObj *Player, targetPlayerId string) ServerResponseObject {
	responseObj := NewServerResponseObject()

	// 获取聊天频道对象
	channelObj, exists := channelMgr.GetPrivateChannel()
	if !exists {
		return *responseObj.SetResultStatus(ChannelNotDefined)
	}

	// 删除私聊历史
	deleteHistoryConfig := channelMgr.NewDeleteHistoryConfig(targetPlayerId)
	channelObj.DeleteHistory(playerObj, deleteHistoryConfig)

	// 组装返回值
	responseObj.SetData(targetPlayerId)

	return *responseObj
}
