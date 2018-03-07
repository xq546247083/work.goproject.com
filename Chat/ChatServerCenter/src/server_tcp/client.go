package server_tcp

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/intAndBytesUtil"
	"work.goproject.com/goutil/timeUtil"
)

const (
	// 包头的长度
	con_HEADER_LENGTH = 4

	// 定义请求、响应数据的前缀的长度
	con_ID_LENGTH = 4

	// 客户端失效的秒数
	con_CLIENT_EXPIRE_SECONDS int64 = 300
)

var (
	// 字节的大小端顺序
	byterOrder = binary.LittleEndian
)

// 定义客户端对象，以实现对客户端连接的封装
type client struct {
	// 客户端连接对象
	conn net.Conn

	// 接收到的消息内容
	receiveData []byte

	// 待发送的数据
	sendData []*sendDataItem

	// 连接是否关闭(通过此字段来协调receiveData和sendData方法)
	closed bool

	// 活跃时间
	activeTime int64

	// 为客户端提供服务的ChatServer对象
	*chatServer

	// 锁对象（用于控制对sendDatap的并发访问；receiveData不需要，因为是同步访问）
	mutex sync.Mutex
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

// 获取待发送的数据
// 返回值：
// 待发送数据项
// 是否含有有效数据
func (this *client) getSendData() (sendDataItemObj *sendDataItem, exists bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 如果没有数据则直接返回
	if len(this.sendData) == 0 {
		return
	}

	// 取出第一条数据,并为返回值赋值
	sendDataItemObj = this.sendData[0]
	exists = true

	// 删除已经取出的数据
	this.sendData = this.sendData[1:]

	return
}

// 追加发送的数据
// sendDataItemObj:待发送数据项
// 返回值：无
func (this *client) appendSendData(sendDataItemObj *sendDataItem) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.sendData = append(this.sendData, sendDataItemObj)
}

// 获取接收到的数据
// 返回值：
// 消息对应客户端的唯一标识
// 消息内容
// 是否含有有效数据
func (this *client) getReceiveData() (id int32, message []byte, exists bool) {
	// 判断是否包含头部信息
	if len(this.receiveData) < con_HEADER_LENGTH {
		return
	}

	// 获取头部信息
	header := this.receiveData[:con_HEADER_LENGTH]

	// 将头部数据转换为内容的长度
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
	content := this.receiveData[con_HEADER_LENGTH : con_HEADER_LENGTH+contentLength]

	// 将对应的数据截断，以得到新的内容
	this.receiveData = this.receiveData[con_HEADER_LENGTH+contentLength:]

	// 将内容分隔为2部分
	idBytes, content := content[:con_ID_LENGTH], content[con_ID_LENGTH:]

	// 提取id、message
	id = intAndBytesUtil.BytesToInt32(idBytes, byterOrder)
	message = content

	// 存在合理的数据
	exists = true

	return
}

// 判断客户端是否超时
// 返回值：
// 是否超时
func (this *client) expired() bool {
	return time.Now().Unix() > this.activeTime+con_CLIENT_EXPIRE_SECONDS
}

// 发送消息
// id：需要添加到内容前的数据
// sendDataItemObj：待发送数据项
// 返回值：
// 错误对象
func (this *client) sendMessage(sendDataItemObj *sendDataItem) error {
	idBytes := intAndBytesUtil.Int32ToBytes(sendDataItemObj.id, byterOrder)

	// 将idByte和内容合并
	content := append(idBytes, sendDataItemObj.data...)

	// 获得数组的长度
	contentLength := len(content)

	// 将长度转化为字节数组
	header := intAndBytesUtil.Int32ToBytes(int32(contentLength), byterOrder)

	// 将头部与内容组合在一起
	message := append(header, content...)

	// 发送消息
	_, err := this.conn.Write(message)

	return err
}

// 登录
func (this *client) login(localAddress, publicAddress, gameServerAddress string) {
	this.chatServer = newChatServer(localAddress, publicAddress, gameServerAddress)
}

// 退出
// 返回值：无
func (this *client) quit() {
	this.conn.Close()
	this.closed = true
	this.chatServer = nil
	unregisterClient(this)
}

// 格式化客户端对象
// 返回值：
// 格式化的字符串
func (this *client) String() string {
	if this.chatServer == nil {
		return fmt.Sprintf("{RemoteAddr:%s, ActiveTime:%s, ChatServer尚未登录}", this.getRemoteAddr(), timeUtil.Format(time.Unix(this.activeTime, 0), "yyyy-MM-dd HH:mm:ss"))
	} else {
		return fmt.Sprintf("{RemoteAddr:%s, ActiveTime:%s, ChatServer:%s}", this.getRemoteAddr(), timeUtil.Format(time.Unix(this.activeTime, 0), "yyyy-MM-dd HH:mm:ss"), this.chatServer.String())
	}
}

// 启动客户端对象，开始接收和发送数据
func (this *client) start() {
	go this.handleReceiveData()
	go this.handleSendData()
}

// 处理客户端收到的数据
// clientObj：客户端对象
func (this *client) handleReceiveData() {
	// 处理goroutine数量
	goroutineName := "server_tcp.handleReceiveData"
	goroutineMgr.MonitorZero(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)
	defer this.quit()

	// 无限循环，不断地读取数据，解析数据，处理数据
	for {
		if this.closed {
			break
		}

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

		// 处理数据
		for {
			// 获取有效的消息
			id, message, exists := this.getReceiveData()
			if !exists {
				break
			}

			handleRequest(this, id, message)
		}
	}
}

// 处理需要客户端发送的数据
func (this *client) handleSendData() {
	// 处理goroutine数量
	goroutineName := "server_tcp.handleSendData"
	goroutineMgr.MonitorZero(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)
	defer this.quit()

	for {
		if this.closed {
			break
		}

		// 如果发送出现错误，表示连接已经断开，则退出方法；
		if sendDataItemObj, exists := this.getSendData(); exists {
			if err := this.sendMessage(sendDataItemObj); err != nil {
				return
			}
		} else {
			time.Sleep(5 * time.Millisecond)
		}
	}
}

// 新建客户端对象
// _conn：连接对象
// 返回值：客户端对象的指针
func newClient(_conn net.Conn) *client {
	return &client{
		conn:        _conn,
		receiveData: make([]byte, 0, 1024),
		sendData:    make([]*sendDataItem, 0, 16),
		closed:      false,
		activeTime:  time.Now().Unix(),
	}
}
