package managecenterMgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	. "work.goproject.com/Framework/managecenterModel/partner"
	. "work.goproject.com/Framework/managecenterModel/returnObject"
	. "work.goproject.com/Framework/managecenterModel/server"
	. "work.goproject.com/Framework/managecenterModel/serverGroup"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
	"work.goproject.com/goutil/zlibUtil"
)

var (
	// 服务器组集合
	serverGroupMap   = make(map[int32]*ServerGroup, 512)
	serverGroupMutex sync.RWMutex

	// 服务器组变化方法集合 (完整列表，新增列表，删除列表，更新列表)
	serverGroupChangeFuncMap = make(map[string]func([]*ServerGroup, []*ServerGroup, []*ServerGroup, []*ServerGroup))
)

// 重新加载服务器组
func reloadServerGroup() error {
	logUtil.DebugLog("开始刷新服务器组信息")

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["GroupType"] = managecenterConfig.GroupType
	postDict["GroupIP"] = managecenterConfig.SpecifiedIp
	postDict["GroupId"] = strconv.Itoa(int(managecenterConfig.SpecifiedGroupId))
	postDict["IsGroupOpen"] = strconv.FormatBool(managecenterConfig.IsGroupOpen)
	postDict["IsResultCompressed"] = strconv.FormatBool(managecenterConfig.IsResultCompressed)

	// 连接服务器，以获取数据
	url := getManageCenterUrl("ServerGroupList.ashx")
	returnBytes, err := webUtil.PostWebData(url, postDict, nil)
	if err != nil {
		logUtil.ErrorLog("获取服务器组列表出错，url:%s,错误信息为：%s, str:%s", url, err, string(returnBytes))
		return err
	}

	// 先进行解压缩
	if managecenterConfig.IsResultCompressed {
		returnBytes, err = zlibUtil.Decompress(returnBytes)
		if err != nil {
			logUtil.ErrorLog("zlib解压缩服务器组列表错误，错误信息为：%s", err)
			return err
		}
	}

	// 解析返回值
	returnObj := new(ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		logUtil.ErrorLog("获取服务器组列表出错，反序列化返回值出错，错误信息为：%s", err)
		return err
	}

	// 判断返回状态是否为成功
	if returnObj.Code != 0 {
		msg := fmt.Sprintf("获取服务器组列表出错，返回状态：%d，信息为：%s", returnObj.Code, returnObj.Message)
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	}

	// 解析Data
	tmpServerGroupList := make([]*ServerGroup, 0, 512)
	if data, ok := returnObj.Data.(string); !ok {
		msg := "获取服务器组列表出错，返回的数据不是string类型"
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	} else {
		if err = json.Unmarshal([]byte(data), &tmpServerGroupList); err != nil {
			logUtil.ErrorLog("获取服务器组列表出错，反序列化数据出错，错误信息为：%s", err)
			return err
		}
	}

	logUtil.DebugLog("刷新服务器组信息结束,服务器组数量:%d", len(tmpServerGroupList))

	tmpServerGroupMap := make(map[int32]*ServerGroup, 512)
	for _, item := range tmpServerGroupList {
		tmpServerGroupMap[item.Id] = item
		// fmt.Println(item.GetGSCallbackUrl("API/LvRank"))
	}

	// 判断服务器组是否有变化，如果有变化
	if addList, deleteList, updateList, changed := isServerGroupChanged(tmpServerGroupMap); changed {
		// 则触发服务器组变化的方法
		triggerServerGroupChangeFunc(tmpServerGroupList, addList, deleteList, updateList)
	}

	// 初始化ip信息
	initIpData(tmpServerGroupMap)

	// 赋值给最终的ServerGroupMap
	serverGroupMutex.Lock()
	defer serverGroupMutex.Unlock()
	serverGroupMap = tmpServerGroupMap

	return nil
}

