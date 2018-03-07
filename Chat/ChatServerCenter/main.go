package main

import (
	_ "work.goproject.com/Chat/ChatServerCenter/src/bll"
	_ "work.goproject.com/Framework/linuxMgr"
)

import (
	"sync"

	"github.com/google/gops/agent"
	"work.goproject.com/Chat/ChatServerCenter/src/bll/dbConfig"
	"work.goproject.com/Chat/ChatServerCenter/src/config"
	"work.goproject.com/Chat/ChatServerCenter/src/server_http"
	"work.goproject.com/Chat/ChatServerCenter/src/server_tcp"
	"work.goproject.com/Framework/managecenterMgr"
	"work.goproject.com/Framework/monitorMgr"
	"work.goproject.com/Framework/signalMgr"
)

var (
	wg sync.WaitGroup
)

func init() {
	// 设置WaitGroup需要等待的数量，只要有一个服务器出现错误都停止服务器
	wg.Add(1)
}

func main() {
	// 启动信号处理程序
	signalMgr.Start()

	// 启动managecentert管理器
	managecenterMgr.Start2(dbConfig.ManageCenterConfig)

	// 启动监控处理程序
	monitorMgr.Start2(config.GetMonitorConfig())

	// 获取配置信息
	baseConfig := config.GetBaseConfig()

	// 启动监控代理程序
	if len(baseConfig.GopsAddr) > 0 {
		if err := agent.Listen(&agent.Options{Addr: baseConfig.GopsAddr}); err != nil {
			panic(err)
		}
	}

	// 启动监听ChatServer的服务器
	go server_tcp.Start(&wg, baseConfig.ChatServerAddress)

	// 设置Web服务器配置，并启动服务器
	go server_http.Start(&wg, baseConfig.WebServerAddress)

	// 阻塞等待，以免main线程退出
	wg.Wait()
}
