package channelMgr

import (
	"sync"
	"time"

	"work.goproject.com/Chat/ChatServer/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
)

// 世界聊天频道历史对象
type worldHistory struct {
	// Id
	Id int `gorm:"column:Id;primary_key"`

	// 服务器组Id信息
	ServerGroupId int32 `gorm:"column:ServerGroupId"`

	// 聊天渠道
	Channel string `gorm:"column:Channel"`

	// 聊天消息
	Message string `gorm:"column:Message"`

	// 语音信息
	Voice string `gorm:"column:Voice"`

	// 说话者
	FromPlayer   string `gorm:"column:FromPlayer"`
	FromPlayerId string `gorm:"column:FromPlayerId"`

	// 创建时间
	Crtime time.Time `gorm:"column:Crtime"`
}

func (this *worldHistory) TableName() string {
	return "history_world"
}

// 转化为ResponseData，以便返回给客户端
func (this *worldHistory) toServerResponseData() *ServerResponseData {
	return ConvertToServerResponseData(
		this.Id,
		this.Channel,
		this.Message,
		this.Voice,
		this.FromPlayer,
		this.FromPlayerId,
		"",
		"",
		this.Crtime.Unix(),
	)
}

func newWorldHistory(id int, serverGroupId int32, channel, message, voice, fromPlayer, fromPlayerId string, crtime time.Time) *worldHistory {
	return &worldHistory{
		Id:            id,
		ServerGroupId: serverGroupId,
		Channel:       channel,
		Message:       message,
		Voice:         voice,
		FromPlayer:    fromPlayer,
		FromPlayerId:  fromPlayerId,
		Crtime:        crtime,
	}
}

type worldHistoryMgr struct {
	*config
	historyMap map[int32][]*worldHistory
	mutex      sync.Mutex
}

func newWorldHistoryMgr(configObj *config) *worldHistoryMgr {
	return &worldHistoryMgr{
		config:     configObj,
		historyMap: make(map[int32][]*worldHistory, 2*configObj.MaxHistoryCount),
	}
}

func (this *worldHistoryMgr) UpdateConfig(configObj *config) {
	this.config = configObj
}

func (this *worldHistoryMgr) InitHistory() (err error) {
	if this.config.MaxHistoryCount == 0 {
		return
	}

	historyList, err := this.getHistoryList()
	if err != nil {
		return
	}

	for _, item := range historyList {
		this.addHistory(item)
	}

	debugUtil.Printf("worldHistoryMgr.HistoryCount:%d\n", len(historyList))
	return
}

func (this *worldHistoryMgr) GetHistoryInfo(playerObj *Player) interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取世界对应的历史消息列表
	if list, exists := this.historyMap[playerObj.ServerGroupId()]; exists {
		if len(list) > 0 {
			return list[len(list)-1].Id
		}
	}

	return 0
}

func (this *worldHistoryMgr) GetHistory(playerObj *Player, configObj *getHistoryConfig) (historyList []*ServerResponseData) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取服务器组对应的历史消息列表
	var list []*worldHistory
	var exists bool
	if list, exists = this.historyMap[playerObj.ServerGroupId()]; !exists {
		return
	}

	// 只取指定数量的数据
	if len(list) > this.config.MaxHistoryCount {
		list = list[len(list)-this.config.MaxHistoryCount:]
	}

	// 找到Id<messageId的第一条数据所在的位置
	maxIndex := len(list) - 1
	for ; maxIndex >= 0; maxIndex-- {
		if list[maxIndex].Id < configObj.MessageId {
			break
		}
	}

	// 根据count来进行过滤
	for index := maxIndex; index >= 0 && configObj.Count > 0; index-- {
		historyList = append(historyList, list[index].toServerResponseData())
		configObj.Count--
	}

	return
}

func (this *worldHistoryMgr) SaveHistory(chatMessageObj *ForwardChatMessage) {
	historyObj := newWorldHistory(chatMessageObj.Id, chatMessageObj.ServerGroupId, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player.String(), chatMessageObj.Player.Id, time.Now())

	// 第一次才需要保存数据库，如果是转发的数据则不再保存数据库，否则会主键冲突
	if chatMessageObj.IsTransfered == false {
		if result := dal.GetDB().Create(historyObj); result.Error != nil {
			dal.WriteLog("worldHistoryMgr.SaveHistory", result.Error)
			return
		}
	}

	this.addHistory(historyObj)
}

func (this *worldHistoryMgr) getHistoryList() (historyList []*worldHistory, err error) {
	result := dal.GetDB().Find(&historyList)
	if err = result.Error; err != nil {
		dal.WriteLog("worldHistoryMgr.getHistoryList", err)
		return
	}

	return
}

func (this *worldHistoryMgr) addHistory(historyObj *worldHistory) {
	if this.config.MaxHistoryCount == 0 {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取世界对应的历史消息列表
	var list []*worldHistory
	var exists bool
	if list, exists = this.historyMap[historyObj.ServerGroupId]; !exists {
		list = make([]*worldHistory, 0, 2*this.config.MaxHistoryCount)
	}

	// 先追加到列表末尾
	list = append(list, historyObj)

	// 判断数量有没有达到指定数量的2倍？避免频繁删除
	if len(list) >= 2*this.config.MaxHistoryCount {
		// 找到临界的index，小于它的都会被删掉
		index := len(list) - this.config.MaxHistoryCount
		id := list[index].Id
		list = list[index:]

		// 删除数据库里面的数据
		this.deleteHistory(id, historyObj.ServerGroupId)
	}

	// 最终保存到集合中
	this.historyMap[historyObj.ServerGroupId] = list
}

func (this *worldHistoryMgr) deleteHistory(id int, serverGroupId int32) (err error) {
	result := dal.GetDB().Where("ServerGroupId = ? AND Id < ?", serverGroupId, id).Delete(worldHistory{})
	if err = result.Error; err != nil {
		dal.WriteLog("worldHistoryMgr.deleteHistory", err)
		return
	}

	return
}
