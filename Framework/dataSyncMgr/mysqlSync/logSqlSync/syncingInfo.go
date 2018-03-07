package logSqlSync

import (
	"database/sql"
	"fmt"
	"time"

	"work.goproject.com/goutil/logUtil"
)

// 同步信息表是否已经被初始化
var ifSyncingTableInited bool = false

// 同步信息项，保存已经处理过的文件的信息
type syncingModel struct {
	// 服务器组Id
	ServerGroupId int32

	// 待处理文件的绝对路径
	FilePath string

	// 待处理文件的偏移量
	FileOffset int64

	// 更新时间
	UpdateTime time.Time
}

// 同步信息对象
type syncingInfo struct {
	// 服务器组Id
	ServerGroupId int32

	// 同步信息项
	item *syncingModel

	// 数据库连接对象
	db *sql.DB
}

// 获取同步信息
// filePath:正在同步的文件
// fileOffset:同步到的位置
func (this *syncingInfo) GetSyncingInfo() (filePath string, fileOffset int64) {
	return this.item.FilePath, this.item.FileOffset
}

// 更新正在同步的位置和文件信息
// filePath:文件路径
// offset:当前同步到的位置
// tran:事务对象，可以为nil
// 返回值:
// error:处理的错误信息
func (this *syncingInfo) Update(filePath string, offset int64, tran *sql.Tx) error {
	this.item.FilePath = filePath
	this.item.FileOffset = offset
	this.item.UpdateTime = time.Now()

	// 更新到数据库
	return this.update(this.item, tran)
}

// 初始化同步信息
// 返回值:
// error:错误信息
func (this *syncingInfo) init() error {
	// 数据表初始化
	if ifSyncingTableInited == false {
		if err := this.initSyncingInfoTable(this.db); err == nil {
			ifSyncingTableInited = true
		} else {
			return err
		}
	}

	// 获取此表的同步信息
	data, exist, err := this.get()
	if err != nil {
		return err
	}

	// 2. 如果同步信息不存在，则初始化一条到此表
	if exist == false {
		data = &syncingModel{
			ServerGroupId: this.ServerGroupId,
			FilePath:      "",
			FileOffset:    0,
			UpdateTime:    time.Now(),
		}
	}

	this.item = data
	return nil
}

// 初始化同步信息表结构
// db:数据库连接对象
func (this *syncingInfo) initSyncingInfoTable(db *sql.DB) error {
	// 创建同步信息表
	createTableSql := `CREATE TABLE IF NOT EXISTS syncing_info (
    ServerGroupId int NOT NULL COMMENT '服务器组Id',
    FilePath varchar(500) NOT NULL COMMENT '正在同步的文件路径',
    FileOffset bigint(20) NOT NULL COMMENT '偏移量',
    UpdateTime datetime NOT NULL COMMENT '最后一次更新时间',
    PRIMARY KEY (ServerGroupId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='正在同步的文件信息';`
	if _, err := db.Exec(createTableSql); err != nil {
		logUtil.ErrorLog("logSqlSync/syncingInfo.initSyncingInfoTable Error:%s", err.Error())
		return err
	}

	return nil
}

// 从数据库获取数据
// 返回值:
// data:获取到的数据
// exist:是否存在此数据
// err:错误信息
func (this *syncingInfo) get() (data *syncingModel, exist bool, err error) {
	//// 从数据库查询
	querySql := fmt.Sprintf("SELECT FilePath,FileOffset,UpdateTime FROM syncing_info WHERE ServerGroupId ='%v'", this.ServerGroupId)
	var rows *sql.Rows
	rows, err = this.db.Query(querySql)
	if err != nil {
		logUtil.ErrorLog("logSqlSync/syncingInfo.get.Query ServerGroupId:%v error:%s", this.ServerGroupId, err.Error())
		return
	}
	defer rows.Close()

	if rows.Next() == false {
		exist = false
		return
	}
	exist = true

	// 读取数据
	data = &syncingModel{
		ServerGroupId: this.ServerGroupId,
	}
	err = rows.Scan(&data.FilePath, &data.FileOffset, &data.UpdateTime)
	if err != nil {
		logUtil.ErrorLog("logSqlSync/syncingInfo.get.Query ServerGroupId:%v error:%s", this.ServerGroupId, err.Error())
		return
	}

	return
}

// 把同步信息更新到数据库
// data:待更新的数据
// tran:事务处理对象
// 返回值:
// error:错误信息
func (this *syncingInfo) update(data *syncingModel, tran *sql.Tx) error {
	updateSql := "REPLACE INTO `syncing_info` SET `ServerGroupId` = ?, `FilePath` = ?,`FileOffset` = ?, `UpdateTime` = ?;"
	var err error
	if tran != nil {
		_, err = tran.Exec(updateSql, data.ServerGroupId, data.FilePath, data.FileOffset, data.UpdateTime)
	} else {
		_, err = this.db.Exec(updateSql, data.ServerGroupId, data.FilePath, data.FileOffset, data.UpdateTime)
	}

	if err != nil {
		logUtil.ErrorLog("logSqlSync/syncingInfo.update ServerGroupId:%v error:%s", this.ServerGroupId, err.Error())
	}

	return err
}

// 创建同步信息对象
// _dirPath:目录的路径
// _identifier:当前数据的唯一标识（可以使用数据库表名）
// _db:数据库连接对象
// 返回值:
// 同步信息对象
func newSyncingInfoObject(serverGroupId int32, _db *sql.DB) (result *syncingInfo, err error) {
	result = &syncingInfo{
		ServerGroupId: serverGroupId,
		db:            _db,
	}

	err = result.init()

	return result, err
}
