package rpcClient

import (
	"fmt"
	"net"
	"time"

	"work.goproject.com/Chat/ChatServer/src/config"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/Framework/initMgr"
	"work.goproject.com/goutil/logUtil"
)

var (
	// 客户端对象
	clientObj = newClient()

	// 启动成功对象
	startSuccessObj = initMgr.NewInitSuccess("rpcClientStart")
)

func init() {
	go func() {
		// 处理goroutine数量
		goroutineName := "rpcClient.heartBeat"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		for {
			time.Sleep(30 * time.Second)

			if clientObj.isLogin {
				updateOnlineCount()
			}
		}
	}()
}

// 启动客户端（连接ChatServerCenter）
func Start() {
	starting := true
	address := config.GetBaseConfig().ChatCenterAddress

	connect := func() {
		defer func() {
			if starting {
				starting = false
			}
		}()

		msg := fmt.Sprintf("rpcClient begins to connect ChatServerCenter: %s...", address)
		fmt.Println(msg)
		logUtil.InfoLog(msg)

		// 连接指定的端口
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			if starting {
				panic(fmt.Errorf("Dial ChatServerCenter Error: %s", err))
			} else {
				return
			}
		}

		msg = fmt.Sprintf("Connect to ChatServerCenter. (local address: %s)", conn.LocalAddr())
		fmt.Println(msg)
		logUtil.InfoLog(msg)

		/*
			1、初始化
			2、通知其它关注者
			3、start
		*/
		clientObj.initialize(conn)
		startSuccessObj.Notify()
		clientObj.start(conn)
	}

	for {
		if clientObj.connected() == false {
			connect()
			time.Sleep(time.Second)
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
