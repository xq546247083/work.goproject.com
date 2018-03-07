package managecenterMgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	. "work.goproject.com/Framework/managecenterModel/partner"
	. "work.goproject.com/Framework/managecenterModel/returnObject"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
	"work.goproject.com/goutil/zlibUtil"
)

var (
	partnerMap   = make(map[int32]*Partner, 128)
	partnerMutex sync.RWMutex
)

// 重新加载合作商
func reloadPartner() error {
	logUtil.DebugLog("开始刷新合作商列表")

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["IsResultCompressed"] = strconv.FormatBool(managecenterConfig.IsResultCompressed)

	// 连接服务器，以获取数据
	url := getManageCenterUrl("PartnerList.ashx")
	returnBytes, err := webUtil.PostWebData(url, postDict, nil)
	if err != nil {
		logUtil.ErrorLog("获取合作商列表出错，url:%s,错误信息为:%s", url, err)
		return err
	}

	// 先进行解压缩
	if managecenterConfig.IsResultCompressed {
		returnBytes, err = zlibUtil.Decompress(returnBytes)
		if err != nil {
			logUtil.ErrorLog("zlib解压缩合作商列表错误，错误信息为：%s", err)
			return err
		}
	}

	// 解析返回值
	returnObj := new(ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		logUtil.ErrorLog("获取合作商列表出错，反序列化返回值出错，错误信息为：%s, str:%s", err, string(returnBytes))
		return err
	}

	// 判断返回状态是否为成功
	if returnObj.Code != 0 {
		msg := fmt.Sprintf("获取合作商列表出错，返回状态：%d，信息为：%s", returnObj.Code, returnObj.Message)
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	}

	// 解析Data
	tmpPartnerList := make([]*Partner, 0, 128)
	if data, ok := returnObj.Data.(string); !ok {
		msg := "获取合作商列表出错，返回的数据不是string类型"
		logUtil.ErrorLog(msg)
		return errors.New(msg)
	} else {
		if err = json.Unmarshal([]byte(data), &tmpPartnerList); err != nil {
			logUtil.ErrorLog("获取合作商列表出错，反序列化数据出错，错误信息为：%s", err)
			return err
		}
	}

	logUtil.DebugLog("刷新合作商信息结束，合作商数量:%d", len(tmpPartnerList))

	tmpPartnerMap := make(map[int32]*Partner)
	for _, item := range tmpPartnerList {
		tmpPartnerMap[item.Id] = item
	}

	// 赋值给最终的partnerMap
	partnerMutex.Lock()
	defer partnerMutex.Unlock()
	partnerMap = tmpPartnerMap

	return nil
}

// 根据合作商Id获取合作商对象
// id：合作商Id
// 返回值：
// 合作商对象
// 是否存在
func GetPartner(id int32) (partnerObj *Partner, exists bool) {
	partnerMutex.RLock()
	defer partnerMutex.RUnlock()

	partnerObj, exists = partnerMap[id]
	return
}

// 获取合作商的其它配置信息
// id:合作商Id
// configKey:其它配置Key
// 返回值
// 配置内容
// 是否存在
// 错误对象
func GetOtherConfigInfo(id int32, configKey string) (configValue string, exists bool, err error) {
	partnerObj, exists := GetPartner(id)
	if !exists {
		return
	}

	if partnerObj.OtherConfigInfo == "" {
		return
	}

	var otherConfigMap map[string]string
	if err = json.Unmarshal([]byte(partnerObj.OtherConfigInfo), &otherConfigMap); err != nil {
		return
	}

	configValue, exists = otherConfigMap[configKey]
	return
}
