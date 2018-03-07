package managecenterMgr

import (
	"sync"

	. "work.goproject.com/Framework/managecenterModel/serverGroup"
)

var (
	ipMap   = make(map[string]bool, 1024)
	ipMutex sync.RWMutex
)

func initIpData(serverGroupMap map[int32]*ServerGroup) {
	ipMutex.Lock()
	defer ipMutex.Unlock()

	ipMap = make(map[string]bool, 1024)
	for _, item := range serverGroupMap {
		for _, ipItem := range item.GetIPList() {
			ipMap[ipItem] = true
		}
	}
}

// 判断IP是否有效
// ip：指定IP地址
// 返回值：
// IP是否有效
func IsIpValid(ip string) bool {
	ipMutex.RLock()
	defer ipMutex.RUnlock()

	_, exists := ipMap[ip]

	return exists
}
