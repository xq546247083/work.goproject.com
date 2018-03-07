package channelMgr

import (
	"fmt"
	"strings"

	"work.goproject.com/Chat/ChatServer/src/bll/player"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

// 分组聊天对象
type Group struct {
	channel             string
	configName          string
	tableName           string
	newGroupHistoryFunc func() IGroupHistory
	notInGroupStatus    *ResultStatus
	*config
	*baseChannel
	*groupHistoryMgr
}

func newGroup(channel, configName, tableName string) *Group {
	return &Group{
		channel:    channel,
		configName: configName,
		tableName:  tableName,
	}
}

func (this *Group) Channel() string {
	return this.channel
}

func (this *Group) ConfigName() string {
	return this.configName
}

func (this *Group) InitConfig() (err error) {
	return initConfig(this)
}

func (this *Group) SetConfig(configObj *config) {
	this.config = configObj
}

func (this *Group) ReloadConfig() error {
	if err := this.InitConfig(); err != nil {
		return err
	}
	this.baseChannel.UpdateConfig(this.config)
	this.groupHistoryMgr.UpdateConfig(this.config)

	return nil
}

func (this *Group) InitBaseChannel() error {
	this.baseChannel = newBaseChannel(this.config, this.Channel())
	return nil
}

func (this *Group) InitHistoryMgr() error {
	this.groupHistoryMgr = newGroupHistoryMgr(this.config, this.tableName)
	this.groupHistoryMgr.setNewGroupHistoryFunc(this.newGroupHistoryFunc)

	return nil
}

func (this *Group) ValidateSendMessage(playerObj *Player, chatMessageObj *ForwardChatMessage, configObj *sendMessageConfig) *ResultStatus {
	// 判断玩家是不在群组内
	propertyValue, err := playerObj.GetProperty(this.config.PlayerProperty)
	if err != nil {
		panic(fmt.Errorf("Player.GetProperty:%s failed. Err:%s", err))
	}

	// 为了避免空的GUID，所以将-,0字符都去掉
	validateValue := strings.Replace(propertyValue, "0", "", -1)
	validateValue = strings.Replace(validateValue, "-", "", -1)

	if validateValue == "" {
		return this.notInGroupStatus
	}

	// 设置属性
	chatMessageObj.TargetProperty = propertyValue

	return Success
}

func (this *Group) GetPlayerList(chatMessageObj *ForwardChatMessage) (playerList []*Player) {
	propertyValue, err := chatMessageObj.Player.GetProperty(this.config.PlayerProperty)
	if err != nil {
		panic(fmt.Errorf("Player.GetProperty:%s failed. Err:%s", err))
	}

	if this.IsCanCrossServer() {
		for _, serverGroupPlayerObj := range player.GetServerGroupList() {
			playerList = append(playerList, serverGroupPlayerObj.GetPropertyPlayerList(this.config.PlayerProperty, propertyValue)...)
		}
	} else {
		serverGroupPlayerObj := player.GetServerGroupPlayer(chatMessageObj.ServerGroupId)
		playerList = serverGroupPlayerObj.GetPropertyPlayerList(this.config.PlayerProperty, propertyValue)
	}

	return
}

func (this *Group) NewServerResponseData(chatMessageObj *ForwardChatMessage) *ServerResponseData {
	return NewServerResponseData(chatMessageObj.Id, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player, nil)
}

// ---------------------------------------非接口方法---------------------------------------

func (this *Group) setNotInGroupStatus(status *ResultStatus) {
	this.notInGroupStatus = status
}

func (this *Group) setNewGroupHistoryFunc(newGroupHistoryFunc func() IGroupHistory) {
	this.newGroupHistoryFunc = newGroupHistoryFunc
}
