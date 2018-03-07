package src

type DBConfig struct {
	// 配置Key
	ConfigKey string `gorm:"column:ConfigKey;primary_key"`

	// 配置内容
	ConfigValue string `gorm:"column:ConfigValue"`
}

func (this *DBConfig) TableName() string {
	return "config"
}

type OtherConfig struct {
	// 最大历史消息保存数量
	MaxMessageLogCount int
}
