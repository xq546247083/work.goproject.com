package server_gs

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/intAndBytesUtil"
	"work.goproject.com/goutil/timeUtil"
)

const (
	// 包头的长度
	con_HEADER_LENGTH = 4
)

var (
	// 字节的大小端顺序
	byterOrder = binary.LittleEndian

	// 客户端失效的秒数
	con_CLIENT_EXPIRE_SECONDS int64 = 300
)

// 定义客户端对象，以实现对客户端连接的封装
type client struct {
	// 客户端连接对象
	conn net.Conn

	// 接收到的消息内容
	receiveData []byte

	// 活跃时间
	activeTime int64
}

// 获取远程地址（IP_Port）
func (this *client) getRemoteAddr() string {
	items := strings.Split(this.conn.RemoteAddr().String(), ":")

	return fmt.Sprintf("%s_%s", items[0], items[1])
}

// 获取远程地址（IP）
func (this *client) getRemoteShortAddr() string {
	items := strings.Split(this.conn.RemoteAddr().String(), ":")

	return items[0]
}

// 获取有效的消息
// 返回值：
// 消息内容
// 是否含有有效数据
func (this *client) getReceiveData() (message []byte, exists bool) {
	// 判断是否包含头部信息
	if len(this.receiveData) < con_HEADER_LENGTH {
		return
	}

	// 获取头部信息
	header := this.receiveData[:con_HEADER_LENGTH]

	// 将头部数据转换为内部的长度
	contentLength := intAndBytesUtil.BytesToInt32(header, byterOrder)

	// 约定contentLength = 0为心跳包
	if contentLength == 0 {
		// 将对应的数据截断，以得到新的内容
		this.receiveData = this.receiveData[con_HEADER_LENGTH:]
		return
	}

	// 判断长度是否满足
	if len(this.receiveData) < con_HEADER_LENGTH+int(contentLength) {
		return
	}

	// 提取消息内容
	message = this.receiveData[con_HEADER_LENGTH : con_HEADER_LENGTH+contentLength]

	// 将对应的数据截断，以得到新的数据
	this.receiveData = this.receiveData[con_HEADER_LENGTH+contentLength:]

	// 存在合理的数据
	exists = true

	return
}

// 判断客户端是否超时
// 返回值：是否超时
func (this *client) expired() bool {
	return time.Now().Unix() > this.activeTime+con_CLIENT_EXPIRE_SECONDS
}

// 客户端退出
// 返回值：无
func (this *client) quit() {
	this.conn.Close()
	unregisterClient(this)
}

// 格式化
func (this *client) String() string {
	return fmt.Sprintf("{RemoteAddr:%s, activeTime:%s}", this.getRemoteAddr(), timeUtil.Format(time.Unix(this.activeTime, 0), "yyyy-MM-dd HH:mm:ss"))
}

// 启动客户端对象，开始接收和发送数据
func (this *client) start() {
	go this.handleReceiveData()
}

// 处理从连接收到的数据
func (this *client) handleReceiveData() {
	goroutineMgr.MonitorZero("server_gs.handleReceiveData")
	defer goroutineMgr.ReleaseMonitor("server_gs.handleReceiveData")
	defer this.quit()

	for {
		// 先读取数据，每次读取1024个字节
		readBytes := make([]byte, 1024)

		// Read方法会阻塞，所以不用考虑异步的方式
		n, err := this.conn.Read(readBytes)
		if err != nil {
			break
		}

		// 将读取到的数据追加到已获得的数据的末尾，并更新activeTime
		this.receiveData = append(this.receiveData, readBytes[:n]...)
		atomic.StoreInt64(&this.activeTime, time.Now().Unix())

		// 处理收到的数据
		for {
			// 获取有效的消息
			message, exists := this.getReceiveData()
			if !exists {
				break
			}

			handleRequest(this, message)
		}
	}
}

// 新建客户端对象
// conn：连接对象
// 返回值：客户端对象的指针
func newClient(_conn net.Conn) *client {
	return &client{
		conn:        _conn,
		receiveData: make([]byte, 0, 1024),
		activeTime:  time.Now().Unix(),
	}
}
