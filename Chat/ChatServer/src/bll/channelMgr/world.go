package channelMgr

import (
	"work.goproject.com/Chat/ChatServer/src/bll/player"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func init() {
	register(newWorld())
}

// 世界聊天对象
type World struct {
	*config
	*baseChannel
	*worldHistoryMgr
}

func newWorld() *World {
	return &World{}
}

func (this *World) Channel() string {
	return "World"
}

func (this *World) ConfigName() string {
	return "ChannelConfig_World"
}

func (this *World) InitConfig() (err error) {
	return initConfig(this)
}

func (this *World) SetConfig(configObj *config) {
	this.config = configObj
}

func (this *World) ReloadConfig() error {
	if err := this.InitConfig(); err != nil {
		return err
	}
	this.baseChannel.UpdateConfig(this.config)
	this.worldHistoryMgr.UpdateConfig(this.config)

	return nil
}

func (this *World) InitBaseChannel() error {
	this.baseChannel = newBaseChannel(this.config, this.Channel())
	return nil
}

func (this *World) InitHistoryMgr() error {
	this.worldHistoryMgr = newWorldHistoryMgr(this.config)
	return nil
}

func (this *World) GetPlayerList(chatMessageObj *ForwardChatMessage) (playerList []*Player) {
	// World频道不可以跨服，所以只获取本服的玩家
	serverGroupPlayerObj := player.GetServerGroupPlayer(chatMessageObj.ServerGroupId)
	playerList = serverGroupPlayerObj.GetPlayerList()

	return
}

func (this *World) NewServerResponseData(chatMessageObj *ForwardChatMessage) *ServerResponseData {
	return NewServerResponseData(chatMessageObj.Id, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player, nil)
}
