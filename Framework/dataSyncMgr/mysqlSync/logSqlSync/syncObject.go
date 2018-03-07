package logSqlSync

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"time"

	"work.goproject.com/Framework/dataSyncMgr/mysqlSync/sqlSync"
	"work.goproject.com/Framework/monitorMgr"
	"work.goproject.com/goutil/logUtil"
)

// 同步对象定义
type SyncObject struct {
	// 服务器组Id
	serverGroupId int32

	// 同步数据的存储路径
	dirPath string

	// 同步数据对象的唯一标识，用于进行重复判断
	identifier string

	// 数据库对象
	dbObj *sql.DB

	// 同步信息对象
	syncingInfoObj *syncingInfo

	// 错误处理对象
	errorHandleObj *syncErrorInfo

	// 同步对象
	syncObj *sqlSync.SyncObject
}

// 初始化
// baseObj:基础同步对象
func (this *SyncObject) Init(baseObj *sqlSync.SyncObject) {
	this.syncObj = baseObj

	// 初始化同步信息对象
	syncingInfoObj, err := newSyncingInfoObject(this.serverGroupId, this.dbObj)
	if err != nil {
		panic(err)
	}
	//// 初始化错误处理对象
	errorHandleObj, err := newSyncErrorInfoObject(this.dbObj)
	if err != nil {
		panic(err)
	}
	this.syncingInfoObj = syncingInfoObj
	this.errorHandleObj = errorHandleObj

	// 初始化当前处理的文件
	fileList := sqlSync.GetDataFileList(this.dirPath)
	filePath, _ := this.syncingInfoObj.GetSyncingInfo()
	if len(filePath) < 0 && len(fileList) > 0 {
		this.syncingInfoObj.Update(fileList[0], 0, nil)
	}
}

// 获取正在同步的信息
// filePath:文件路径
// offset:文件偏移量
func (this *SyncObject) GetSyncingInfo() (filePath string, offset int64) {
	return this.syncingInfoObj.GetSyncingInfo()
}

// 更新
// filePath:文件路径
// offset:文件偏移量
// tran:事务对象
// 返回值:
// error:错误对象
func (this *SyncObject) Update(filePath string, offset int64, tx *sql.Tx) error {
	return this.syncingInfoObj.Update(filePath, offset, tx)
}

// 同步一条sql语句
// command:待执行的命令
// filePath:保存路径
// offset:文件偏移量
// 返回值:
// error:错误信息
func (this *SyncObject) SyncOneSql(command string, filePath string, offset int64) {
	var err error
	for {
		err = sqlSync.ExecuteByTran(this.dbObj, func(tran *sql.Tx) (isCommit bool, err error) {
			// 保存sql到数据库
			err = this.syncToMysql(command, tran)
			if err != nil {
				return
			}

			// 保存进度信息到数据库
			err = this.syncingInfoObj.Update(filePath, offset, tran)
			if err != nil {
				return
			}

			isCommit = true
			return
		})

		// 如果是连接出错，则仍然循环执行
		if err != nil {
			monitorMgr.Report2("logSqlSync/syncObject.SyncOneSql Error", err.Error())

			if sqlSync.CheckIfConnectionError(err.Error()) {
				time.Sleep(5 * time.Second)
				continue
			}
		}

		// 如果不是数据库连接出错，则算是执行完成
		break
	}

	// 如果存在错误，则循环尝试执行
	if err != nil {
		this.recordSqlError(command, filePath, offset, err.Error())
	}

	return
}

// 同步数据到mysql中
// command:sql语句
// tx:事务处理对象
// 返回值:
// error:错误信息
func (this *SyncObject) syncToMysql(command string, tx *sql.Tx) error {
	_, err := tx.Exec(command)
	if err != nil {
		logUtil.ErrorLog("mysqlSync/logSqlSync/syncObject.syncToMysql error:%s", err.Error())
		return err
	}

	return nil
}

// 错误处理
// cmd:待执行的命令
// filePath:保存路径
// offset:文件偏移量
// errMsg:错误信息
func (this *SyncObject) recordSqlError(command string, filePath string, offset int64, errMsg string) {
	errMsg = sqlSync.GetSimpleErrorMessage(errMsg)

	for {
		err := sqlSync.ExecuteByTran(this.dbObj, func(tran *sql.Tx) (isCommit bool, err error) {
			// 保存sql到数据库
			err = this.errorHandleObj.AddErrorSql(tran, command, errMsg)
			if err != nil {
				return
			}

			// 保存进度信息到数据库
			err = this.syncingInfoObj.Update(filePath, offset, tran)
			if err != nil {
				return
			}

			isCommit = true
			return
		})

		if err == nil {
			return
		}

		monitorMgr.Report2("logSqlSync/syncObject.recordSqlError Error", err.Error())
		time.Sleep(5 * time.Second)
	}
}

// 创新新的mysql同步对象
// dirPath:存放数据的目录
// identifier:当前数据的唯一标识（可以使用数据库表名）
// dbObj:数据库对象
// syncingInfoObj:同步信息记录对象
// errorHandleObj:错误处理对象
// 返回值:
// mysql同步对象
func NewSyncObject(serverGroupId int32, dirPath, identifier string, dbObj *sql.DB) *SyncObject {
	dirPath = filepath.Join(dirPath, identifier)

	// 创建更新目录
	err := os.MkdirAll(dirPath, os.ModePerm|os.ModeTemporary)
	if err != nil {
		err = fmt.Errorf("%s-%s-make dir failed:%s", identifier, "SyncObject.newSyncObject.os.MkdirAll", err)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}

	// 构造同步信息对象
	result := &SyncObject{
		serverGroupId: serverGroupId,
		dirPath:       dirPath,
		identifier:    identifier,
		dbObj:         dbObj,
	}

	return result
}
