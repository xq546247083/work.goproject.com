package word

type ForbidWord struct {
	Word string `gorm:"column:Word"`
}

func (this *ForbidWord) TableName() string {
	return "config_word_forbid"
}

type SensitiveWord struct {
	Word string `gorm:"column:Word"`
}

func (this *SensitiveWord) TableName() string {
	return "config_word_sensitive"
}
