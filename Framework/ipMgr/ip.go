/*
提供统一的ip验证的逻辑；包括两部分：
1、ManageCenter中配置到ServerGroup中的IP；系统内部自动处理，外部无需关注；
2、各个应用程序中配置的ip；通过调用方法来进行初始化；
*/
package ipMgr

import (
	"sync"

	"work.goproject.com/Framework/managecenterMgr"
	"work.goproject.com/goutil/stringUtil"
)

var (
	// ip集合
	ipMap = make(map[string]bool, 32)
	mutex sync.RWMutex
)

// 初始化IP列表(obsolete)
func Init(ipList []string) {
	mutex.Lock()
	defer mutex.Unlock()

	// 先清空再初始化
	ipMap = make(map[string]bool, 32)

	for _, item := range ipList {
		ipMap[item] = true
	}
}

// 初始化ip字符串（以分隔符分割的）
func InitString(ipStr string) {
	mutex.Lock()
	defer mutex.Unlock()

	// 先清空再初始化
	ipMap = make(map[string]bool, 32)

	for _, item := range stringUtil.Split(ipStr, nil) {
		ipMap[item] = true
	}
}

func isLocalIpValid(ip string) bool {
	mutex.RLock()
	defer mutex.RUnlock()

	_, exists := ipMap[ip]

	return exists
}

func isServerGroupIpValid(ip string) bool {
	return managecenterMgr.IsIpValid(ip)
}

// 判断传入的Ip是否有效
// ip:ip
// 返回值：
// 是否有效
func IsIpValid(ip string) bool {
	if isLocalIpValid(ip) {
		return true
	} else if isServerGroupIpValid(ip) {
		return true
	}

	return false
}
