package rpc

import (
	"encoding/binary"
	"work.goproject.com/goutil/intAndBytesUtil"
	"net"
)

const (
	// 包头的长度
	con_HEADER_LENGTH = 4
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
	content []byte
}

// 追加内容
// content：新的内容
// 返回值：无
func (clientObj *client) appendContent(content []byte) {
	clientObj.content = append(clientObj.content, content...)
}

// 获取有效的消息
// 返回值：
// 消息内容
// 是否含有有效数据
func (clientObj *client) getValidMessage() ([]byte, bool) {
	// 判断是否包含头部信息
	if len(clientObj.content) < con_HEADER_LENGTH {
		return nil, false
	}

	// 获取头部信息
	header := clientObj.content[:con_HEADER_LENGTH]

	// 将头部数据转换为内部的长度
	contentLength := intAndBytesUtil.BytesToInt32(header, byterOrder)

	// 判断长度是否满足
	if len(clientObj.content) < con_HEADER_LENGTH+int(contentLength) {
		return nil, false
	}

	// 提取消息内容
	content := clientObj.content[con_HEADER_LENGTH : con_HEADER_LENGTH+contentLength]

	// 将对应的数据截断，以得到新的数据
	clientObj.content = clientObj.content[con_HEADER_LENGTH+contentLength:]

	// 判断是否为心跳包，如果是心跳包，则不解析，直接返回
	if contentLength == 0 || len(content) == 0 {
		return nil, false
	}

	return content, true
}

// 发送字节数组消息
// message：待发送的字节数组
func (clientObj *client) sendByteMessage(message []byte) {
	// 获得数组的长度
	contentLength := len(message)

	// 将长度转化为字节数组
	header := intAndBytesUtil.Int32ToBytes(int32(contentLength), byterOrder)

	// 将头部与内容组合在一起
	message = append(header, message...)

	// 发送消息
	clientObj.conn.Write(message)
}

// 发送字符串消息
// s：待发送的字符串
func (clientObj *client) sendStringMessage(s string) {
	clientObj.sendByteMessage([]byte(s))
}

// 发送心跳包信息
func (clientObj *client) sendHeartBeatMessage() {
	clientObj.sendByteMessage([]byte{})
}

// 新建客户端对象
// conn：连接对象
// 返回值：客户端对象的指针
func newClient(_conn net.Conn) *client {
	return &client{
		conn:    _conn,
		content: make([]byte, 0, 1024),
	}
}
