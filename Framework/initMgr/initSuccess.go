package initMgr

import (
	"fmt"
	"sync"

	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

// 初始化成功对象
type InitSuccess struct {
	name        string
	registerMap map[string]chan bool
	mutex       sync.RWMutex
}

// 注册需要被通知的对象
// name:唯一标识
// ch:用于通知的通道
func (this *InitSuccess) Register(name string, ch chan bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, exists := this.registerMap[name]; exists {
		panic(fmt.Errorf("InitSuccess.Register-%s-%s已经存在，请检查", this.name, name))
	}

	this.registerMap[name] = ch

	logUtil.DebugLog("InitSuccess.Register-%s，当前共有%d个注册", this.name, len(this.registerMap))
}

// 取消启动成功通知注册
// name:唯一标识
func (this *InitSuccess) Unregister(name string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	delete(this.registerMap, name)
}

// 通知所有已注册的对象
func (this *InitSuccess) Notify() {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	for name, ch := range this.registerMap {
		ch <- true
		msg := fmt.Sprintf("通知:%s-%s初始化成功", this.name, name)
		debugUtil.Println(msg)
		logUtil.DebugLog(msg)
	}
}

// 创建初始化成功对象
func NewInitSuccess(name string) *InitSuccess {
	return &InitSuccess{
		name:        name,
		registerMap: make(map[string]chan bool),
	}
}
