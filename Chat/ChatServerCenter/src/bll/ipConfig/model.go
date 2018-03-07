package ipConfig

type IpConfig struct {
	IP string `gorm:"column:IP"`
}

func (this *IpConfig) TableName() string {
	return "config_ip"
}
