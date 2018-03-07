package channelMgr

import (
	"sync"
	"time"

	"work.goproject.com/Chat/ChatServer/src/bll/player"
	"work.goproject.com/Chat/ChatServer/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
)

// 私聊历史对象
type privateHistory struct {
	// Id
	Id int `gorm:"column:Id;primary_key"`

	// 私聊消息的接收者Id
	PlayerId string `gorm:"column:PlayerId;primary_key"`

	// 聊天渠道
	Channel string `gorm:"column:Channel"`

	// 聊天消息
	Message string `gorm:"column:Message"`

	// 语音信息
	Voice string `gorm:"column:Voice"`

	// 源玩家
	FromPlayer string `gorm:"column:FromPlayer"`

	// 源玩家Id
	FromPlayerId string `gorm:"column:FromPlayerId"`

	// 目标玩家对象
	ToPlayer string `gorm:"column:ToPlayer"`

	// 目标玩家Id
	ToPlayerId string `gorm:"column:ToPlayerId"`

	// 创建时间
	Crtime time.Time `gorm:"column:Crtime"`
}

func (this *privateHistory) TableName() string {
	return "history_private"
}

// 转化为ResponseData，以便返回给客户端
func (this *privateHistory) toServerResponseData() *ServerResponseData {
	return ConvertToServerResponseData(
		this.Id,
		this.Channel,
		this.Message,
		this.Voice,
		this.FromPlayer,
		this.FromPlayerId,
		this.ToPlayer,
		this.ToPlayerId,
		this.Crtime.Unix(),
	)
}

func (this *privateHistory) setPlayerId(playerId string) {
	this.PlayerId = playerId
}

func (this *privateHistory) isTargetPlayer(targetPlayerId string) bool {
	if this.FromPlayerId == targetPlayerId || this.ToPlayerId == targetPlayerId {
		return true
	}

	return false
}

func newPrivateHistory(id int, channel, message, voice, fromPlayer, fromPlayerId, toPlayer, toPlayerId string, crtime time.Time) *privateHistory {
	return &privateHistory{
		Id:           id,
		PlayerId:     fromPlayerId,
		Channel:      channel,
		Message:      message,
		Voice:        voice,
		FromPlayer:   fromPlayer,
		FromPlayerId: fromPlayerId,
		ToPlayer:     toPlayer,
		ToPlayerId:   toPlayerId,
		Crtime:       crtime,
	}
}

type privateHistoryMgr struct {
	*config
	historyMap map[string][]*privateHistory
	mutex      sync.Mutex
}

func newPrivateHistoryMgr(configObj *config) *privateHistoryMgr {
	this := &privateHistoryMgr{
		config:     configObj,
		historyMap: make(map[string][]*privateHistory, 32),
	}

	// 注册玩家注册、反注册触发方法
	player.AddPlayerRegisterTriggerFunc("initPlayerHistory", this.initPlayerHistory)
	player.AddPlayerUnregisterTriggerFunc("clearPlayerHistory", this.clearPlayerHistory)

	return this
}

func (this *privateHistoryMgr) UpdateConfig(configObj *config) {
	this.config = configObj
}

func (this *Private) InitHistory() (err error) {
	// 私聊不提供全局的初始化方法
	return nil
}

func (this *privateHistoryMgr) GetHistoryInfo(playerObj *Player) interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	/*
		获取每天私聊数据，获取TargetPlayerId以及对应的Id
		用后面的数据覆盖前面的数据
	*/
	data := make(map[string]int)
	if list, exists := this.historyMap[playerObj.Id]; exists {
		for _, item := range list {
			if playerObj.Id != item.FromPlayerId {
				data[item.FromPlayerId] = item.Id
			} else {
				data[item.ToPlayerId] = item.Id
			}
		}
	}

	return data
}

func (this *privateHistoryMgr) GetHistory(playerObj *Player, configObj *getHistoryConfig) (historyList []*ServerResponseData) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取玩家对应的历史消息列表
	var list []*privateHistory
	var exists bool
	if list, exists = this.historyMap[playerObj.Id]; !exists {
		return
	}

	// 找到与目标玩家的所有聊天记录
	var targetList []*privateHistory = make([]*privateHistory, 0, 32)
	for _, item := range list {
		if item.isTargetPlayer(configObj.TargetPlayerId) {
			targetList = append(targetList, item)
		}
	}

	if len(targetList) == 0 {
		return
	}

	// 找到Id<messageId的第一条数据所在的位置
	maxIndex := len(targetList) - 1
	for ; maxIndex >= 0; maxIndex-- {
		if targetList[maxIndex].Id < configObj.MessageId {
			break
		}
	}

	// 根据count来进行过滤
	for index := maxIndex; index >= 0 && configObj.Count > 0; index-- {
		historyList = append(historyList, targetList[index].toServerResponseData())
		configObj.Count--
	}

	return
}

func (this *privateHistoryMgr) SaveHistory(chatMessageObj *ForwardChatMessage) {
	var toPlayerId, toPlayer string
	if toPlayerObj, ok := chatMessageObj.TargetProperty.(*Player); ok {
		toPlayerId = toPlayerObj.Id
		toPlayer = toPlayerObj.String()
	}

	now := time.Now()
	//(id int, channel, message, voice, fromPlayer, fromPlayerId, toPlayer, toPlayerId string, crtime time.Time)
	senderHistoryObj := newPrivateHistory(chatMessageObj.Id, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player.String(), chatMessageObj.Player.Id, toPlayer, toPlayerId, now)
	senderHistoryObj.setPlayerId(senderHistoryObj.FromPlayerId)
	receiverHistoryObj := newPrivateHistory(chatMessageObj.Id, chatMessageObj.Channel, chatMessageObj.Message, chatMessageObj.Voice, chatMessageObj.Player.String(), chatMessageObj.Player.Id, toPlayer, toPlayerId, now)
	receiverHistoryObj.setPlayerId(senderHistoryObj.ToPlayerId)

	// 第一次才需要保存数据库，如果是转发的数据则不再保存数据库，否则会主键冲突
	if chatMessageObj.IsTransfered == false {
		this.insert(senderHistoryObj)
		this.insert(receiverHistoryObj)
	}

	this.addHistory(senderHistoryObj)
	this.addHistory(receiverHistoryObj)
}

