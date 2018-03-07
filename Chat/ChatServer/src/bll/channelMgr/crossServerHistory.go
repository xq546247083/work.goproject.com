package channelMgr

import (
	"sync"
	"time"

	"work.goproject.com/Chat/ChatServer/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
)

// 跨服聊天频道历史对象
type crossServerHistory struct {
	// Id
	Id int `gorm:"column:Id;primary_key"`

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

func (this *crossServerHistory) TableName() string {
	return "history_crossserver"
}

// 转化为ResponseData，以便返回给客户端
func (this *crossServerHistory) toServerResponseData() *ServerResponseData {
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

func newCrossServerHistory(id int, channel, message, voice, fromPlayer, fromPlayerId string, crtime time.Time) *crossServerHistory {
	return &crossServerHistory{
		Id:           id,
		Channel:      channel,
		Message:      message,
		Voice:        voice,
		FromPlayer:   fromPlayer,
		FromPlayerId: fromPlayerId,
		Crtime:       crtime,
	}
}

type crossServerHistoryMgr struct {
	*config
	historyList []*crossServerHistory
	mutex       sync.Mutex
}

func newCrossServerHistoryMgr(configObj *config) *crossServerHistoryMgr {
	return &crossServerHistoryMgr{
		config:      configObj,
		historyList: make([]*crossServerHistory, 0, 2*configObj.MaxHistoryCount),
	}
}

func (this *crossServerHistoryMgr) UpdateConfig(configObj *config) {
	this.config = configObj
}

func (this *crossServerHistoryMgr) InitHistory() (err error) {
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

	debugUtil.Printf("crossServerHistoryMgr.HistoryCount:%d\n", len(historyList))
	return
}

func (this *crossServerHistoryMgr) GetHistoryInfo(playerObj *Player) interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取世界对应的历史消息列表
	if len(this.historyList) > 0 {
		return this.historyList[len(this.historyList)-1].Id
	}

	return 0
}

func (this *crossServerHistoryMgr) GetHistory(playerObj *Player, configObj *getHistoryConfig) (historyList []*ServerResponseData) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取跨服信息列表
	list := this.historyList

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

func (this *crossServerHistoryMgr) SaveHistory(chatMessageObj *ForwardChatMessage) {
	historyObj := newCrossServerHistory(chatMessageObj.Id, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player.String(), chatMessageObj.Player.Id, time.Now())

	// 第一次才需要保存数据库，如果是转发的数据则不再保存数据库，否则会主键冲突
	if chatMessageObj.IsTransfered == false {
		if result := dal.GetDB().Create(historyObj); result.Error != nil {
			dal.WriteLog("crossServerHistoryMgr.SaveHistory", result.Error)
			return
		}
	}

	this.addHistory(historyObj)
}

func (this *crossServerHistoryMgr) getHistoryList() (historyList []*crossServerHistory, err error) {
	result := dal.GetDB().Find(&historyList)
	if err = result.Error; err != nil {
		dal.WriteLog("crossServerHistoryMgr.getHistoryList", err)
		return
	}

	return
}

func (this *crossServerHistoryMgr) addHistory(historyObj *crossServerHistory) {
	if this.config.MaxHistoryCount == 0 {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取跨服信息列表
	list := this.historyList

	// 先追加到列表末尾
	list = append(list, historyObj)

	// 判断数量有没有达到指定数量的2倍？避免频繁删除
	if len(list) >= 2*this.config.MaxHistoryCount {
		// 找到临界的index，小于它的都会被删掉
		index := len(list) - this.config.MaxHistoryCount
		id := list[index].Id
		list = list[index:]

		// 删除数据库里面的数据
		this.deleteHistory(id)
	}

	// 最终保存到集合中
	this.historyList = list
}

func (this *crossServerHistoryMgr) deleteHistory(id int) (err error) {
	result := dal.GetDB().Where("Id < ?", id).Delete(crossServerHistory{})
	if err = result.Error; err != nil {
		dal.WriteLog("crossServerHistoryMgr.deleteHistory", err)
		return
	}

	return
}
