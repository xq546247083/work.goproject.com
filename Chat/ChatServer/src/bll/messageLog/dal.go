package messageLog

import (
	"work.goproject.com/Chat/ChatServer/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

// 插入数据
func insert(messageLogObj *MessageLog) (id int, err error) {
	result := dal.GetDB().Create(messageLogObj)
	if err = result.Error; err != nil {
		dal.WriteLog("messageLog.insert", err)
		return
	}

	id = messageLogObj.Id

	return
}
