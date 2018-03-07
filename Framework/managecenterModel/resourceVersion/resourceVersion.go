package resourceVersion

import (
	"work.goproject.com/goutil/stringUtil"
)

// 游戏版本
type ResourceVersion struct {
	// 资源版本唯一标识
	Id int32 `json:"ResourceVersionID"`

	// 资源版本名称
	Name string `json:"ResourceVersionName"`

	// 资源版本的url地址
	Url string `json:"ResourceVersionUrl"`

	// 资源大小
	Size int32 `json:"Size"`

	// 资源文件MD5加密的结果
	MD5 string `json:"MD5"`

	// 资源生效时间
	StartTime     string `json:"StartTime"`
	StartTimeTick int64  `json:"StartTimeTick"`

	// 资源失效时间
	EndTime     string `json:"EndTime"`
	EndTimeTick int64  `json:"EndTimeTick"`

	// 添加时间
	Crdate     string `json:"Crdate"`
	CrdateTick int64  `json:"CrdateTick"`

	// 更新时间
	UpdateTime     string `json:"UpdateTime"`
	UpdateTimeTick int64  `json:"UpdateTimeTick"`

	// 是否重启客户端
	IfRestart int32 `json:"IfRestart"`

	// 是否禁用
	IfDelete int32 `json:"IfDelete"`

	// 是否审核服下载
	IfAuditServiceDownload int32 `json:"IfAuditServiceDownload"`

	// 资源所属的合作商ID集合
	PartnerIds string `json:"PartnerIDs"`

	// 资源所属的游戏版本ID集合
	GameVersionIds string `json:"GameVersionIDs"`
}

// 判断资源是否包含指定合作商
// partnerId：合作商Id
// 返回值
// 是否包含
func (this *ResourceVersion) ContainsPartner(partnerId int32) bool {
	partnerIdList, _ := stringUtil.SplitToInt32Slice(this.PartnerIds, ",")
	for _, item := range partnerIdList {
		if item == partnerId {
			return true
		}
	}

	return false
}

// 判断资源是否包含指定游戏版本
// gameVersionId：游戏版本Id
// 返回值
// 是否包含
func (this *ResourceVersion) ContainsGameVersion(gameVersionId int32) bool {

	gameVersionIdList, _ := stringUtil.SplitToInt32Slice(this.GameVersionIds, ",")
	for _, item := range gameVersionIdList {
		if item == gameVersionId {
			return true
		}
	}

	return false
}

// 按照Id进行升序排序
// target:另一个资源对象
// 是否是小于
func (this *ResourceVersion) SortByIdAsc(target *ResourceVersion) bool {
	return this.Id < target.Id
}

// 按照Id进行降序排序
// target:另一个资源对象
// 是否是大于
func (this *ResourceVersion) SortByIdDesc(target *ResourceVersion) bool {
	return this.Id > target.Id
}
