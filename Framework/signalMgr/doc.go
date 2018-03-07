package signalMgr

// 系统信号管理包
// 提供对操作系统信号的管理，支持三种信号：syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP
// 其中syscall.SIGTERM, syscall.SIGINT表示程序终止信号，可以绑定一个方法在系统终止时进行调用
// syscall.SIGHUP表示程序重启信号，可以绑定一个方法在系统重启时进行调用

// 使用方法：
// 调用func Start(reloadFunc func() []error, exitFunc func() error)
// 传入在重启和终止时需要调用的方法，如果不需要则传入nil
