package sqlSync

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"work.goproject.com/Framework/monitorMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/fileUtil"
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

	// 处理数据写入的文件
	sqlFileObj *SqlFile

	// 同步处理对象
	syncHandleObj SyncHandler
}

// 进行同步对象初始化
// maxFileSize:每个大文件的最大写入值（单位：Byte）
func (this *SyncObject) Init(maxFileSize int) {
	// 启动时同步所有数据(然后才能从数据库中查询数据，以免数据丢失)
	this.syncHandleObj.Init(this)

	// 构造同步数据对象
	fileName, _ := this.syncHandleObj.GetSyncingInfo()
	this.sqlFileObj = NewSqlFile(this.dirPath, this.identifier, fileName, maxFileSize)

	// 当前没有正在同步的文件，则指向当前正在写的文件
	if len(fileName) <= 0 {
		this.syncHandleObj.Update(this.sqlFileObj.FileFullName(), 0, nil)
	}

	// 启动一个新goroutine来负责同步数据
	go func() {
		/* 此处不使用goroutineMgr.Monitor/ReleaseMonitor，因为此处不能捕获panic，需要让外部进程终止执行，
		因为此模块的文件读写为核心逻辑，一旦出现问题必须停止进程，否则会造成脏数据
		*/
		this.Sync()
	}()
}

// 保存数据到本地文件
// command:待保存的指令
func (this *SyncObject) Save(command string) {
	this.sqlFileObj.Write(command)
}

// 循环同步多个文件
func (this *SyncObject) Sync() {
	// 开始循环同步
	for {
		// 同步当前文件
		this.syncOneFile()

		// 当前文件同步完成，记录同步日志
		nowFilePath, _ := this.syncHandleObj.GetSyncingInfo()

		// 删除已同步完成的文件
		WaitForOk(func() bool {
			fileExist, err := fileUtil.IsFileExists(nowFilePath)
			if err != nil {
				logUtil.ErrorLog("mysqlSync/syncObject IsFileExists error:%s", err.Error())
				monitorMgr.Report2("mysqlSync/syncObject IsFileExists error", err.Error())
				return false
			}
			if fileExist == false {
				return true
			}

			err = fileUtil.DeleteFile(nowFilePath)
			if err != nil {
				logUtil.ErrorLog("mysqlSync/syncObject delete file error:%s", err.Error())
				monitorMgr.Report2("mysqlSync/syncObject delete file error", err.Error())

				return false
			}

			return true
		}, 10*time.Second)

		// 当前文件同步完成，获取下个文件
		nextFileName := NewFileName("", nowFilePath)
		filePath := filepath.Join(this.dirPath, nextFileName)
		exist, err := fileUtil.IsFileExists(filePath)
		if err != nil {
			logUtil.ErrorLog("mysqlSync/syncObject IsFileExists error:%s", err.Error())
			monitorMgr.Report2("mysqlSync/syncObject IsFileExists error", err.Error())
			panic(err)
		}

		// 如果文件不存在，退出
		if !exist {
			// fmt.Println("协程退出了")
			return
		}

		// 更新同步的位置信息 此处忽略错误是因为，哪怕是出错了，也不会影响整体逻辑
		this.syncHandleObj.Update(filePath, 0, nil)
	}
}

// 同步单个文件
func (this *SyncObject) syncOneFile() {
	// 获取信息同步项对象
	filePath, offset := this.syncHandleObj.GetSyncingInfo()

	// 打开待读取的文件
	f, exist := this.openFile(filePath)
	if exist == false {
		// logUtil.WarnLog("待同步的文件不存在，跳过此文件:%s", filePath)
		return
	}
	defer f.Close()

	for {
		// 移动到需要读取的位置
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			prefix := fmt.Sprintf("%s-%s", this.identifier, "SyncObject.Seek")
			err = fmt.Errorf("%s-Seek offset for header failed:%s", prefix, err)
			logUtil.ErrorLog(err.Error())
			monitorMgr.Report(fmt.Sprintf("SyncObject.Seek error， Identifier:%v", this.identifier))
			panic(err)
		}

		command, readLen, err := Read(f)
		if err != nil {
			// 如果读取到文件末尾，判断是否等待
			if err == io.EOF {
				if this.sqlFileObj != nil && strings.Contains(filePath, this.sqlFileObj.FileName()) {
					time.Sleep(20 * time.Millisecond)
					continue
				}

				// 如果该文件是空文件,同步更新信息
				return
			}

			prefix := fmt.Sprintf("%s-%s", this.identifier, "SyncObject.syncOneFile.f.Read")
			err = fmt.Errorf("%s-Read header failed:%s", prefix, err)
			logUtil.ErrorLog(err.Error())
			monitorMgr.Report(fmt.Sprintf("SyncObject.syncOneFile.f.Read error Identifier:%v", this.identifier))
			panic(err)
		}

		// 3. 同步到mysql中,并更新同步位置
		this.syncHandleObj.SyncOneSql(command, filePath, offset+readLen)

		// 4. 更新内存中的同步位置
		offset += readLen
	}
}

// 打开待读取的文件
// filePath:待打开的文件
// 返回值:
// *os.File:文件句柄，
func (this *SyncObject) openFile(filePath string) (f *os.File, exist bool) {
	var err error
	for {
		exist, err = fileUtil.IsFileExists(filePath)
		if err != nil {
			err = fmt.Errorf("check file error,filePath:%v  error:%v", filePath, err.Error())
			logUtil.ErrorLog(err.Error())
			monitorMgr.Report("SyncObject.OpenFile check file error")
			time.Sleep(time.Second * 5)
			continue
		}
		if exist == false {
			// 如果文件不存在，则跳过此文件
			logUtil.WarnLog("file no exist, skip file:%v", filePath)
			exist = false
			return
		}
		exist = true

		// 打开当前处理文件
		f, err = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm|os.ModeTemporary)
		if err != nil {
			prefix := fmt.Sprintf("%s-%s", this.identifier, "SyncObject.syncOneFile.os.OpenFile")
			err = fmt.Errorf("%s-Open file:%s failed:%s", prefix, filePath, err)
			logUtil.ErrorLog(err.Error())
			monitorMgr.Report("SyncObject.OpenFile open file error")

			time.Sleep(time.Second * 5)
			continue
		}

		return
	}
}

// 同步数据到mysql中
// command:sql语句
// tx:事务处理对象
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

// 创新新的mysql同步对象
// dirPath:存放数据的目录
// identifier:当前数据的唯一标识（可以使用数据库表名）
// dbObj:数据库对象
// _syncHandleObj:同步处理对象
// 返回值:
// mysql同步对象
func NewSyncObject(dirPath, identifier string, dbObj *sql.DB, _syncHandleObj SyncHandler) *SyncObject {
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
		dirPath:       dirPath,
		identifier:    identifier,
		dbObj:         dbObj,
		syncHandleObj: _syncHandleObj,
	}

	return result
}

// 同步处理接口
type SyncHandler interface {
	// 初始化
	Init(baseObj *SyncObject)

	// 获取正在同步的信息
	// filePath:文件路径
	// offset:文件偏移量
	GetSyncingInfo() (filePath string, offset int64)

	// 更新
	// filePath:文件路径
	// offset:文件偏移量
	// tran:事务对象
	// 返回值:
	// error:错误对象
	Update(filePath string, offset int64, tran *sql.Tx) error

	// 同步一条sql
	// command:指令数据
	// filePath:文件路径
	// offset:文件偏移量
	SyncOneSql(command string, filePath string, offset int64)
}
