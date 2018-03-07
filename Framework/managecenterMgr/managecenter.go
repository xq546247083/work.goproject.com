package managecenterMgr

import (
	"fmt"
	"strings"
	"time"

	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/Framework/initMgr"
	"work.goproject.com/goutil/logUtil"
)

var (
	managecenterConfig *ManageCenterConfig
	initSuccessObj     = initMgr.NewInitSuccess("managecenterMgr")
)

// 注册初始化成功的通道
// name:模块名称
// ch:通道对象
func RegisterInitSuccess(name string, ch chan bool) {
	initSuccessObj.Register(name, ch)
}

func setConfig(config *ManageCenterConfig) {
	logUtil.DebugLog("ManageCenterAPIUrl:%s", config.ManageCenterAPIUrl)
	logUtil.DebugLog("GroupType:%s", config.GroupType)
	logUtil.DebugLog("SpecifiedIp:%s", config.SpecifiedIp)
	logUtil.DebugLog("SpecifiedGroupId:%d", config.SpecifiedGroupId)
	logUtil.DebugLog("IsGroupOpen:%v", config.IsGroupOpen)
	logUtil.DebugLog("RefreshInterval:%d", config.RefreshInterval)
	logUtil.DebugLog("IsResultCompressed:%v", config.IsResultCompressed)

	managecenterConfig = config
}

// Start ...启动ManageCenter管理器(obsolete，建议使用Start2)
// manageCenterAPIUrl:ManageCenter对外提供的API
// groupType:服务器组类型
// specifiedIP:指定的IP，如果指定则只获取对应IP的服务器组（默认为""）
// specifiedGroupId:指定的服务器组Id，如果指定则只获取对应GroupId的服务器组（默认为0）
// isGroupOpen:是否只处理已经开服的
// refreshInterval:刷新时间间隔，单位：分钟
func Start(manageCenterAPIUrl string, groupType string, specifiedIp string, specifiedGroupId int32, isGroupOpen bool, refreshInterval int) {
	config := NewManageCenterConfig(manageCenterAPIUrl, groupType, specifiedIp, specifiedGroupId, isGroupOpen, refreshInterval)
	Start2(config)
}

// Start ...启动ManageCenter管理器
// config：ManagerCenter配置对象
func Start2(config *ManageCenterConfig) {
	setConfig(config)

	// 先初始化一次服务器组
	if err := reload(); err != nil {
		panic(err)
	}

	// 通知初始化成功
	initSuccessObj.Notify()

	// 如果config.RefreshInterval==0，表示不需要定时刷新；
	// 只有config.RefreshInterval>0，才需要定时更新
	if config.RefreshInterval > 0 {
		go timelyRefresh()
	}
}

// Refresh ...刷新ManageCenter信息
func Refresh() error {
	return reload()
}

// 定时刷新
func timelyRefresh() {
	goroutineName := "managecenterMgr.timelyRefresh"
	goroutineMgr.Monitor(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	for {
		// 每5分钟刷新一次（这样就不需要ManageCenter主动推送了）
		time.Sleep(time.Duration(managecenterConfig.RefreshInterval) * time.Minute)

		// 刷新服务器组
		reload()
	}
}

// 重新加载/初始化
func reload() error {
	var err error

	if err = reloadPartner(); err != nil {
		return err
	}

	if err = reloadServer(); err != nil {
		return err
	}

	if err = reloadServerGroup(); err != nil {
		return err
	}

	// 只有当不需要定时刷新才加载资源；
	// 因为当不需要定时刷新，表示进程是游戏服务器，而只有游戏服务器才需要处理资源，其它的诸如聊天、跨服、GameServerCenter等都不需要关注资源信息
	if managecenterConfig.RefreshInterval == 0 {
		if err = reloadResourceVersion(); err != nil {
			return err
		}
	}

	return nil
}

// 获取可访问的ManageCenter地址
// suffix:Url后缀
// 返回值:
// 可访问的ManageCenter地址
func getManageCenterUrl(suffix string) string {
	if strings.HasSuffix(managecenterConfig.ManageCenterAPIUrl, "/") {
		return fmt.Sprintf("%s%s", managecenterConfig.ManageCenterAPIUrl, suffix)
	} else {
		return fmt.Sprintf("%s/%s", managecenterConfig.ManageCenterAPIUrl, suffix)
	}
}
