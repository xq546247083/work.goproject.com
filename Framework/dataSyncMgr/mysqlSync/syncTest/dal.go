package main

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"work.goproject.com/Framework/dataSyncMgr/mysqlSync"
	"work.goproject.com/goutil/logUtil"
)

var _ = mysql.DeregisterLocalFile

var (
	connectionString = "root:moqikaka3309@tcp(10.1.0.10:3309)/develop_liujun?charset=utf8&parseTime=true&loc=Local&timeout=60s"
	maxOpenConns     = 10
	maxIdleConns     = 10

	syncFileSize = 1024 * 1024
)

var (
	// 数据库对象
	dbObj *gorm.DB

	// 同步管理对象
	syncMgr *mysqlSync.SyncMgr
)

func init() {
	// 初始化数据库连接
	dbObj = initMysql()

	// 构造同步管理对象
	syncMgr = mysqlSync.NewSyncMgr(1, "Sync", syncFileSize, 1, dbObj.DB())
}

// 初始化Mysql
func initMysql() *gorm.DB {
	dbObj, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic(fmt.Errorf("初始化数据库:%s失败，错误信息为：%s", connectionString, err))
	}
	logUtil.DebugLog("连接mysql:%s成功", connectionString)

	if maxOpenConns > 0 && maxIdleConns > 0 {
		dbObj.DB().SetMaxOpenConns(maxOpenConns)
		dbObj.DB().SetMaxIdleConns(maxIdleConns)
	}

	return dbObj
}

// 注册同步对象
func registerSyncObj(identifier string) {
	syncMgr.RegisterSyncObj(identifier)
}

// 保存sql数据
func save(identifier string, command string) {
	syncMgr.Save(identifier, command)
}
