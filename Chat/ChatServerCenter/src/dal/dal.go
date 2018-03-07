package dal

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"work.goproject.com/Chat/ChatServerCenter/src/config"
	"work.goproject.com/Framework/monitorMgr"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/mysqlUtil"
)

var (
	// 数据库对象
	dbObj *gorm.DB
)

func init() {
	// 初始化数据库连接
	dbConfig := config.GetDBConfig()

	var err error
	logUtil.DebugLog("开始连接Mysql数据库")
	if dbObj, err = gorm.Open("mysql", dbConfig.ConnectionString); err != nil {
		panic(fmt.Errorf("初始化数据库失败，错误信息为：%s", err))
	}
	logUtil.DebugLog("连接Mysql数据库成功")

	if dbConfig.MaxOpenConns > 0 && dbConfig.MaxIdleConns > 0 {
		dbObj.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
		dbObj.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
	}

	// 注册监控方法
	monitorMgr.RegisterMonitorFunc(monitor)
}

// 获取数据库对象
// 返回值：
// 数据库对象
func GetDB() *gorm.DB {
	return dbObj
}

// 读取数据表的所有数据
// dataList:用户保存数据的
func GetAll(dataList interface{}) error {
	if result := dbObj.Find(dataList); result.Error != nil {
		WriteLog("dal.GetAll", result.Error)
		return result.Error
	}

	return nil
}

// 记录日志
// funcName:方法名称
// err:错误对象
func WriteLog(funcName string, err error) {
	logUtil.ErrorLog("%s出错，错误信息：%s", funcName, err)
}

// 监控数据库的可用情况
func monitor() error {
	return mysqlUtil.TestConnection(dbObj.DB())
}
