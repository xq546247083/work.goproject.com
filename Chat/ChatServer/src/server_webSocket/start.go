package server_webSocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"work.goproject.com/Chat/ChatServer/src/clientMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logUtil.ErrorLog("websocket.handleConn获取连接出错，err:%v", err)
		return
	}

	// 创建客户端对象
	clientObj := newClient(conn)
	clientObj.start()
	clientMgr.RegisterClient(clientObj)

	debugUtil.Printf("收到连接请求:remoteAdd:%s\n", conn.RemoteAddr())
}

// 启动服务器
func Start(wg *sync.WaitGroup, address string) {
	defer wg.Done()

	msg := fmt.Sprintf("server_websocket begins to listen on:%s...", address)
	fmt.Println(msg)
	logUtil.InfoLog(msg)

	http.HandleFunc("/", handleConn)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(fmt.Sprintf("server_websocket.ListenAndServe, err:%v", err))
	}
}

// 启动https服务器
func Start2(wg *sync.WaitGroup, address,crt,key string) {
	defer wg.Done()

	msg := fmt.Sprintf("server_websocket begins to listen TLS on:%s...", address)
	fmt.Println(msg)
	logUtil.InfoLog(msg)

	http.HandleFunc("/", handleConn)
	err := http.ListenAndServeTLS(address,crt,key, nil)
	if err != nil {
		panic(fmt.Sprintf("server_websocket.ListenAndServeTLS, err:%v", err))
	}
}