// 判断服务器组是否有变化
// tmpServerGroupMap：临时服务器组集合
func isServerGroupChanged(tmpServerGroupMap map[int32]*ServerGroup) (
	addList []*ServerGroup, deleteList []*ServerGroup, updateList []*ServerGroup,
	changed bool) {

	serverGroupMutex.RLock()
	defer serverGroupMutex.RUnlock()

	// 判断是否有新增的数据
	for k, v := range tmpServerGroupMap {
		if _, exists := serverGroupMap[k]; !exists {
			addList = append(addList, v)
			changed = true
		}
	}

	// 判断是否有删除的数据
	for k, v := range serverGroupMap {
		if _, exists := tmpServerGroupMap[k]; !exists {
			deleteList = append(deleteList, v)
			changed = true
		}
	}

	// 判断是否有更新的数据
	for k, oldItem := range serverGroupMap {
		if newItem, exists := tmpServerGroupMap[k]; exists {
			if oldItem.IsEqual(newItem) == false {
				updateList = append(updateList, newItem)
				changed = true
			}
		}
	}

	if changed {
		logUtil.DebugLog("addList count:%d, as follows:", len(addList))
		if debugUtil.IsDebug() {
			for _, item := range addList {
				logUtil.DebugLog("%v", item)
			}
		}
		logUtil.DebugLog("deleteList count:%d, as follows:", len(deleteList))
		if debugUtil.IsDebug() {
			for _, item := range deleteList {
				logUtil.DebugLog("%v", item)
			}
		}
		logUtil.DebugLog("updateList count:%d, as follows:", len(updateList))
		if debugUtil.IsDebug() {
			for _, item := range updateList {
				logUtil.DebugLog("%v", item)
			}
		}
	}

	return
}

// 触发服务器组变化的方法
func triggerServerGroupChangeFunc(allList []*ServerGroup, addList []*ServerGroup, deleteList []*ServerGroup, updateList []*ServerGroup) {
	// 如果有注册服务器组变化的方法
	if len(serverGroupChangeFuncMap) > 0 {
		for funcName, serverGroupChangeFunc := range serverGroupChangeFuncMap {
			logUtil.DebugLog("开始触发服务器组变化的方法：%s", funcName)
			serverGroupChangeFunc(allList, addList, deleteList, updateList)
			logUtil.DebugLog("触发服务器组变化的方法：%s结束", funcName)
		}
	}
}

// 注册服务器组变化方法
// funcName：方法名称
// serverGroupChangeFunc：服务器组变化方法
func RegisterServerGroupChangeFunc(funcName string, serverGroupChangeFunc func([]*ServerGroup, []*ServerGroup, []*ServerGroup, []*ServerGroup)) {
	if _, exists := serverGroupChangeFuncMap[funcName]; exists {
		panic(fmt.Errorf("RegisterServerGroupChange:%s已经存在，请检查", funcName))
	}
	serverGroupChangeFuncMap[funcName] = serverGroupChangeFunc

	logUtil.DebugLog("注册服务器组变化方法 funcName:%s，当前共有%d个注册", funcName, len(serverGroupChangeFuncMap))
}

// 获取服务器组集合
// 返回值：
// 服务器组集合
func GetServerGroupMap() (retServerGroupMap map[int32]*ServerGroup) {
	serverGroupMutex.RLock()
	defer serverGroupMutex.RUnlock()

	retServerGroupMap = make(map[int32]*ServerGroup, 128)
	for k, v := range serverGroupMap {
		retServerGroupMap[k] = v
	}

	return
}

// 获取服务器组列表
// 返回值：
// 服务器组列表
func GetServerGroupList() (serverGroupList []*ServerGroup) {
	serverGroupMutex.RLock()
	defer serverGroupMutex.RUnlock()

	for _, item := range serverGroupMap {
		serverGroupList = append(serverGroupList, item)
	}

	sort.Slice(serverGroupList, func(i, j int) bool {
		return serverGroupList[i].SortByIdAsc(serverGroupList[j])
	})

	return
}

// 获取服务器组项
// id
// 返回值：
// 服务器组对象
// 是否存在
func GetServerGroupItem(id int32) (serverGroupObj *ServerGroup, exists bool) {
	serverGroupMutex.RLock()
	defer serverGroupMutex.RUnlock()

	serverGroupObj, exists = serverGroupMap[id]
	return
}

// 根据合作商Id、服务器Id获取服务器组对象
// partnerId：合作商Id
// serverId：服务器Id
// 返回值：
// 服务器组对象
// 服务器对象
// 是否存在
func GetServerGroup(partnerId, serverId int32) (serverGroupObj *ServerGroup, serverObj *Server, exists bool) {
	var partnerObj *Partner

	// 获取合作商对象
	partnerObj, exists = GetPartner(partnerId)
	if !exists {
		return
	}

	// 获取服务器对象
	serverObj, exists = GetServer(partnerObj, serverId)
	if !exists {
		return
	}

	// 获取服务器组对象
	serverGroupObj, exists = GetServerGroupItem(serverObj.GroupId)
	return
}
