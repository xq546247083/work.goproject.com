package channelMgr

import (
	"work.goproject.com/Chat/ChatServer/src/bll/player"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/logUtil"
)

func init() {
	register(newPrivate())
}

// 私聊聊天对象
type Private struct {
	*config
	*baseChannel
	*privateHistoryMgr
}

func newPrivate() *Private {
	return &Private{}
}

func (this *Private) Channel() string {
	return "Private"
}

func (this *Private) ConfigName() string {
	return "ChannelConfig_Private"
}

func (this *Private) InitConfig() (err error) {
	return initConfig(this)
}

func (this *Private) SetConfig(configObj *config) {
	this.config = configObj
}

func (this *Private) ReloadConfig() error {
	if err := this.InitConfig(); err != nil {
		return err
	}
	this.baseChannel.UpdateConfig(this.config)
	this.privateHistoryMgr.UpdateConfig(this.config)

	return nil
}

func (this *Private) InitBaseChannel() error {
	this.baseChannel = newBaseChannel(this.config, this.Channel())
	return nil
}

func (this *Private) InitHistoryMgr() error {
	this.privateHistoryMgr = newPrivateHistoryMgr(this.config)
	return nil
}

func (this *Private) ValidateSendMessage(playerObj *Player, chatMessageObj *ForwardChatMessage, configObj *sendMessageConfig) *ResultStatus {
	// 目标玩家Id不能为空
	if configObj.ToPlayerId == "" {
		return NotFoundTarget
	}

	// 不能给自己发送消息
	if configObj.ToPlayerId == playerObj.Id {
		return CantSendMessageToSelf
	}

	if toPlayerObj, exists, err := player.GetPlayer(configObj.ToPlayerId, true); err != nil {
		return DataError
	} else if !exists {
		return PlayerNotExist
	} else {
		// 如果私聊不能跨服，则判断两个玩家的ServerGroup是不相等
		if this.IsCanCrossServer() == false && toPlayerObj.ServerGroupId() != playerObj.ServerGroupId() {
			return CantCrossServerTalk
		}

		chatMessageObj.TargetProperty = toPlayerObj
	}

	return Success
}

func (this *Private) GetPlayerList(chatMessageObj *ForwardChatMessage) (playerList []*Player) {
	// 如果没有被转发过，则需要添加发送者；如果是转发的，则发送者肯定不在此服务器中，就无需添加到列表中
	if chatMessageObj.IsTransfered == false {
		playerList = append(playerList, chatMessageObj.Player)
	}

	if toPlayer, ok := chatMessageObj.TargetProperty.(*Player); ok {
		playerList = append(playerList, toPlayer)
	} else {
		logUtil.ErrorLog("将chatMessageObj.TargetProperty转换为Player对象失败。")
	}

	return
}

func (this *Private) NewServerResponseData(chatMessageObj *ForwardChatMessage) *ServerResponseData {
	toPlayer, _ := chatMessageObj.TargetProperty.(*Player)
	return NewServerResponseData(chatMessageObj.Id, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player, toPlayer)
}

func (this *Private) DeleteHistory(playerObj *Player, configObj *deleteHistoryConfig) {
	this.privateHistoryMgr.DeleteHistory(playerObj, configObj)
}
