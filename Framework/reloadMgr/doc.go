package reloadMgr

// 重新加载包，提供重新加载的功能
// 使用方法：
// 1、先调用RegisterReloadFunc方法，将重新加载时需要调用的方法进行注册。
// 2、在需要重新加载时调用Reload()方法
