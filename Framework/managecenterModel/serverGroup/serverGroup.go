package serverGroup

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"work.goproject.com/goutil/stringUtil"
)

// 服务器组
type ServerGroup struct {
	// 服务器组Id
	Id int32 `json:"GroupID"`

	// 服务器组名称
	Name string `json:"GroupName"`

	// 服务器组Url
	Url string `json:"GroupUrl"`

	// 聊天服务器Url
	ChatServerUrl string `json:"ChatServerUrl"`

	// 数据库连接配置
	DBConnectionConfig string `json:"DBConnectionConfig"`

	// 服务器组状态（1：正常；2：维护）
	GroupState int32 `json:"GroupState"`

	// 服务器组热度（1：正常；2：新服；3：推荐）
	GroupHeat int32 `json:"GroupHeat"`

	// 服务器组负载（1：正常；2：火爆）
	GroupLoad int32 `json:"GroupLoad"`

	// 服务器开服时间对应的Unix时间戳
	OpenTimeTick int64 `json:OpenTimeTick`

	// 服务器组Ip(外网IP;内网IP;回调GS内网端口)
	Ip string `json:"GroupIp"`

	// 正式服或测试服；1：正式服；2：测试服
	OfficialOrTest int32 `json:"OfficialOrTest"`

	// 服务器组类型
	Type int32 `json:"GroupType"`

	// 服务器组排序
	Order int32 `json:"GroupOrder"`

	// 服务器组维护开始时间对应的时间戳
	MaintainBeginTimeTick int64 `json:MaintainBeginTimeTick`

	// 维护持续分钟数
	MaintainMinutes int32 `json:"MaintainMinutes"`

	// 维护信息
	MaintainMessage string `json:"MaintainMessage"`

	// 游戏监听地址
	GameListenAddr string `json:"GameListenAddr"`

	// 回调监听地址
	CallbackListenAddr string `json:"CallbackListenAddr"`

	// 外网回调地址
	ExternalCallbackUrl string `json:"ExternalCallbackUrl"`

	// 内网回调地址
	InternalCallbackUrl string `json:"InternalCallbackUrl"`

	// 是否在主群组（机房）内
	IsInMainGroup bool `json:"IsInMainGroup"`

	// 监控端口
	GopsPort string `json:"GopsPort"`
}

// 排序方法(默认按照Id进行升序排序)
// target:另一个服务器组对象
// 是否是小于
func (this *ServerGroup) SortByIdAsc(target *ServerGroup) bool {
	return this.Id < target.Id
}

// 按照开服时间进行升序排序
// target:另一个服务器组对象
// 是否是小于
func (this *ServerGroup) SortByOpenTimeAsc(target *ServerGroup) bool {
	return this.OpenTimeTick < target.OpenTimeTick
}

// 获取数据库配置对象
// 返回值:
// 数据库配置对象
// 错误对象
func (this *ServerGroup) GetDBConfig() (*DBConnectionConfig, error) {
	var dbConfig *DBConnectionConfig
	if err := json.Unmarshal([]byte(this.DBConnectionConfig), &dbConfig); err != nil {
		return nil, err
	}

	return dbConfig, nil
}

// 获取ip列表
// 返回值：
// ip列表
func (this *ServerGroup) GetIPList() []string {
	return stringUtil.Split(this.Ip, nil)
}

// 服务器组是否开启
// 返回值：
// 是否开启
func (this *ServerGroup) IsOpen() bool {
	return this.OpenTimeTick < time.Now().Unix()
}

// 获取游戏服务器的回调地址
// suffix:地址后缀
// 返回值
// 游戏服务器的回调地址
func (this *ServerGroup) GetGSCallbackUrl(suffix string) string {
	// 如果是在主群组（机房）内，则使用内网地址，否则使用外网地址
	url := ""
	if this.IsInMainGroup {
		url = this.InternalCallbackUrl
	} else {
		url = this.ExternalCallbackUrl
	}

	if url != "" {
		if strings.HasSuffix(url, "/") {
			return fmt.Sprintf("%s%s", url, suffix)
		} else {
			return fmt.Sprintf("%s/%s", url, suffix)
		}
	}

	// 兼容旧的ManageCenter版本
	ipList := this.GetIPList()

	// 外网IP;内网IP;回调GS内网端口；如果数量小于3，则直接使用配置的GroupUrl；否则使用第3个值
	if len(ipList) < 3 {
		if strings.HasSuffix(this.Url, "/") {
			return fmt.Sprintf("%s%s", this.Url, suffix)
		} else {
			return fmt.Sprintf("%s/%s", this.Url, suffix)
		}
	} else {
		return fmt.Sprintf("http://%s:%s/%s", ipList[1], ipList[2], suffix)
	}
}

// 判断服务器组是否相同
// target:目标服务器组
// 是否相同
func (this *ServerGroup) IsEqual(target *ServerGroup) bool {
	return this.Id == target.Id &&
		this.Name == target.Name &&
		this.Url == target.Url &&
		this.ChatServerUrl == target.ChatServerUrl &&
		this.DBConnectionConfig == target.DBConnectionConfig &&
		this.GroupState == target.GroupState &&
		this.GroupHeat == target.GroupHeat &&
		this.GroupLoad == target.GroupLoad &&
		this.OpenTimeTick == target.OpenTimeTick &&
		this.Ip == target.Ip &&
		this.OfficialOrTest == target.OfficialOrTest &&
		this.Type == target.Type &&
		this.Order == target.Order &&
		this.MaintainBeginTimeTick == target.MaintainBeginTimeTick &&
		this.MaintainMinutes == target.MaintainMinutes &&
		this.MaintainMessage == target.MaintainMessage &&
		this.GameListenAddr == target.GameListenAddr &&
		this.CallbackListenAddr == target.CallbackListenAddr &&
		this.ExternalCallbackUrl == target.ExternalCallbackUrl &&
		this.InternalCallbackUrl == target.InternalCallbackUrl &&
		this.IsInMainGroup == target.IsInMainGroup &&
		this.GopsPort == target.GopsPort
}
