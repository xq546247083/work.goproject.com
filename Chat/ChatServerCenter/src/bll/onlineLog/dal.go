package onlineLog

import (
	"work.goproject.com/Chat/ChatServerCenter/src/dal"
)

// 插入数据
func insert(onlineLogObj *OnlineLog) (err error) {
	result := dal.GetDB().Create(&onlineLogObj)
	if err = result.Error; err != nil {
		dal.WriteLog("onlineLog.insert", err)
		return
	}

	return nil
}
