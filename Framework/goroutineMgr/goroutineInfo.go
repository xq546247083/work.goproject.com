package goroutineMgr

import (
	"fmt"
	"sync"

	"work.goproject.com/Framework/monitorMgr"
	"work.goproject.com/goutil/logUtil"
)

var (
	goroutineInfoMap   map[string]int = make(map[string]int)
	goroutineInfoMutex sync.RWMutex
)

func init() {
	monitorMgr.RegisterMonitorFunc(monitor)
}

func registerGoroutineInfo(goroutineName string, count int) {
	goroutineInfoMutex.Lock()
	defer goroutineInfoMutex.Unlock()

	goroutineInfoMap[goroutineName] = count
}

// 监控指定的goroutine
func Monitor(goroutineName string) {
	increaseCount(goroutineName)
	registerGoroutineInfo(goroutineName, 1)
}

// 只添加数量，不监控
func MonitorZero(goroutineName string) {
	increaseCount(goroutineName)
}

// 释放监控
func ReleaseMonitor(goroutineName string) {
	if r := recover(); r != nil {
		logUtil.LogUnknownError(r)
	}

	decreaseCount(goroutineName)
}

func monitor() error {
	/*
		先记录活跃的goroutine的数量信息
		然后再判断数量是否匹配
	*/
	logGoroutineCountInfo()

	goroutineInfoMutex.RLock()
	defer goroutineInfoMutex.RUnlock()

	for goroutineName, count := range goroutineInfoMap {
		if currCount := getGoroutineCount(goroutineName); currCount != count {
			return fmt.Errorf("%s需要%d个goroutine，现在有%d个", goroutineName, count, currCount)
		}
	}

	return nil
}
