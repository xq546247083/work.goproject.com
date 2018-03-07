package partner

// 合作商
type Partner struct {
	// 合作商Id
	Id int32 `json:"PartnerID"`

	// 合作商名称
	Name string `json:"PartnerName"`

	// 合作商别名
	Alias string `json:"PartnerAlias"`

	// 应用Id
	AppId string `json:"AppID"`

	// 登陆加密Key
	LoginKey string `json:"LoginKey"`

	// 充值配置
	ChargeConfig string `json:"ChargeConfig"`

	// 其它配置
	OtherConfigInfo string `json:"OtherConfigInfo"`

	// 游戏版本下载Url
	GameVersionUrl string `json:"GameVersionUrl"`

	// 充值服务器Url
	ChargeServerUrl string `json:"ChargeServerUrl"`

	// 合作商类型
	PartnerType int32 `json:"PartnerType"`

	// 权重
	Weight int32 `json:"Weight"`
}
