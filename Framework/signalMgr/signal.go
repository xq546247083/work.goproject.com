package signalMgr

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"work.goproject.com/Framework/exitMgr"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

// Start ...启动信号管理器
func Start() {
	go func() {
		goroutineName := "signalMgr.Start"
		goroutineMgr.Monitor(goroutineName)
		defer goroutineMgr.ReleaseMonitor(goroutineName)

		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

		for {
			// 准备接收信息
			sig := <-sigs

			// 输出信号
			debugUtil.Println("sig:", sig)

			if sig == syscall.SIGHUP {
				logUtil.NormalLog("收到重启的信号，准备重新加载配置", logUtil.Info)

				// 重新加载
				errList := reloadMgr.Reload()
				for _, err := range errList {
					logUtil.NormalLog(fmt.Sprintf("重启失败，错误信息为:%s", err), logUtil.Error)
				}

				logUtil.NormalLog("收到重启的信号，重新加载配置完成", logUtil.Info)
			} else {
				logUtil.NormalLog("收到退出程序的信号，开始退出……", logUtil.Info)

				// 调用退出的方法
				exitMgr.Exit()

				logUtil.NormalLog("收到退出程序的信号，退出完成……", logUtil.Info)

				// 一旦收到信号，则表明管理员希望退出程序，则先保存信息，然后退出
				os.Exit(0)
			}
		}
	}()
}
