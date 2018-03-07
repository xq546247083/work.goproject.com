package channelMgr

import (
	"work.goproject.com/Chat/ChatServer/src/bll/player"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func init() {
	register(newCrossServer())
}

// 跨服聊天对象
type CrossServer struct {
	*config
	*baseChannel
	*crossServerHistoryMgr
}

func newCrossServer() *CrossServer {
	return &CrossServer{}
}

func (this *CrossServer) Channel() string {
	return "CrossServer"
}

func (this *CrossServer) ConfigName() string {
	return "ChannelConfig_CrossServer"
}

func (this *CrossServer) InitConfig() (err error) {
	return initConfig(this)
}

func (this *CrossServer) SetConfig(configObj *config) {
	this.config = configObj
}

func (this *CrossServer) ReloadConfig() error {
	if err := this.InitConfig(); err != nil {
		return err
	}
	this.baseChannel.UpdateConfig(this.config)
	this.crossServerHistoryMgr.UpdateConfig(this.config)

	return nil
}

func (this *CrossServer) InitBaseChannel() error {
	this.baseChannel = newBaseChannel(this.config, this.Channel())
	return nil
}

func (this *CrossServer) InitHistoryMgr() error {
	this.crossServerHistoryMgr = newCrossServerHistoryMgr(this.config)
	return nil
}

func (this *CrossServer) GetPlayerList(chatMessageObj *ForwardChatMessage) (playerList []*Player) {
	for _, serverGroupPlayerObj := range player.GetServerGroupList() {
		playerList = append(playerList, serverGroupPlayerObj.GetPlayerList()...)
	}

	return
}

func (this *CrossServer) NewServerResponseData(chatMessageObj *ForwardChatMessage) *ServerResponseData {
	return NewServerResponseData(chatMessageObj.Id, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player, nil)
}
