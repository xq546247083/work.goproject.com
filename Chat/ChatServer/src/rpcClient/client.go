package rpcClient

import (
	"encoding/binary"
	"net"
	"sync"

	"work.goproject.com/goutil/intAndBytesUtil"
	"work.goproject.com/goutil/logUtil"
)

const (
	// 包头的长度
	con_HEADER_LENGTH = 4

	// 定义请求、响应数据的前缀的长度
	con_ID_LENGTH = 4
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

	// 是否已经登陆
	isLogin bool

	mutex sync.Mutex
}

func (this *client) connected() bool {
	return this.conn != nil
}

func (this *client) initialize(conn net.Conn) {
	this.conn = conn
	this.receiveData = make([]byte, 0, 1024)
	this.isLogin = false
}

func (this *client) start(conn net.Conn) {
	defer this.quit()

	// 死循环，不断地读取数据，解析数据，发送数据
	for {
		// 先读取数据，每次读取1024个字节
		readBytes := make([]byte, 1024)

		// Read方法会阻塞，所以不用考虑异步的方式
		n, err := this.conn.Read(readBytes)
		if err != nil {
			break
		}

		// 将读取到的数据追加到已获得的数据的末尾
		this.receiveData = append(this.receiveData, readBytes[:n]...)

		// 处理数据
		for {
			id, content, ok := this.getValidMessage()
			if !ok {
				break
			}

			handleMessage(id, content)
		}
	}
}

func (this *client) login() {
	this.isLogin = true
}

func (this *client) quit() {
	this.conn.Close()
	this.conn = nil
	this.receiveData = make([]byte, 0, 1024)
	this.isLogin = false
}

// 获取有效的消息
// 返回值：
// 消息对应客户端的唯一标识
// 消息内容
// 是否含有有效数据
func (this *client) getValidMessage() (id int32, message []byte, exists bool) {
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
	content := this.receiveData[con_HEADER_LENGTH : con_HEADER_LENGTH+contentLength]

	// 将对应的数据截断，以得到新的数据
	this.receiveData = this.receiveData[con_HEADER_LENGTH+contentLength:]

	// 判断内容的长度是否足够
	if len(content) < con_ID_LENGTH {
		logUtil.ErrorLog("内容数据不正确；con_ID_LENGTH=%d,len(content)=%d", con_ID_LENGTH, len(content))
		return
	}

	// 截取内容的前4位
	idBytes, content := content[:con_ID_LENGTH], content[con_ID_LENGTH:]

	// 组装返回值
	id = intAndBytesUtil.BytesToInt32(idBytes, byterOrder)
	message = content
	exists = true

	return
}

// 发送字节数组消息
// id：需要添加到b前发送的数据
// message：待发送的字节数组
func (this *client) sendByteMessage(id int32, message []byte) error {
	idBytes := intAndBytesUtil.Int32ToBytes(id, byterOrder)

	// 将idByte和b合并
	message = append(idBytes, message...)

	// 获得数组的长度
	contentLength := len(message)

	// 将长度转化为字节数组
	header := intAndBytesUtil.Int32ToBytes(int32(contentLength), byterOrder)

	// 将头部与内容组合在一起
	message = append(header, message...)

	// 发送消息
	_, err := this.conn.Write(message)

	return err
}

// 新建客户端对象
// conn：连接对象
// 返回值：客户端对象的指针
func newClient() *client {
	return &client{}
}
