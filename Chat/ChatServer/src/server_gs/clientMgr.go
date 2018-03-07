package server_gs

import (
	"sync"
)

var (
	// 客户端连接列表
	clientMap = make(map[*client]bool, 32)

	// 读写锁
	clientMutex sync.RWMutex
)

// 注册客户端对象
// clientObj：客户端对象
// 返回值：
// 无
func registerClient(clientObj *client) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	clientMap[clientObj] = true
}

// 取消客户端注册
// clientObj：客户端对象
// 返回值：
// 无
func unregisterClient(clientObj *client) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	delete(clientMap, clientObj)
}

// 获取客户端的数量
// 返回值：
// 客户端数量
func getClientCount() int {
	clientMutex.RLock()
	defer clientMutex.RUnlock()

	return len(clientMap)
}

// 获取客户端列表
// 返回值：
// 客户端列表
func getClientList() (clientList []*client) {
	clientMutex.RLock()
	defer clientMutex.RUnlock()

	for k, _ := range clientMap {
		clientList = append(clientList, k)
	}

	return
}

func getExpiredClientList() (expiredList []*client) {
	clientMutex.RLock()
	defer clientMutex.RUnlock()

	for k, _ := range clientMap {
		if k.expired() {
			expiredList = append(expiredList, k)
		}
	}

	return
}
