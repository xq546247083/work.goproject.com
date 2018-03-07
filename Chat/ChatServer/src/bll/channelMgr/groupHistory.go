package channelMgr

import (
	"fmt"
	"reflect"
	"sync"

	"work.goproject.com/Chat/ChatServer/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
)

// 所有以分组形式进行的聊天，均需要实现此接口。如国家聊天、公会聊天、组队聊天等。
// 其中的GetIdentifier()返回的数据必须是全局唯一的。
// 如值为GUID的公会Id；
// 如果是国家Id这样的非全局唯一的，则以ServerGroupId_NationCode这样的形式来生成全局唯一的Identifier。
type IGroupHistory interface {
	GetId() int
	GetIdentifier() string
	SetIdentifier(int, string)
	SetChannel(string)
	SetMessage(string, string)
	SetFromPlayer(string, string)
	ToServerResponseData() *ServerResponseData
}

type groupHistoryMgr struct {
	*config
	tableName           string
	newGroupHistoryFunc func() IGroupHistory
	historyMap          map[string][]IGroupHistory
	mutex               sync.Mutex
}

func newGroupHistoryMgr(configObj *config, tableName string) *groupHistoryMgr {
	return &groupHistoryMgr{
		config:     configObj,
		tableName:  tableName,
		historyMap: make(map[string][]IGroupHistory, 32),
	}
}

func (this *groupHistoryMgr) mapKey(playerObj *Player) string {
	propertyValue, err := playerObj.GetProperty(this.config.PlayerProperty)
	if err != nil {
		panic(fmt.Errorf("Player.GetProperty:%s failed. Err:%s", err))
	}

	return propertyValue
}

func (this *groupHistoryMgr) UpdateConfig(configObj *config) {
	this.config = configObj
}

func (this *groupHistoryMgr) InitHistory() (err error) {
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

	debugUtil.Printf("groupHistoryMgr.HistoryCount:%d\n", len(historyList))
	return
}

func (this *groupHistoryMgr) GetHistoryInfo(playerObj *Player) interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取世界对应的历史消息列表
	if list, exists := this.historyMap[this.mapKey(playerObj)]; exists {
		if len(list) > 0 {
			return list[len(list)-1].GetId()
		}
	}

	return 0
}

func (this *groupHistoryMgr) GetHistory(playerObj *Player, configObj *getHistoryConfig) (historyList []*ServerResponseData) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取服务器组对应的历史消息列表
	var list []IGroupHistory
	var exists bool
	if list, exists = this.historyMap[this.mapKey(playerObj)]; !exists {
		return
	}

	// 只取指定数量的数据
	if len(list) > this.config.MaxHistoryCount {
		list = list[len(list)-this.config.MaxHistoryCount:]
	}

	// 找到Id<messageId的第一条数据所在的位置
	maxIndex := len(list) - 1
	for ; maxIndex >= 0; maxIndex-- {
		if list[maxIndex].GetId() < configObj.MessageId {
			break
		}
	}

	// 根据count来进行过滤
	for index := maxIndex; index >= 0 && configObj.Count > 0; index-- {
		historyList = append(historyList, list[index].ToServerResponseData())
		configObj.Count--
	}

	return
}

func (this *groupHistoryMgr) SaveHistory(chatMessageObj *ForwardChatMessage) {
	historyObj := this.newGroupHistoryFunc()
	historyObj.SetIdentifier(chatMessageObj.Id, chatMessageObj.TargetProperty.(string))
	historyObj.SetChannel(chatMessageObj.Channel)
	historyObj.SetMessage(chatMessageObj.Message, chatMessageObj.Voice)
	historyObj.SetFromPlayer(chatMessageObj.Player.String(), chatMessageObj.Player.Id)

	// 第一次才需要保存数据库，如果是转发的数据则不再保存数据库，否则会主键冲突
	if chatMessageObj.IsTransfered == false {
		if result := dal.GetDB().Table(this.tableName).Create(historyObj); result.Error != nil {
			dal.WriteLog("groupHistoryMgr.SaveHistory", result.Error)
			return
		}
	}

	this.addHistory(historyObj)
}

func (this *groupHistoryMgr) getHistoryList() (historyList []IGroupHistory, err error) {
	historyObj := this.newGroupHistoryFunc()
	historyObjType := reflect.TypeOf(historyObj)
	historySlice := reflect.New(reflect.SliceOf(historyObjType)).Interface()

	result := dal.GetDB().Table(this.tableName).Find(historySlice)
	if err = result.Error; err != nil {
		dal.WriteLog("groupHistoryMgr.getHistoryList", err)
		return
	}

	resultSlice := reflect.ValueOf(historySlice).Elem()
	if resultSlice.IsNil() || resultSlice.Len() <= 0 {
		return
	}

	for i := 0; i < resultSlice.Len(); i++ {
		historyList = append(historyList, resultSlice.Index(i).Interface().(IGroupHistory))
	}

	return
}

func (this *groupHistoryMgr) addHistory(historyObj IGroupHistory) {
	if this.config.MaxHistoryCount == 0 {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 获取分组对应的历史消息列表
	var list []IGroupHistory
	var exists bool
	if list, exists = this.historyMap[historyObj.GetIdentifier()]; !exists {
		list = make([]IGroupHistory, 0, 2*this.config.MaxHistoryCount)
	}

	// 先追加到列表末尾
	list = append(list, historyObj)

	// 判断数量有没有达到指定数量的2倍？避免频繁删除
	if len(list) >= 2*this.config.MaxHistoryCount {
		// 找到临界的index，小于它的都会被删掉
		index := len(list) - this.config.MaxHistoryCount
		id := list[index].GetId()
		list = list[index:]

		// 删除数据库里面的数据
		this.deleteHistory(id, historyObj.GetIdentifier())
	}

	// 最终保存到集合中
	this.historyMap[historyObj.GetIdentifier()] = list
}

func (this *groupHistoryMgr) deleteHistory(id int, identifier string) (err error) {
	result := dal.GetDB().Table(this.tableName).Where("Identifier = ? AND Id < ?", identifier, id).Delete(this.newGroupHistoryFunc())
	if err = result.Error; err != nil {
		dal.WriteLog("groupHistoryMgr.deleteHistory", err)
		return
	}

	return
}

// ---------------------------------------非接口方法---------------------------------------

func (this *groupHistoryMgr) setNewGroupHistoryFunc(newGroupHistoryFunc func() IGroupHistory) {
	this.newGroupHistoryFunc = newGroupHistoryFunc
}
