package goroutineMgr

/*
goroutine的管理包，提供了对goroutine的监控机制
对外提供的方法为：

// 监控指定的goroutine
Monitor(goroutineName string)

// 只添加数量，不监控
MonitorZero(goroutineName string)

// 释放监控
ReleaseMonitor(goroutineName string)

*/
