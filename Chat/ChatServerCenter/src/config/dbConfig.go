package config

import (
	"work.goproject.com/goutil/configUtil"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/mysqlUtil"
)

var (
	dbConfig *mysqlUtil.DBConfig
)

func initDBConfig(config *configUtil.XmlConfig) error {
	tempConfig := new(mysqlUtil.DBConfig)
	err := config.Unmarshal("root/DBConnection", tempConfig)
	if err != nil {
		return err
	}

	dbConfig = tempConfig
	debugUtil.Printf("dbConfig:%v\n", dbConfig)

	return nil
}

// GetDBConfig 获取mysql数据库配置
func GetDBConfig() *mysqlUtil.DBConfig {
	return dbConfig
}
