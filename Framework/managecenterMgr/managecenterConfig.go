package managecenterMgr

import (
	"encoding/json"
)

// ManageCenter配置对象
type ManageCenterConfig struct {
	// API地址
	ManageCenterAPIUrl string

	// 服务器组类型
	GroupType string

	// 指定的IP
	SpecifiedIp string

	// 指定的服务器组Id
	SpecifiedGroupId int32

	// 是否只获取已经开服的服务器组
	IsGroupOpen bool

	// 刷新间隔（单位：分钟）
	RefreshInterval int

	// 返回的结果是否被压缩
	IsResultCompressed bool
}

func (this *ManageCenterConfig) String() string {
	bytes, _ := json.Marshal(this)
	return string(bytes)
}

func NewManageCenterConfig(manageCenterAPIUrl, groupType, specifiedIp string,
	specifiedGroupId int32, isGroupOpen bool, refreshInterval int) *ManageCenterConfig {

	return &ManageCenterConfig{
		ManageCenterAPIUrl: manageCenterAPIUrl,
		GroupType:          groupType,
		SpecifiedIp:        specifiedIp,
		SpecifiedGroupId:   specifiedGroupId,
		IsGroupOpen:        isGroupOpen,
		RefreshInterval:    refreshInterval,
		IsResultCompressed: false,
	}
}

func NewManageCenterConfig2(manageCenterAPIUrl, groupType, specifiedIp string,
	specifiedGroupId int32, isGroupOpen bool, refreshInterval int, isResultCompressed bool) *ManageCenterConfig {

	return &ManageCenterConfig{
		ManageCenterAPIUrl: manageCenterAPIUrl,
		GroupType:          groupType,
		SpecifiedIp:        specifiedIp,
		SpecifiedGroupId:   specifiedGroupId,
		IsGroupOpen:        isGroupOpen,
		RefreshInterval:    refreshInterval,
		IsResultCompressed: isResultCompressed,
	}
}
