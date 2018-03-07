package exitMgr

// 程序退出包，提供程序退出时的功能
// 使用方法：
// 1、先调用RegisterExitFunc方法，将系统退出时需要调用的方法进行注册。
// 2、在程序退出时调用Exit()方法
