package messageLog

import (
	"time"

	"work.goproject.com/Chat/ChatServerCenter/src/bll/dbConfig"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/debugUtil"
)

func init() {
	go func() {
		// 处理goroutine数量
		goroutineName := "messageLog.clear"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		for {
			time.Sleep(time.Hour)

			// 获取数量对象
			maxIdAndMinId, err := getId()
			if err != nil {
				continue
			}

			// 删除多余的数据
			maxMessageLogCount := int64(dbConfig.GetOtherConfig().MaxMessageLogCount)
			if maxIdAndMinId.MaxId-maxIdAndMinId.MinId > maxMessageLogCount {
				debugUtil.Printf("MaxId:%d, MinId:%d, DeleteId:%d\n", maxIdAndMinId.MaxId, maxIdAndMinId.MinId, maxIdAndMinId.MaxId-maxMessageLogCount)
				delete(maxIdAndMinId.MaxId - maxMessageLogCount)
			}
		}
	}()
}
