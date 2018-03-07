package server

// 服务器
type Server struct {
	// 服务器Id
	Id int32 `json:"ServerID"`

	// 服务器名称
	Name string `json:"ServerName"`

	// 合作商Id
	PartnerId int32 `json:"PartnerID"`

	// 服务器组Id
	GroupId int32 `json:"GroupID"`

	// 对应的游戏版本号
	GameVersionId int32 `json:"GameVersionID"`

	// 需要的最低游戏版本号
	MinGameVersionId int32 `json:"MinGameVersionID"`
}
