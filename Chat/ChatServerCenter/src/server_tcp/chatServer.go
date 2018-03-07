package server_tcp

import (
	"fmt"
	"sort"

	"work.goproject.com/Chat/ChatServerCenter/src/server_http"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func init() {
	server_http.RegisterHandler("/API/server/getAvailable", getAvailableHandler)
}

// 为客户端提供服务的chatServer对象
type chatServer struct {
	// ChatServer监听地址
	chatServerAddress string

	// ChatServer为GS提供的监听地址
	gameServerAddress string

	// ChatServer为为GS提供的Web监听地址
	gameServerWebAddress string

	// 连上ChatServer的客户端数量
	clientCount int

	// 连上ChatServer的玩家数量
	playerCount int
}

func (this *chatServer) SortByClientCountAsc(target *chatServer) bool {
	return this.clientCount < target.clientCount
}

// 更新客户端数量
// clientCount：客户端数量
// playerCount：玩家数量
// 返回值：无
func (this *chatServer) updateOnlineCount(clientCount int, playerCount int) {
	this.clientCount = clientCount
	this.playerCount = playerCount
}

// 格式化字符串
// 返回值：
// 格式化后的字符串
func (this *chatServer) String() string {
	return fmt.Sprintf("chatServerAddress:%s, gameServerAddress:%s, gameServerWebAddress:%s, ClientCount:%d, PlayerCount:%d",
		this.chatServerAddress, this.gameServerAddress, this.gameServerWebAddress, this.clientCount, this.playerCount)
}

// 创建新的chatServer对象
func newChatServer(chatServerAddress, gameServerAddress, gameServerWebAddress string) *chatServer {
	return &chatServer{
		chatServerAddress:    chatServerAddress,
		gameServerAddress:    gameServerAddress,
		gameServerWebAddress: gameServerWebAddress,
		clientCount:          0,
		playerCount:          0,
	}
}

func getAvailableHandler(context *server_http.Context) (responseObj *CenterResponseObject) {
	responseObj = NewCenterResponseObject()

	currAddress, exists := context.GetFormValue("ChatServerAddress")
	if !exists {
		return responseObj.SetResultStatus(ParamNotMatch)
	}

	// 获得Client列表
	list := make([]*chatServer, 0, 4)
	for _, item := range getClientList() {
		if item.chatServer != nil {
			list = append(list, item.chatServer)
		}
	}

	// 如果没有可用的服务器，则返回空
	if len(list) == 0 {
		return responseObj.SetResultStatus(NoAvailableServer)
	}

	var targetChatServer *chatServer
	defer func() {
		data := make(map[string]interface{})
		data["ChatServerAddress"] = targetChatServer.chatServerAddress
		data["GameServerAddress"] = targetChatServer.gameServerAddress
		data["GameServerWebAddress"] = targetChatServer.gameServerWebAddress
		responseObj.SetData(data)
	}()

	// 如果当前服务器所连接的聊天服务器不为空，则判断是否有与当前地址相同的服务器，如果有则直接返回
	if currAddress != "" {
		for _, item := range list {
			if item.chatServerAddress == currAddress {
				targetChatServer = item
			}
		}
	}

	// 运行到此处，表明当前服务器对应的聊天服务器为空；或者该服务器不再存在；则选择客户端连接最少的
	sort.Slice(list, func(i, j int) bool {
		return list[i].SortByClientCountAsc(list[j])
	})
	targetChatServer = list[0]

	return
}