// 删除私聊历史
// playerId：玩家Id
// targetPlayerId：发送玩家Id
func (this *privateHistoryMgr) DeleteHistory(playerObj *Player, configObj *deleteHistoryConfig) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 查找玩家的私聊消息列表
	if list, exists := this.historyMap[playerObj.Id]; exists {
		// 查找和目标玩家的私聊消息
		count := 0
		newList := make([]*privateHistory, 0, 32)
		for _, item := range list {
			if item.isTargetPlayer(configObj.TargetPlayerId) {
				count++
			} else {
				newList = append(newList, item)
			}
		}

		// 判断是否找到有匹配的数据
		if count > 0 {
			this.historyMap[playerObj.Id] = newList
			this.deleteHistory(playerObj.Id, configObj.TargetPlayerId)
		}
	}
}

// 添加私聊历史
func (this *privateHistoryMgr) addHistory(historyObj *privateHistory) {
	if this.config.MaxHistoryCount == 0 {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取玩家对应的历史消息列表
	var list []*privateHistory
	var exists bool
	if list, exists = this.historyMap[historyObj.PlayerId]; !exists {
		// 如果找不到对应的玩家，则证明玩家不在此服务器
		return
	}

	// 先追加到列表末尾
	list = append(list, historyObj)

	// 判断有没有过期的数据(超过20条才进行判断，否则太频繁了)
	if len(list) > 20 {
		expireTime := time.Now()
		expireCount := 0
		expireIndex := 0
		for index, item := range list {
			expireIndex = index
			if item.Crtime.AddDate(0, 0, this.config.MaxHistoryCount).Unix() <= time.Now().Unix() {
				expireCount++
				expireTime = item.Crtime
			} else {
				break
			}
		}

		// 将过期的数据截断
		if expireCount > 0 {
			list = list[expireIndex:]
		}

		// 过期的数量大于5才删除数据库，避免太频繁地操作数据库
		if expireCount > 5 {
			this.deleteExpiredHistory(historyObj.PlayerId, expireTime)
		}
	}

	// 最终保存到集合中
	this.historyMap[historyObj.PlayerId] = list
}

// 获取历史消息列表
func (this *privateHistoryMgr) getHistoryList(playerId string, expireTime time.Time) (historyList []*privateHistory, err error) {
	result := dal.GetDB().Where("PlayerId = ? AND Crtime > ?", playerId, expireTime).Find(&historyList)
	if err = result.Error; err != nil {
		dal.WriteLog("privateHistoryMgr.getHistoryList", err)
		return
	}

	return
}

func (this *privateHistoryMgr) insert(historyObj *privateHistory) (err error) {
	result := dal.GetDB().Create(&historyObj)
	if err = result.Error; err != nil {
		debugUtil.Printf("history:%v\n", historyObj)
		dal.WriteLog("privateHistoryMgr.insert", err)
		return
	}

	return
}

// 删除过期的私聊历史
func (this *privateHistoryMgr) deleteExpiredHistory(playerId string, expireTime time.Time) (err error) {
	result := dal.GetDB().Where("PlayerId = ? AND Crtime < ?", playerId, expireTime).Delete(privateHistory{})
	if err = result.Error; err != nil {
		dal.WriteLog("privateHistoryMgr.deleteExpiredHistory", err)
		return
	}

	return
}

// 删除私聊历史
func (this *privateHistoryMgr) deleteHistory(playerId, targetPlayerId string) (err error) {
	result := dal.GetDB().Where("PlayerId = ? AND (FromPlayerId = ? OR ToPlayerId = ?)", playerId, targetPlayerId, targetPlayerId).Delete(privateHistory{})
	if err = result.Error; err != nil {
		dal.WriteLog("privateHistoryMgr.deleteHistory", err)
		return
	}

	return
}

// 初始化玩家私聊历史
func (this *privateHistoryMgr) initPlayerHistory(playerObj *Player) {
	if this.config.MaxHistoryCount == 0 {
		return
	}

	initPrivateHistory := func(playerId string) {
		this.mutex.Lock()
		defer this.mutex.Unlock()

		this.historyMap[playerId] = make([]*privateHistory, 0, 32)
	}

	expireTime := time.Now().AddDate(0, 0, -1*this.config.MaxHistoryCount)
	historyList, err := this.getHistoryList(playerObj.Id, expireTime)
	if err != nil {
		return
	}
	debugUtil.Printf("playerId:%s,historyList count:%d\n", playerObj.Id, len(historyList))

	// 初始化空的列表进行占位处理
	initPrivateHistory(playerObj.Id)
	for _, item := range historyList {
		this.addHistory(item)
	}

	debugUtil.Printf("privateHistoryMgr.HistoryCount: %d\n", len(historyList))
}

// 清理玩家私聊历史
// playerObj:玩家对象
func (this *privateHistoryMgr) clearPlayerHistory(playerObj *Player) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.historyMap, playerObj.Id)
}
