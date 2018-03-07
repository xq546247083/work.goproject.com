package managecenterMgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	. "work.goproject.com/Framework/managecenterModel/resourceVersion"
	. "work.goproject.com/Framework/managecenterModel/returnObject"
	. "work.goproject.com/Framework/managecenterModel/serverGroup"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
	"work.goproject.com/goutil/zlibUtil"
)

var (
	resourceVersionList  = make([]*ResourceVersion, 128)
	resourceVersionMutex sync.RWMutex
)

// 重新加载资源列表
func reloadResourceVersion() error {
	logUtil.DebugLog("开始刷新资源列表")

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["IsResultCompressed"] = strconv.FormatBool(managecenterConfig.IsResultCompressed)

	// 连接服务器，以获取数据
	url := getManageCenterUrl("ResourceVersionList.ashx")
	returnBytes, err := webUtil.PostWebData(url, postDict, nil)
	if err != nil {
		logUtil.ErrorLog("获取资源列表出错，url:%s,错误信息为:%s", url, err)
		return err
	}

	// 先进行解压缩
	if managecenterConfig.IsResultCompressed {
		returnBytes, err = zlibUtil.Decompress(returnBytes)
		if err != nil {
			logUtil.ErrorLog("zlib解压缩资源列表错误，错误信息为：%s", err)
			return err
		}
	}

	// 解析返回值
	returnObj := new(ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		logUtil.ErrorLog("获取资源列表出错，反序列化返回值出错，错误信息为：%s, str:%s", err, string(returnBytes))
		return err
	}

	// 判断返回状态是否为成功
	if returnObj.Code != 0 {
		msg := fmt.Sprintf("获取资源列表出错，返回状态：%d，信息为：%s", returnObj.Code, returnObj.Message)
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	}

	// 解析Data
	tmpResourceVersionList := make([]*ResourceVersion, 0, 128)
	if data, ok := returnObj.Data.(string); !ok {
		msg := "获取资源列表出错，返回的数据不是string类型"
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	} else {
		if err = json.Unmarshal([]byte(data), &tmpResourceVersionList); err != nil {
			logUtil.ErrorLog("获取资源列表出错，反序列化数据出错，错误信息为：%s", err)
			return err
		}
	}

	logUtil.DebugLog("刷新资源信息结束，资源数量:%d", len(tmpResourceVersionList))

	// 赋值给最终的partnerMap
	resourceVersionMutex.Lock()
	defer resourceVersionMutex.Unlock()
	resourceVersionList = tmpResourceVersionList

	return nil
}

// 判断资源版本名称是否有效
// name：资源版本名称
// 返回值
// 时间戳
// HashCode
// 是否有效
func isNameValid(name string) (timeTick int64, hashCode string, valid bool) {
	if name == "" {
		return
	}

	// 判断resourceVersionName的是否包含_
	if strings.Index(name, "_") < 0 {
		logUtil.ErrorLog("IsNameValid:%s Format Error.", name)
		return
	}

	// 分割resourceVersionName
	itemList := strings.Split(name, "_")
	if len(itemList) != 2 {
		logUtil.ErrorLog("IsNameValid:%s Format Error.", name)
		return
	}

	timeTick, err := strconv.ParseInt(itemList[0], 10, 64)
	if err != nil {
		logUtil.ErrorLog("IsNameValid TimeTick: %s Format Error.", itemList[0])
		return
	}

	hashCode = itemList[1]
	valid = true

	return
}

// 从传入的资源列表中筛选出可用的资源
// partnerId：合作商Id
// gameVersionId：游戏版本Id
// resourceVersionName：资源版本名称
// 返回值
// 时间戳
// HashCode
// 是否有效
func GetAvailableResource(partnerId, gameVersionId int32, resourceVersionName string, officialOrTest OfficialOrTest) (newName string, url string, exists bool) {
	resourceVersionMutex.RLock()
	defer resourceVersionMutex.RUnlock()

	if len(resourceVersionList) == 0 {
		return
	}

	// 判断资源版本名称是否有效
	timeTick, hashCode, valid := isNameValid(resourceVersionName)
	if !valid {
		return
	}

	// 根据合作商Id、游戏版本Id和开始时间来过滤数据
	validList := make([]*ResourceVersion, 0, 32)
	nowTick := time.Now().Unix()
	for _, item := range resourceVersionList {
		if item.ContainsPartner(partnerId) && item.ContainsGameVersion(gameVersionId) && item.StartTimeTick <= nowTick {
			validList = append(validList, item)
		}
	}

	if len(validList) == 0 {
		return
	}

	// 按照资源Id进行降序排序
	sort.Slice(validList, func(i, j int) bool {
		return validList[i].SortByIdDesc(validList[j])
	})

	// 取出资源号最大的资源，如果与传入的资源名称相等，则表示没有新资源
	newestResource := validList[0]
	if newestResource.Name == resourceVersionName {
		return
	}

	// 判断资源号中的HashCode是否相等，如果相等，则表示没有新资源；如果传入的timeTick>最新的timeTick说明服务器没有被刷新，表示没有新资源
	newestTimeTick, newestHashCode, newestValid := isNameValid(newestResource.Name)
	if !newestValid {
		return
	}

	if hashCode == newestHashCode || timeTick > newestTimeTick {
		return
	}

	// 如果是测试服，且指定不下载资源，则表示没有新资源
	if officialOrTest == Con_Test && newestResource.IfAuditServiceDownload == 0 {
		return
	}

	// 返回数据
	newName = newestResource.Name
	url = newestResource.Url
	exists = true

	return
}
