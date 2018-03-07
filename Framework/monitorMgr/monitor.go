package monitorMgr

import (
	"fmt"
	"sync"
	"time"

	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
)

var (
	// 报告监控信息的URL
	remoteURL = "http://maintenance.work.goproject.com/Manage/Monitor.ashx"

	// 服务器IP
	serverIP string

	// 服务器名称
	serverName string

	// 监控时间间隔(单位：分钟)
	monitorInterval = 5

	// 重复消息发送的时间间隔（单位：分钟）
	duplicateInterval = 5

	// 已经发送的消息
	sentMessageMap = make(map[string]int64)

	// 已经发送消息的锁对象
	sentMessageMutex sync.Mutex

	// 监控方法列表
	monitorFuncList = make([]func() error, 0, 4)

	// 监控方法锁对象
	monitorFuncMutex sync.Mutex
)

// SetURL ...设置监控信息发送的URL
// url:监控信息发送的URL
func SetURL(url string) {
	remoteURL = url
}

// SetDuplicateInterval ...设置重复消息发送的时间间隔（单位：分钟）
// _duplicateInterval:重复消息发送的时间间隔（单位：分钟）
func SetDuplicateInterval(_duplicateInterval int) {
	duplicateInterval = _duplicateInterval
}

// SetParam ...设置参数
// _serverIP:服务器IP
// _serverName:服务器名称
// _monitorInterval:监控时间间隔(单位：分钟)
func SetParam(_serverIP, _serverName string, _monitorInterval int) {
	serverIP = _serverIP
	serverName = _serverName
	monitorInterval = _monitorInterval
}

// RegisterMonitorFunc ...注册监控方法
// f:监控方法
func RegisterMonitorFunc(f func() error) {
	monitorFuncMutex.Lock()
	defer monitorFuncMutex.Unlock()

	monitorFuncList = append(monitorFuncList, f)
}

// Start ...启动监控服务(obsolete，建议使用Start2)
// serverIp:服务器IP
// serverName:服务器名称
// monitorInterval:监控时间间隔(单位：分钟)
func Start(serverIp, serverName string, monitorInterval int) {
	monitorConfig := NewMonitorConfig(serverIp, serverName, monitorInterval)
	Start2(monitorConfig)
}

// Start ...启动监控服务2
// monitorConfig:监控配置对象
func Start2(monitorConfig *MonitorConfig) {
	// 设置参数
	SetParam(monitorConfig.ServerIp, monitorConfig.ServerName, monitorConfig.Interval)

	// 实际的监控方法调用
	monitorFunc := func() {
		monitorFuncMutex.Lock()
		defer monitorFuncMutex.Unlock()

		for _, item := range monitorFuncList {
			if err := item(); err != nil {
				Report(err.Error())
			}
		}
	}

	go func() {
		// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
		defer func() {
			if r := recover(); r != nil {
				logUtil.LogUnknownError(r)
			}
		}()

		for {
			// 先休眠，避免系统启动时就进行报警
			time.Sleep(time.Minute * time.Duration(monitorInterval))

			// 实际的监控方法调用
			monitorFunc()
		}
	}()
}

// 判断指定时间内是否已经处理过
// conent:报告内容
func isDuplicate(content string) bool {
	if ts, exists := sentMessageMap[content]; exists && time.Now().Unix()-ts < int64(60*duplicateInterval) {
		return true
	}

	return false
}

// 添加到已发送集合中
// conent:报告内容
func addToMap(content string) {
	sentMessageMap[content] = time.Now().Unix()
}

// Report ...报告异常信息
// format:报告内容格式
// args:具体参数
func Report(format string, args ...interface{}) {
	if len(args) <= 0 {
		Report2(format, "")
	} else {
		Report2(fmt.Sprintf(format, args...), "")
	}
}

// Report ...报告异常信息
// title:上报的标题
// contentFmt:报告内容
// args:参数列表
func Report2(title, contentFmt string, args ...interface{}) {
	content := contentFmt
	if len(args) > 0 {
		content = fmt.Sprintf(contentFmt, args...)
	}

	sentMessageMutex.Lock()
	defer sentMessageMutex.Unlock()

	// 判断指定时间内是否已经处理过
	if isDuplicate(title) {
		return
	}

	logUtil.NormalLog(fmt.Sprintf("MonitorReport:ServerIP:%s,ServerName:%s,Title:%s Content:%s", serverIP, serverName, title, content), logUtil.Warn)

	// 判断是否是DEBUG模式
	if debugUtil.IsDebug() {
		return
	}

	detailMsg := title
	if len(content) > 0 {
		detailMsg = fmt.Sprintf("Title:%s Content:%s", title, content)
	}

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["ServerIp"] = serverIP
	postDict["ServerName"] = serverName
	postDict["Content"] = detailMsg

	// 连接服务器，以获取数据
	returnBytes, err := webUtil.PostWebData(remoteURL, postDict, nil)
	if err != nil {
		logUtil.NormalLog(fmt.Sprintf("MonitorReport:，错误信息为：%s", err), logUtil.Error)
		return
	}

	result := string(returnBytes)
	logUtil.NormalLog(fmt.Sprintf("MonitorReport:Result:%s", result), logUtil.Warn)

	if result != "200" {
		logUtil.NormalLog(fmt.Sprintf("返回值不正确，当前返回值为：%s", result), logUtil.Error)
	}

	// 添加到已发送集合中
	addToMap(title)
}
