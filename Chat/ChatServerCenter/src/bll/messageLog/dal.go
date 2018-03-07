package messageLog

import (
	"work.goproject.com/Chat/ChatServerCenter/src/dal"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func getId() (maxIdAndMinId *MaxIdAndMinId, err error) {
	maxIdAndMinId = &MaxIdAndMinId{}

	result := dal.GetDB().Raw("SELECT MAX(Id) AS MaxId, MIN(Id) AS MinId FROM log_message").Scan(maxIdAndMinId)
	if err = result.Error; err != nil {
		dal.WriteLog("messageLog.getCount", err)
		return
	}

	return
}

func delete(id int64) (err error) {
	result := dal.GetDB().Where("Id < ?", id).Delete(MessageLog{})
	if err = result.Error; err != nil {
		dal.WriteLog("messageLog.delete", err)
		return
	}

	return
}
