package mysqlSync

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"work.goproject.com/Framework/dataSyncMgr/mysqlSync/sqlSync"
	"work.goproject.com/Framework/monitorMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

// 同步对象定义
type SyncObject struct {
	// 同步数据的存储路径
	dirPath string

	// 同步数据对象的唯一标识，用于进行重复判断
	identifier string

	// 数据库对象
	dbObj *sql.DB

	// 同步信息对象
	syncingInfoObj *syncingInfo

	// 错误处理对象
	errorHandleObj *errorFile

	// 同步对象
	syncObj *sqlSync.SyncObject
}

// 进行同步对象初始化
// maxFileSize:每个大文件的最大写入值（单位：Byte）
func (this *SyncObject) Init(baseObj *sqlSync.SyncObject) {
	this.syncObj = baseObj

	// 创建同步信息记录对象
	syncingInfoObj, err := newSyncingInfoObject(this.identifier, this.dbObj)
	if err != nil {
		panic(err)
	}
	this.syncingInfoObj = syncingInfoObj

	// 启动时同步所有数据(然后才能从数据库中查询数据，以免数据丢失)
	this.syncOldData()
}

// 同步完成之前未同步完的数据
func (this *SyncObject) syncOldData() {
	// 获取文件列表（有序的列表）
	fileList := sqlSync.GetDataFileList(this.dirPath)
	filePath, _ := this.syncingInfoObj.GetSyncingInfo()

	// 判断是否有文件
	if len(fileList) == 0 {
		return
	}

	// 判断当前文件是否为空，如果为空则将第一个文件赋给它
	if filePath == "" {
		this.syncingInfoObj.Update(fileList[0], 0, nil)
	}

	// 开始同步数据
	this.syncObj.Sync()

	return
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
func (this *SyncObject) Update(filePath string, offset int64, tran *sql.Tx) error {
	return this.syncingInfoObj.Update(filePath, offset, tran)
}

// 同步一条sql语句
// command:待执行的命令
// filePath:保存路径
// offset:文件偏移量
// 返回值:
// error:错误信息
func (this *SyncObject) SyncOneSql(command string, filePath string, offset int64) {
	err := this.syncOneSqlDetail(command, filePath, offset)
	if err == nil {
		return
	}

	//  发送监控报警
	monitorMgr.Report(fmt.Sprintf("SyncObject.SyncOneSql error Identifier:%v", this.identifier))
	this.handleError(command, filePath, offset, err)

	return
}

// 同步一条sql语句的具体逻辑
// command:待执行的命令
// filePath:保存路径
// offset:文件偏移量
// 返回值:
// error:错误信息
func (this *SyncObject) syncOneSqlDetail(command string, filePath string, offset int64) error {
	return sqlSync.ExecuteByTran(this.dbObj, func(tx *sql.Tx) (isCommit bool, err error) {
		// 保存sql到数据库
		err = this.syncToMysql(command, tx)
		if err != nil {
			return false, err
		}

		// 保存进度信息到数据库
		err = this.syncingInfoObj.Update(filePath, offset, tx)
		if err != nil {
			return false, err
		}

		return true, nil
	})
}

// 同步数据到mysql中
// command:待执行的命令
// tx:事务对象
// 返回值:
// error:错误信息
func (this *SyncObject) syncToMysql(command string, tx *sql.Tx) error {
	_, err := tx.Exec(command)
	if err != nil {
		prefix := fmt.Sprintf("%s-%s", this.identifier, "SyncObject.syncToMysql")
		err = fmt.Errorf("%s-%s Update to mysql failed:%s", prefix, command, err)
		logUtil.ErrorLog(err.Error())
		debugUtil.Printf("fatal Error:%v", err.Error())
		return err
	}

	return nil
}

// 进行错误处理
// command:存在异常的数据
// filePath:文件路径
// offset:文件偏移量
// err:错误信息
func (this *SyncObject) handleError(command string, filePath string, offset int64, err error) {
	defer this.errorHandleObj.Delete()

	// 保存当前sql命令
	this.errorHandleObj.SaveCommand(command)

	// 循环处理当前命令，直到没有错误
	beginTime := time.Now().Unix()
	for {
		// 每隔5分钟，发送警报
		if time.Now().Unix()-beginTime > 5*60 {
			monitorMgr.Report2(fmt.Sprintf("SyncObject.handleError error Identifier:%v", this.identifier), err.Error())
			beginTime = time.Now().Unix()
		}

		// 每次循环休眠20秒
		time.Sleep(5 * time.Second)
		command = this.errorHandleObj.ReadCommand()
		err = this.syncOneSqlDetail(command, filePath, offset)
		if err != nil {
			continue
		}

		break
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
func newSyncObject(dirPath, identifier string, dbObj *sql.DB) *SyncObject {
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
		dirPath:        dirPath,
		identifier:     identifier,
		dbObj:          dbObj,
		errorHandleObj: newErrorFile(dirPath, identifier),
	}

	return result
}
