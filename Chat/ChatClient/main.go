package main

import (
	"sync"

	"work.goproject.com/Chat/ChatClient/src/config"
	"work.goproject.com/Chat/ChatClient/src/rpc"
	"work.goproject.com/Framework/managecenterMgr"
)

var (
	wg sync.WaitGroup
)

func init() {
	wg.Add(1)
}

func main() {
	// 启动managecentert管理器
	managecenterMgr.Start2(config.GetManageCenterConfig())

	// 启动客户端
	startCh := make(chan int)
	go rpc.StartClient(startCh)
	<-startCh

	// 与用户交互
	ch := make(chan int)
	go rpc.Interaction(ch)

	<-ch
}
