package onlineLog

// 保存在线日志
func Save(onlineLogList []*OnlineLog) {
	sid := 0
	totalCount := 0

	// 重新计算sid和totalCount
	for i := 0; i < len(onlineLogList); i++ {
		item := onlineLogList[i]
		sid += 1
		item.SetSid(sid)
		totalCount += item.ClientCount
	}

	// 保存数据
	for i := 0; i < len(onlineLogList); i++ {
		item := onlineLogList[i]
		item.SetTotalCount(totalCount)
		insert(item)
	}
}
