package clientMgr

import (
	"sync"

	. "work.goproject.com/Chat/ChatServerModel/src"
)

var (
	clientMap = make(map[int32]IClient, 1024)
	mutex     sync.RWMutex
)

func RegisterClient(clientObj IClient) {
	mutex.Lock()
	defer mutex.Unlock()

	clientMap[clientObj.GetId()] = clientObj
}

func UnregisterClient(clientObj IClient) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(clientMap, clientObj.GetId())
}

func GetClient(id int32) (clientObj IClient, exists bool) {
	mutex.Lock()
	defer mutex.Unlock()
	
	clientObj, exists = clientMap[id]
	return
}

func getClientCount() int {
	mutex.RLock()
	defer mutex.RUnlock()

	return len(clientMap)
}

func getExpiredClientList() (expiredList []IClient) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, item := range clientMap {
		if item.Expired() {
			expiredList = append(expiredList, item)
		}
	}

	return
}

func ResponseResult(clientObj IClient, responseObj *ServerResponseObject) {
	clientObj.AppendSendData(responseObj)
}
