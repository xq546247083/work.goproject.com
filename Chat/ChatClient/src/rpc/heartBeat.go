package rpc

import (
	"work.goproject.com/goutil/logUtil"
	"time"
)

func heartBeat() {
	// 处理内部未处理的异常，避免导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
		}
	}()

	for {
		// 由于连接刚刚建立，所以无需发心跳包；等待一段时间之后再发
		time.Sleep(30 * time.Second)

		clientObj.sendHeartBeatMessage()
	}
}
