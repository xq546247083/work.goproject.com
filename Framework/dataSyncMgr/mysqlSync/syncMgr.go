package mysqlSync

import (
	"database/sql"
	"fmt"
	"sync"

	"work.goproject.com/Framework/dataSyncMgr/mysqlSync/logSqlSync"
	"work.goproject.com/Framework/dataSyncMgr/mysqlSync/sqlSync"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

// 数据同步管理
type SyncMgr struct {
	// 服务器组Id
	serverGroupId int32

	// 同步数据的存储路径
	dirPath string

	// 大文件对象size
	maxFileSize int

	// 数据库对象
	dbObj *sql.DB

	// 同步对象集合
	syncObjMap map[string]*sqlSync.SyncObject

	// 同步对象锁
	mutex sync.RWMutex

	// 新建实例对象的函数
	newInstanceFunc func(mgr *SyncMgr, identifier string) *sqlSync.SyncObject
}

// 注册同步对象
// identifier:当前数据的唯一标识（可以使用数据库表名）
func (this *SyncMgr) RegisterSyncObj(identifier string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 判断是否设置了相同的唯一标识，以免弄混淆
	if _, exists := this.syncObjMap[identifier]; exists {
		prefix := fmt.Sprintf("%s-%s", identifier, "SyncMgr.RegisterSyncObj")
		err := fmt.Errorf("%s has already existed, please change another identifier", prefix)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}

	syncObj := this.newInstanceFunc(this, identifier)
	syncObj.Init(this.maxFileSize)
	this.syncObjMap[identifier] = syncObj

	if debugUtil.IsDebug() {
		fmt.Printf("%s同步对象成功注册进SyncMgr, 当前有%d个同步对象\n", identifier, len(this.syncObjMap))
	}
}

// 保存数据
// identifier:当前数据的唯一标识（可以使用数据库表名）
// command:sql命令
func (this *SyncMgr) Save(identifier string, command string) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	syncObj, exists := this.syncObjMap[identifier]
	if !exists {
		err := fmt.Errorf("syncObj:%s does not existed, please register first", identifier)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}

	syncObj.Save(command)
}

// 构造同步管理对象
// serverGroupId:服务器组Id
// dirPath: 文件目录
// maxFileSize: 大文件对象大小
// survivalTime: 同步数据存活时间 (单位：hour)
// dbObj: 数据库对象
func NewSyncMgr(serverGroupId int32, dirPath string, maxFileSize int, survivalTime int, dbObj *sql.DB) *SyncMgr {
	result := &SyncMgr{
		serverGroupId: serverGroupId,
		dirPath:       dirPath,
		maxFileSize:   maxFileSize,
		dbObj:         dbObj,
		syncObjMap:    make(map[string]*sqlSync.SyncObject),
		newInstanceFunc: func(mgr *SyncMgr, identifier string) *sqlSync.SyncObject {
			handler := newSyncObject(mgr.dirPath, identifier, mgr.dbObj)
			return sqlSync.NewSyncObject(mgr.dirPath, identifier, mgr.dbObj, handler)
		},
	}

	return result
}

// 新建日志同步管理对象
// serverGroupId:服务器组Id
// dirPath: 文件目录
// maxFileSize: 大文件对象大小
// dbObj: 数据库对象
func NewLogSyncMgr(serverGroupId int32, dirPath string, maxFileSize int, dbObj *sql.DB) *SyncMgr {
	result := &SyncMgr{
		serverGroupId: serverGroupId,
		dirPath:       dirPath,
		maxFileSize:   maxFileSize,
		dbObj:         dbObj,
		syncObjMap:    make(map[string]*sqlSync.SyncObject),
		newInstanceFunc: func(mgr *SyncMgr, identifier string) *sqlSync.SyncObject {
			handler := logSqlSync.NewSyncObject(mgr.serverGroupId, mgr.dirPath, identifier, mgr.dbObj)
			return sqlSync.NewSyncObject(mgr.dirPath, identifier, mgr.dbObj, handler)
		},
	}

	return result
}
