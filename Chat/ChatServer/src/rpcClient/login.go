package rpcClient

import (
	"work.goproject.com/Chat/ChatServer/src/config"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/Framework/initMgr"
	"work.goproject.com/goutil/debugUtil"
)

var (
	// 启动成功的通知通道
	serverStartSuccessName_login = "login"
	serverStartSuccessCh_login   = make(chan bool)

	// 启动成功对象
	LoginSuccessObj = initMgr.NewInitSuccess("rpcClientLogin")
)

func init() {
	// 注册启动成功通知
	startSuccessObj.Register(serverStartSuccessName_login, serverStartSuccessCh_login)

	// 登陆
	go login()
}

func login() {
	// 处理goroutine数量
	goroutineName := "rpcClient.login"
	goroutineMgr.Monitor(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)

	for {
		// 等待客户端连接成功
		<-serverStartSuccessCh_login
		debugUtil.Println("收到来自serverStartSuccessCh的消息")

		params := make([]interface{}, 3, 3)
		params[0] = config.GetBaseConfig().GetPublicChatServerAddress()
		params[1] = config.GetBaseConfig().GetPublicGameServerAddress()
		params[2] = config.GetBaseConfig().GetPublicGameServerWebAddress()

		// 发送Login消息
		debugUtil.Println("Send Login Command")

		request(Con_ChatServerCenter_Method_Login, params, loginCallback)
	}
}

func loginCallback(data interface{}) {
	debugUtil.Println("Login success")
	clientObj.login()
	LoginSuccessObj.Notify()
}
