package managecenterMgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	. "work.goproject.com/Framework/managecenterModel/partner"
	. "work.goproject.com/Framework/managecenterModel/returnObject"
	. "work.goproject.com/Framework/managecenterModel/server"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
	"work.goproject.com/goutil/zlibUtil"
)

var (
	serverMap         = make(map[int32]map[int32]*Server, 128)
	serverDistinctMap = make(map[int32]*Server, 1024)
	serverMutex       sync.RWMutex
)

// 重新加载服务器
func reloadServer() error {
	logUtil.NormalLog("开始刷新服务器列表", logUtil.Debug)

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["GroupType"] = managecenterConfig.GroupType
	postDict["GroupId"] = strconv.Itoa(int(managecenterConfig.SpecifiedGroupId))
	postDict["IsResultCompressed"] = strconv.FormatBool(managecenterConfig.IsResultCompressed)

	// 连接服务器，以获取数据
	url := getManageCenterUrl("ServerList.ashx")
	returnBytes, err := webUtil.PostWebData(url, postDict, nil)
	if err != nil {
		logUtil.ErrorLog("获取服务器列表出错，url:%s,错误信息为:%s", url, err)
		return err
	}

	// 先进行解压缩
	if managecenterConfig.IsResultCompressed {
		returnBytes, err = zlibUtil.Decompress(returnBytes)
		if err != nil {
			logUtil.ErrorLog("zlib解压缩服务器列表错误，错误信息为：%s", err)
			return err
		}
	}

	// 解析返回值
	returnObj := new(ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		logUtil.ErrorLog("获取服务器列表出错，反序列化返回值出错，错误信息为：%s, str:%s", err, string(returnBytes))
		return err
	}

	// 判断返回状态是否成功
	if returnObj.Code != 0 {
		msg := fmt.Sprintf("获取服务器列表出错，返回状态：%d，信息为：%s", returnObj.Code, returnObj.Message)
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	}

	// 解析Data
	tmpServerList := make([]*Server, 0, 1024)
	if data, ok := returnObj.Data.(string); !ok {
		msg := "获取服务器列表出错，返回的数据不是string类型"
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	} else {
		if err = json.Unmarshal([]byte(data), &tmpServerList); err != nil {
			logUtil.ErrorLog("获取服务器列表出错，反序列化数据出错，错误信息为：%s", err)
			return err
		}
	}

	logUtil.DebugLog("刷新服务器信息结束，服务器数量:%d", len(tmpServerList))

	tmpServerMap := make(map[int32]map[int32]*Server, 128)
	tmpServerDistinctMap := make(map[int32]*Server, 1024)
	for _, item := range tmpServerList {
		// 构造tmpServerMap数据
		if _, ok := tmpServerMap[item.PartnerId]; !ok {
			tmpServerMap[item.PartnerId] = make(map[int32]*Server, 1024)
		}
		tmpServerMap[item.PartnerId][item.Id] = item

		// 构造tmpServerDistinctMap数据
		tmpServerDistinctMap[item.Id] = item
	}

	// 赋值给最终的serverMap、serverDistinctMap
	serverMutex.Lock()
	defer serverMutex.Unlock()

	serverMap = tmpServerMap
	serverDistinctMap = tmpServerDistinctMap

	return nil
}

// 获取服务器组对应的所有服务器列表
// 返回值
// 服务器列表
func GetServerList(serverGroupId int32) (serverList []*Server) {
	serverMutex.RLock()
	defer serverMutex.RUnlock()

	for _, subMap := range serverMap {
		for _, item := range subMap {
			if item.GroupId == serverGroupId {
				serverList = append(serverList, item)
			}
		}
	}

	return
}

// 根据合作商对象、服务器Id获取服务器对象
// partnerObj：合作商对象
// serverId：服务器Id
// 返回值：
// 服务器对象
// 是否存在
func GetServer(partnerObj *Partner, serverId int32) (serverObj *Server, exists bool) {
	serverMutex.RLock()
	defer serverMutex.RUnlock()

	if subServerMap, exists1 := serverMap[partnerObj.Id]; exists1 {
		serverObj, exists = subServerMap[serverId]
	}

	return
}

// 根据合作商Id、服务器Id获取服务器对象
// partnerId：合作商Id
// serverId：服务器Id
// 返回值：
// 服务器对象
// 是否存在
func GetServerItem(partnerId, serverId int32) (serverObj *Server, exists bool) {
	serverMutex.RLock()
	defer serverMutex.RUnlock()

	if subServerMap, exists1 := serverMap[partnerId]; exists1 {
		serverObj, exists = subServerMap[serverId]
	}

	return
}

// 根据服务器组Id获取对应的服务器列表
// groupId:服务器组Id
// 返回值:
// 服务器列表
func GetServerListByGroupId(groupId int32) (serverList []*Server) {
	serverMutex.RLock()
	defer serverMutex.RUnlock()

	for _, item := range serverDistinctMap {
		if item.GroupId == groupId {
			serverList = append(serverList, item)
		}
	}

	return
}

// 获取不重复的服务器Id列表
// 返回值:
// 不重复的服务器Id列表
func GetDistinctServerIdList() (distinctServerIdList []int32) {
	serverMutex.RLock()
	defer serverMutex.RUnlock()

	for _, item := range serverDistinctMap {
		distinctServerIdList = append(distinctServerIdList, item.Id)
	}

	return
}
