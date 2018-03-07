package sqlSync

import (
	"database/sql"
	"strings"
	"time"

	"work.goproject.com/goutil/logUtil"
)

// 以事务的方式执行
// db:数据库对象
// funcObj:对应的具体处理函数
// 返回值:
// error:处理是否存在错误
func ExecuteByTran(db *sql.DB, funcObj func(tran *sql.Tx) (isCommit bool, err error)) error {
	tran, err := db.Begin()
	if err != nil {
		logUtil.ErrorLog("start transaction error:%v", err.Error())
		return err
	}

	// 事务处理
	isCommit := false
	defer func() {
		if isCommit {
			err = tran.Commit()
		} else {
			err = tran.Rollback()
		}

		if err != nil {
			logUtil.ErrorLog("transaction end error:%v", err.Error())
		}
	}()

	isCommit, err = funcObj(tran)

	return err
}

// 循环执行知道返回成功为止
// funcObj:待执行的函数
// interval:执行间隔时间
func WaitForOk(funcObj func() bool, interval time.Duration) {
	for {
		if funcObj() == false {
			time.Sleep(interval)
		}

		break
	}
}

// 检查是否是连接错误
// errMsg:错误信息
// 返回值:
// bool:true：连接错误 false:其他异常
func CheckIfConnectionError(errMsg string) bool {
	//// 连接被关闭
	ifConnectionClose := strings.Contains(errMsg, "A connection attempt failed because the connected party did not properly respond")
	if ifConnectionClose {
		return true
	}

	// 使用过程中连接断开
	ifConnectionClose = strings.Contains(errMsg, "No connection could be made")
	if ifConnectionClose {
		return true
	}

	// 事务处理过程中连接断开的提示
	ifConnectionClose = strings.Contains(errMsg, "bad connection")
	if ifConnectionClose {
		return true
	}

	// socket压根儿连不上的处理
	ifConnectionClose = strings.Contains(errMsg, "A socket operation was attempted to an unreachable network")
	if ifConnectionClose {
		return true
	}

	// 用户无法访问
	return strings.Contains(errMsg, "Access denied for user")
}

// 获取比较简洁的错误信息
// errMsg:错误信息
// 返回值:
// string:比较简洁的错误信息
func GetSimpleErrorMessage(errMsg string) string {
	if strings.Contains(errMsg, "Error 1064: You have an error in your SQL syntax") {
		return "SqlError"
	}

	return errMsg
}
