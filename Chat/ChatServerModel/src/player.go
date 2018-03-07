package src

import (
	"encoding/json"
	"sync"
	"time"

	"work.goproject.com/Framework/managecenterMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/typeUtil"
)

// 定义玩家对象
type Player struct {
	// 玩家Id
	Id string `gorm:"column:Id;primary_key"`

	// 合作商Id
	PartnerId int32 `gorm:"column:PartnerId"`

	// 服务器Id
	ServerId int32 `gorm:"column:ServerId"`

	// 玩家名称
	Name string `gorm:"column:Name"`

	// 玩家等级
	Lv int `gorm:"column:Lv"`

	// 玩家Vip等级
	Vip int `gorm:"column:Vip"`

	// 玩家登录令牌
	Token string `gorm:"column:Token"`

	// 扩展信息
	ExtendInfo string `gorm:"column:ExtendInfo"`

	// 注册时间
	RegisterTime time.Time `json:"-" gorm:"column:RegisterTime"`

	// 禁言结束时间
	SilentEndTime time.Time `json:"-" gorm:"column:SilentEndTime"`

	// 发言管理器
	*speakManager `json:"-" gorm:"-"`

	// 锁对象
	Mutex sync.Mutex `json:"-" gorm:"-"`

	// 客户端Id(由于在数据传输过程当中，需要用到ClientId来进行判断，所以需要在JSON序列化时保留)
	ClientId int32 `gorm:"-"`
}

func (this *Player) TableName() string {
	return "player"
}

func (this *Player) ServerGroupId() int32 {
	serverGroupObj, _, exists := managecenterMgr.GetServerGroup(this.PartnerId, this.ServerId)
	if !exists {
		logUtil.ErrorLog("PartnerId=%d, ServerId=%d对应的服务器组不存在", this.PartnerId, this.ServerId)
		return 0
	}
	return serverGroupObj.Id
}

func (this *Player) ServerName() string {
	_, serverObj, exists := managecenterMgr.GetServerGroup(this.PartnerId, this.ServerId)
	if !exists {
		logUtil.ErrorLog("PartnerId=%d, ServerId=%d对应的服务器不存在", this.PartnerId, this.ServerId)
		return ""
	}
	return serverObj.Name
}

func (this *Player) GetProperty(key string) (value string, err error) {
	var extendInfoMap map[string]interface{}
	if err = json.Unmarshal([]byte(this.ExtendInfo), &extendInfoMap); err != nil {
		logUtil.ErrorLog("反序列化ExtendInfo:%s失败，错误信息为：%s", this.ExtendInfo, err)
		return
	}

	mapData := typeUtil.NewMapData(extendInfoMap)
	if value, err = mapData.String(key); err != nil {
		logUtil.ErrorLog("获取玩家属性:%s失败，错误信息为：%s", key, err)
	}

	return
}

func (this *Player) ClientLogin(clientId int32) {
	this.ClientId = clientId
}

// 将玩家对象转变为字符串，以便于与客户端进行通信
func (this *Player) String() string {
	data := make(map[string]interface{})

	data["Id"] = this.Id
	data["Name"] = this.Name
	data["Lv"] = this.Lv
	data["Vip"] = this.Vip
	data["ServerName"] = this.ServerName()
	data["ServerGroupId"] = this.ServerGroupId()
	data["ExtendInfo"] = this.ExtendInfo
	if debugUtil.IsDebug() {
		data["PartnerId"] = this.PartnerId
		data["ServerId"] = this.ServerId
	}

	bytes, _ := json.Marshal(data)

	return string(bytes)
}

// 初始化玩家其它信息
func (this *Player) Init() {
	this.speakManager = newSpeakManager()
}

// 判断玩家是否处于禁言状态
// 返回值：
// 是否处于禁言状态
// 禁言剩余分钟数
func (this *Player) IsInSilent() (bool, int) {
	leftSeconds := this.SilentEndTime.Unix() - time.Now().Unix()
	if leftSeconds <= 0 {
		return false, 0
	} else {
		if leftSeconds%60 == 0 {
			return true, int(leftSeconds / 60)
		} else {
			return true, int(leftSeconds/60) + 1
		}
	}
}

func (this *Player) IsInfoChanged(target *Player) bool {
	if this.Name != target.Name ||
		this.PartnerId != target.PartnerId ||
		this.ServerId != target.ServerId ||
		this.Lv != target.Lv ||
		this.Vip != target.Vip ||
		this.Token != target.Token ||
		this.ExtendInfo != target.ExtendInfo {

		return true
	}

	return false
}

func (this *Player) UpdateInfoFromGS(gsPlayer *Player) {
	this.Name = gsPlayer.Name
	this.PartnerId = gsPlayer.PartnerId
	this.ServerId = gsPlayer.ServerId
	this.Lv = gsPlayer.Lv
	this.Vip = gsPlayer.Vip
	this.Token = gsPlayer.Token
	this.ExtendInfo = gsPlayer.ExtendInfo
}

// 初始化空对象
func NewEmptyPlayer() *Player {
	return &Player{}
}

// 利用从GS传过来的玩家对象来构造新的玩家对象
func NewPlayerFromGS(gsPlayer *Player) *Player {
	return &Player{
		Id:            gsPlayer.Id,
		PartnerId:     gsPlayer.PartnerId,
		ServerId:      gsPlayer.ServerId,
		Name:          gsPlayer.Name,
		Lv:            gsPlayer.Lv,
		Vip:           gsPlayer.Vip,
		Token:         gsPlayer.Token,
		ExtendInfo:    gsPlayer.ExtendInfo,
		RegisterTime:  time.Now(),
		SilentEndTime: time.Now(),
		speakManager:  newSpeakManager(),
	}
}

type speakManager struct {
	data  map[string]int64
	mutex sync.Mutex
}

func (this *speakManager) IsSpeakFast(channel string, interval int) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if lastTime, exists := this.data[channel]; !exists {
		return false
	} else {
		return time.Now().Unix() < lastTime+int64(interval)
	}
}

func (this *speakManager) Speak(channel string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.data[channel] = time.Now().Unix()
}

func newSpeakManager() *speakManager {
	return &speakManager{
		data: make(map[string]int64),
	}
}
