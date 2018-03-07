package gameServerUtil

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"work.goproject.com/Framework/managecenterMgr"
	. "work.goproject.com/Framework/managecenterModel/serverGroup"
	"work.goproject.com/goutil/timeUtil"
	"work.goproject.com/goutil/typeUtil"
)

type ManageCenterUtil struct{}

// 检测是否有游戏版本更新
// partnerId：合作商Id
// serverId：服务器Id
// gameVersionId：游戏版本号
// 返回值
// 游戏版本地址
// 是否有游戏版本更新
func (this *ManageCenterUtil) CheckNewGameVersion(partnerId, serverId, gameVersionId int32) (gameVersionUrl string, exists bool) {
	// 获取服务器对象
	serverObj, exists1 := managecenterMgr.GetServerItem(partnerId, serverId)
	if !exists1 {
		return
	}

	// 判断版本是否有更新
	if gameVersionId < serverObj.MinGameVersionId {
		partnerObj, exists1 := managecenterMgr.GetPartner(partnerId)
		if !exists1 {
			return
		}

		gameVersionUrl = partnerObj.GameVersionUrl
		exists = true
	}

	return
}

// 检测是否有资源版本更新
// partnerId：合作商Id
// serverId：服务器Id
// gameVersionId：游戏版本号
// resourceVersionName：资源版本名称
// 返回值
// 资源版本名称
// 资源版本地址
// 是否有游戏版本更新
func (this *ManageCenterUtil) CheckNewResourceVersion(partnerId, serverId, gameVersionId int32, resourceVersionName string) (newName string, url string, exists bool) {
	// 获取服务器对象
	serverGroupObj, _, exists := managecenterMgr.GetServerGroup(partnerId, serverId)
	if !exists {
		return
	}

	newName, url, exists = managecenterMgr.GetAvailableResource(partnerId, gameVersionId, resourceVersionName, OfficialOrTest(serverGroupObj.OfficialOrTest))
	return
}

// 检查服务器是否在维护
// serverGroupId：服务器组Id
// 返回值
// 维护消息
// 服务器是否在维护
func (this *ManageCenterUtil) CheckMaintainStatus(serverGroupId int32) (maintainMessage string, isMaintaining bool) {
	serverGroupObj, exists := managecenterMgr.GetServerGroupItem(serverGroupId)
	if !exists {
		return
	}

	// 判断维护状态
	nowTick := time.Now().Unix()
	if serverGroupObj.GroupState == int32(Con_GroupState_Maintain) || (serverGroupObj.MaintainBeginTimeTick <= nowTick && nowTick <= serverGroupObj.MaintainBeginTimeTick+int64(60*serverGroupObj.MaintainMinutes)) {
		maintainMessage = serverGroupObj.MaintainMessage
		isMaintaining = true
		return
	}

	return
}

// 获取服务器组的开服信息
// serverGroupId：服务器组Id
// 返回值
// 服务器开服日期
// 服务器开服天数
func (this *ManageCenterUtil) GetServerOpenDateInfo(serverGroupId int32) (openDate time.Time, openDays int) {
	serverGroupObj, exists := managecenterMgr.GetServerGroupItem(serverGroupId)
	if !exists {
		return
	}

	openDate, _ = typeUtil.DateTime(serverGroupObj.OpenTimeTick)
	openDays = timeUtil.SubDay(time.Now(), openDate) + 1

	return
}

// 获取合作商、服务器的组合字符串
// serverGroupId：服务器组Id
// 返回值
// 合作商、服务器的组合字符串
func (this *ManageCenterUtil) GetPartnerServerPairString(serverGroupId int32) string {
	var buf bytes.Buffer

	serverList := managecenterMgr.GetServerList(serverGroupId)
	for _, item := range serverList {
		buf.WriteString(fmt.Sprintf("%d_%d|", item.PartnerId, item.Id))
	}

	return buf.String()
}

// 是否是有效的合作商、服务器组合
// partnerId：合作商Id
// serverId：服务器Id
// parnterServerPairString：合作商、服务器的组合字符串,格式为:PartnerId_ServerId|PartnerId_ServerId|
// 返回值
// 是否有效
func (this *ManageCenterUtil) IfValidPartnerServerPair(partnerId, serverId int32, parnterServerPairString string) bool {
	if parnterServerPairString == "" {
		return false
	}

	partnerServerPairStringList := strings.Split(parnterServerPairString, "|")
	if len(partnerServerPairStringList) == 0 {
		return false
	}

	// 获得玩家的合作商、服务器组合字符串
	partnerServerPair := fmt.Sprintf("%d_%d", partnerId, serverId)

	// 遍历寻找
	for _, item := range partnerServerPairStringList {
		if item == partnerServerPair {
			return true
		}
	}

	return false
}

// ------------------类型定义和业务逻辑的分隔符-------------------------

var (
	ManageCenterUtilObj = new(ManageCenterUtil)
)
