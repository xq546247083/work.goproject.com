package main

import (
	_ "work.goproject.com/Chat/ChatServer/src/bll"
	_ "work.goproject.com/Framework/linuxMgr"
)

import (
	"fmt"
	"sync"

	"github.com/google/gops/agent"
	"work.goproject.com/Chat/ChatServer/src/bll/dbConfig"
	"work.goproject.com/Chat/ChatServer/src/config"
	"work.goproject.com/Chat/ChatServer/src/rpcClient"
	"work.goproject.com/Chat/ChatServer/src/server_gs"
	"work.goproject.com/Chat/ChatServer/src/server_http"
	"work.goproject.com/Chat/ChatServer/src/server_tcp"
	"work.goproject.com/Chat/ChatServer/src/server_webSocket"
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
	if len(baseConfig.GetPrivateGopsAddress()) > 0 {
		if err := agent.Listen(&agent.Options{Addr: baseConfig.GetPrivateGopsAddress()}); err != nil {
			panic(err)
		}
	}

	// 注册rpcClient启动成功的通知
	rpcClientLoginSuccessName := "rpcClientLoginSuccess"
	rpcClientLoginSuccessCh := make(chan bool)
	rpcClient.LoginSuccessObj.Register(rpcClientLoginSuccessName, rpcClientLoginSuccessCh)

	// 设置rpcClient配置，并启动服务器
	go rpcClient.Start()

	// 等待rpcClient启动成功的通知通道
	<-rpcClientLoginSuccessCh
	fmt.Println("收到来自rpcClientLoginSuccessCh的消息")
	rpcClient.LoginSuccessObj.Unregister(rpcClientLoginSuccessName)

	// 设置chatServer配置，并启动服务器
	if baseConfig.Protocol == "tcp" {
		go server_tcp.Start(&wg, baseConfig.GetPrivateChatServerAddress())
	} else if baseConfig.Protocol == "websocket" {
		go server_webSocket.Start(&wg, baseConfig.GetPrivateChatServerAddress())
	}else{
		crtConfig := config.GetCrtConfig()
		go server_webSocket.Start2(&wg, baseConfig.GetPrivateChatServerAddress(),crtConfig.Crt,crtConfig.Key)
	}

	// 启动监听游戏服务器的服务器
	go server_gs.Start(&wg, baseConfig.GetPrivateGameServerAddress())

	// 启动监听游戏服务器的Web服务器
	go server_http.Start(&wg, baseConfig.GetPrivateGameServerWebAddress())

	// 阻塞等待，以免main线程退出
	wg.Wait()
}
