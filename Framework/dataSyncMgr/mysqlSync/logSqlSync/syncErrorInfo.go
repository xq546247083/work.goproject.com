package logSqlSync

import (
	"database/sql"
	"time"

	"work.goproject.com/goutil/logUtil"
)

// 错误信息记录表是否已经初始化
var ifSyncErrorInfoTableInited bool = false

// 同步的错误信息处理对象
type syncErrorInfo struct {
	// 数据库连接对象
	db *sql.DB
}

// 初始化表信息
func (this *syncErrorInfo) init() error {
	// 初始化表结构
	if ifSyncErrorInfoTableInited == false {
		err := this.initTable(this.db)
		if err == nil {
			ifSyncErrorInfoTableInited = true
		}

		return err
	}

	return nil
}

// 把同步信息更新到数据库
// data:待更新的数据
// 返回值:
// error:错误信息
func (this *syncErrorInfo) AddErrorSql(tran *sql.Tx, data string, errMsg string) error {
	updateSql := "INSERT INTO `sync_error_info` (`SqlString`,`ExecuteTime`,`RetryCount`,`ErrMessage`) VALUES(?,?,?,?);"

	var err error
	if tran != nil {
		_, err = tran.Exec(updateSql, data, time.Now(), 0, errMsg)
	} else {
		_, err = this.db.Exec(updateSql, data, time.Now(), 0, errMsg)
	}

	if err != nil {
		logUtil.ErrorLog("logSqlSync/syncErrorInfo.AddErrorSql Error:%s", err.Error())
	}

	return err
}

// 初始化同步信息表结构
// db:数据库连接对象
func (this *syncErrorInfo) initTable(db *sql.DB) error {
	// 创建同步信息表
	createTableSql := `CREATE TABLE IF NOT EXISTS sync_error_info (
    Id bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增Id',
    SqlString varchar(1024) NOT NULL COMMENT '执行的sql',
    ExecuteTime datetime NOT NULL COMMENT '最近一次执行时间',
    RetryCount int NOT NULL COMMENT '重试次数',
    ErrMessage text NULL COMMENT '执行错误的信息',
    PRIMARY KEY (Id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='未执行成功的sql数据';`
	if _, err := db.Exec(createTableSql); err != nil {
		logUtil.ErrorLog("logSqlSync/syncErrorInfo.initTable Error:%s", err.Error())
		return err
	}

	return nil
}

// 创建同步信息对象
// _db:数据库连接对象
// 返回值:
// 同步信息对象
func newSyncErrorInfoObject(_db *sql.DB) (result *syncErrorInfo, err error) {
	result = &syncErrorInfo{
		db: _db,
	}

	err = result.init()

	return result, err
}
