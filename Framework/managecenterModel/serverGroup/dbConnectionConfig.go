package serverGroup

import (
	. "work.goproject.com/goutil/mysqlUtil"
	. "work.goproject.com/goutil/redisUtil"
)

// 数据库连接字符串配置
type DBConnectionConfig struct {
	// 模型数据库内网连接字符串
	GameModelDB string

	// 游戏数据库内网连接字符串
	GameDB string

	// 日志数据库内网连接字符串
	LogDB string

	// Redis连接字符串
	RedisConfig string
}

// 获取游戏模型数据库连接
// 返回值:
// 数据库连接配置对象
// 错误对象
func (this *DBConnectionConfig) GetGameModelDBConn() (dbConfig *DBConfig, err error) {
	dbConfig, err = NewDBConfig2(this.GameModelDB)
	return
}

// 获取游戏数据库连接
// 返回值:
// 数据库连接配置对象
// 错误对象
func (this *DBConnectionConfig) GetGameDBConn() (dbConfig *DBConfig, err error) {
	dbConfig, err = NewDBConfig2(this.GameDB)
	return
}

// 获取游戏日志数据库连接
// 返回值:
// 数据库连接配置对象
// 错误对象
func (this *DBConnectionConfig) GetLogDBConn() (dbConfig *DBConfig, err error) {
	dbConfig, err = NewDBConfig2(this.LogDB)
	return
}

// 获取Redis配置
// 返回值:
// redis配置对象
// 错误对象
func (this *DBConnectionConfig) GetRedisConfig() (redisConfig *RedisConfig, err error) {
	redisConfig, err = NewRedisConfig(this.RedisConfig)
	return
}
