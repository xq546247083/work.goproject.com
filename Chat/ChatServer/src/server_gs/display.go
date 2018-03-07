package server_gs

import (
	"bytes"
	"sync/atomic"
	"time"

	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/mathUtil"
)

var (
	// 数据量总大小
	totalSize int64 = 0
)

func init() {
	// 显示数据大小信息
	go displayDataSize()
}

// 增加接收到的数据大小
func increaseSize(size int) {
	atomic.AddInt64(&totalSize, int64(size))
}

// 显示数据大小信息
func displayDataSize() {
	// 处理goroutine数量
	goroutineName := "server_gs.displayDataSize"
	goroutineMgr.Monitor(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	for {
		// 先等待5分钟，以便服务器启动
		time.Sleep(5 * time.Minute)

		clientList := getClientList()
		var buf bytes.Buffer
		buf.WriteString("server_gs:所有客户端的地址为：")
		for _, item := range clientList {
			buf.WriteString(item.getRemoteAddr())
			buf.WriteString(",")
		}
		logUtil.DebugLog(buf.String())

		// 组装需要记录的信息
		logUtil.DebugLog("server_gs:当前客户端数量：%d, 本次服务器运行期间，共收到%s的数据", getClientCount(), mathUtil.GetSizeDesc(totalSize))
	}
}
