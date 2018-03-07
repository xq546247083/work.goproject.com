package server_webSocket

import (
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"work.goproject.com/Chat/ChatServer/src/clientMgr"
	"work.goproject.com/Chat/ChatServer/src/config"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/goroutineMgr"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/timeUtil"
	"work.goproject.com/goutil/zlibUtil"
)

const (
	// 消息每次读取的上限值
	con_MaxMessageSize = 1024

	// 客户端失效的秒数
	con_CLIENT_EXPIRE_SECONDS = 20
)

var (
	// 全局客户端的id，从1开始进行自增
	globalClientId int32 = 0

	// 字节的大小端顺序
	byterOrder = binary.LittleEndian
)

// 实现IConnection接口的所有方法
// 定义客户端对象，以实现对客户端连接的封装
type Client struct {
	// 唯一标识
	id int32

	// The websocket connection.
	conn *websocket.Conn

	// 接收到的消息内容
	receiveData []byte

	// 待发送的数据
	sendData []*ServerResponseObject

	// 连接是否关闭
	closed bool

	// 锁对象（用于控制对sendDatap的并发访问；receiveData不需要，因为是同步访问）
	mutex sync.Mutex

	// 玩家Id
	playerId string

	// 上次活跃时间
	activeTime int64
}

// 获取唯一标识
func (this *Client) GetId() int32 {
	return this.id
}

// 获取玩家Id
// 返回值：
// 玩家Id
func (this *Client) GetPlayerId() string {
	return this.playerId
}

// 获取远程地址（IP_Port）
func (this *Client) getRemoteAddr() string {
	items := strings.Split(this.conn.RemoteAddr().String(), ":")

	return fmt.Sprintf("%s_%s", items[0], items[1])
}

// 获取远程地址（IP）
func (this *Client) getRemoteShortAddr() string {
	items := strings.Split(this.conn.RemoteAddr().String(), ":")

	return items[0]
}

// 获取待发送的数据
// 返回值：
// 待发送数据项
// 是否含有有效数据
func (this *Client) getSendData() (responseObj *ServerResponseObject, exists bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 如果没有数据则直接返回
	if len(this.sendData) == 0 {
		return
	}

	// 取出第一条数据,并为返回值赋值
	responseObj = this.sendData[0]
	exists = true

	// 删除已经取出的数据
	this.sendData = this.sendData[1:]

	return
}

// 追加发送的数据
// sendDataItemObj:待发送数据项
// 返回值：无
func (this *Client) AppendSendData(responseObj *ServerResponseObject) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.sendData = append(this.sendData, responseObj)
}

// 向客户端发送消息
// message：最终要发送的二进制流
func (this *Client) sendMessage(responseObj *ServerResponseObject) error {
	beforeTime := time.Now().Unix()

	// 序列化发送的数据
	content, _ := json.Marshal(responseObj)

	// 进行zlib压缩
	if config.GetBaseConfig().IfCompressData {
		content, _ = zlibUtil.Compress(content, zlib.DefaultCompression)
	}

	// 发送消息
	if err := this.sendMessageToClient(content); err != nil {
		return err
	}

	// 如果发送的时间超过3秒，则记录下来
	if time.Now().Unix()-beforeTime > 3 {
		logUtil.WarnLog("消息Size:%d, UseTime:%d", len(content), time.Now().Unix()-beforeTime)
	}

	return nil
}

func (this *Client) sendMessageToClient(message []byte) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if err := this.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
		return err
	}

	return nil
}

// 判断客户端是否超时（超过300秒不活跃算作超时）
// 返回值：是否超时
func (this *Client) Expired() bool {
	return time.Now().Unix() > this.activeTime+con_CLIENT_EXPIRE_SECONDS
}

// 格式化
func (this *Client) String() string {
	return fmt.Sprintf("{Id:%d, RemoteAddr:%s, activeTime:%s, playerId:%s}", this.id, this.getRemoteAddr(), timeUtil.Format(time.Unix(this.activeTime, 0), "yyyy-MM-dd HH:mm:ss"), this.playerId)
}

// 玩家登陆
// playerId：玩家Id
// 返回值：无
func (this *Client) PlayerLogin(playerId string) {
	this.playerId = playerId
}

// 客户端退出
// 返回值：无
func (this *Client) Quit() {
	this.conn.Close()
	this.closed = true
	// 注销客户端连接，并从缓存中移除
	clientMgr.UnregisterClient(this)
}

// -----------------------------实现IConnection接口方法end-------------------------

// 客户端启动函数
func (this *Client) start() {
	go this.handleReceiveData()
	go this.handleSendData()
}

// 处理从连接收到的数据
func (this *Client) handleReceiveData() {
	// 处理goroutine数量
	goroutineName := "server_webSocket.handleSendData"
	goroutineMgr.MonitorZero(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)
	defer this.Quit()

	// 无限循环，不断地读取数据，解析数据，处理数据
	for {
		if this.closed {
			break
		}

		_, message, err := this.conn.ReadMessage()
		if err != nil {
			break
		}

		// 更新activeTime
		atomic.StoreInt64(&this.activeTime, time.Now().Unix())

		// 约定len(message) == 0,为心跳请求
		if len(message) == 0 {
			if err := this.sendMessageToClient([]byte{}); err != nil {
				return
			}
			continue
		}

		clientMgr.HandleRequest(this, message)
	}
}

// 处理向客户端发送的数据
func (this *Client) handleSendData() {
	// 处理goroutine数量
	goroutineName := "server_webSocket.handleSendData"
	goroutineMgr.MonitorZero(goroutineName)
	defer goroutineMgr.ReleaseMonitor(goroutineName)
	defer this.Quit()

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
// conn：连接对象
// 返回值：客户端对象的指针
func newClient(conn *websocket.Conn) *Client {
	conn.SetReadLimit(con_MaxMessageSize)

	// 获得自增的id值
	getIncrementId := func() int32 {
		atomic.AddInt32(&globalClientId, 1)
		return globalClientId
	}

	return &Client{
		id:          getIncrementId(),
		conn:        conn,
		receiveData: make([]byte, 0, 1024),
		sendData:    make([]*ServerResponseObject, 0, 16),
		activeTime:  time.Now().Unix(),
		playerId:    "",
	}
}
