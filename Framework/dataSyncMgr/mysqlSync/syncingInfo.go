package mysqlSync

import (
	"database/sql"
	"fmt"
	"time"

	"work.goproject.com/goutil/logUtil"
)

var (
	// 是否已经初始化了正在同步的表信息
	ifSyncingInfoTableInited = false

	// 表初始化错误信息
	initTableError error = nil

	// 表是否已经初始化
	isTableInited bool = false
)

// 同步信息项，保存已经处理过的文件的信息
type syncingModel struct {
	// 唯一标识
	Identifier string

	// 待处理文件的绝对路径
	FilePath string

	// 待处理文件的偏移量
	FileOffset int64

	// 更新时间
	UpdateTime time.Time
}

// 同步信息对象
type syncingInfo struct {
	// 同步数据对象的唯一标识，用于进行重复判断
	identifier string

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
	if ifSyncingInfoTableInited == false {
		err := initSyncingInfoTable(this.db)
		if err != nil {
			return err
		}

		ifSyncingInfoTableInited = true
	}

	// 获取此表的同步信息
	data, exist, err := this.get()
	if err != nil {
		return err
	}

	// 2. 如果同步信息不存在，则初始化一条到此表
	if exist == false {
		data = &syncingModel{
			Identifier: this.identifier,
			FilePath:   "",
			FileOffset: 0,
			UpdateTime: time.Now(),
		}

		err = this.insert(data)
		if err != nil {
			return err
		}
	}

	this.item = data
	return nil
}

// 从数据库获取数据
// 返回值:
// data:获取到的数据
// exist:是否存在此数据
// err:错误信息
func (this *syncingInfo) get() (data *syncingModel, exist bool, err error) {
	//// 从数据库查询
	querySql := fmt.Sprintf("SELECT FilePath,FileOffset,UpdateTime FROM syncing_info WHERE Identifier ='%v';", this.identifier)
	var rows *sql.Rows
	rows, err = this.db.Query(querySql)
	if err != nil {
		logUtil.ErrorLog("mysqlSync/syncingInfo get.query error:%v", err.Error())
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
		Identifier: this.identifier,
	}
	err = rows.Scan(&data.FilePath, &data.FileOffset, &data.UpdateTime)
	if err != nil {
		logUtil.ErrorLog("mysqlSync/syncingInfo get.scan error:%v", err.Error())
		return
	}

	return
}

// 把同步信息写入到数据库
// data:待插入的数据
// 返回值:
// error:错误信息
func (this *syncingInfo) insert(data *syncingModel) error {
	insertSql := "INSERT INTO syncing_info(Identifier,FilePath,FileOffset,UpdateTime) VALUES(?,?,?,?);"
	_, err := this.db.Exec(insertSql, data.Identifier, data.FilePath, data.FileOffset, data.UpdateTime)

	if err != nil {
		logUtil.ErrorLog("mysqlSync/syncingInfo insert error:%v", err.Error())
	}

	return err
}

// 把同步信息更新到数据库
// data:待更新的数据
// tran:事务对象
// 返回值:
// error:错误信息
func (this *syncingInfo) update(data *syncingModel, tran *sql.Tx) error {
	updateSql := "UPDATE syncing_info SET FilePath=?, FileOffset=?, UpdateTime=? WHERE Identifier=?;"
	var err error
	if tran != nil {
		_, err = tran.Exec(updateSql, data.FilePath, data.FileOffset, data.UpdateTime, data.Identifier)
	} else {
		_, err = this.db.Exec(updateSql, data.FilePath, data.FileOffset, data.UpdateTime, data.Identifier)
	}

	if err != nil {
		logUtil.ErrorLog("mysqlSync/syncingInfo update error:%v", err.Error())
	}

	return err
}

// 创建同步信息对象
// _dirPath:目录的路径
// _identifier:当前数据的唯一标识（可以使用数据库表名）
// _db:数据库连接对象
// 返回值:
// 同步信息对象
func newSyncingInfoObject(identifier string, _db *sql.DB) (result *syncingInfo, err error) {
	result = &syncingInfo{
		identifier: identifier,
		db:         _db,
	}

	err = result.init()

	return result, err
}

// 初始化同步信息表结构
// db:数据库连接对象
func initSyncingInfoTable(db *sql.DB) error {
	if isTableInited {
		return initTableError
	}

	defer func() {
		isTableInited = true
	}()

	// 创建同步信息表
	createTableSql := `CREATE TABLE IF NOT EXISTS syncing_info (
	Identifier  varchar(30) NOT NULL COMMENT '同步唯一标识(数据库表名)',
		FilePath  varchar(500) NOT NULL COMMENT '正在同步的文件路径',
		FileOffset  bigint NOT NULL COMMENT '文件偏移量',
		UpdateTime  datetime NOT NULL COMMENT '最后一次更新时间',
		PRIMARY KEY (Identifier)
	) COMMENT 'P表同步信息';`
	if _, initTableError = db.Exec(createTableSql); initTableError != nil {
		logUtil.ErrorLog("mysqlSync/syncingInfo initSyncingInfoTable error:%v", initTableError.Error())
		return initTableError
	}

	return nil
}
