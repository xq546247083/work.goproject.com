package player

import (
	"work.goproject.com/Chat/ChatServer/src/clientMgr"
	"work.goproject.com/Chat/ChatServer/src/rpcClient"
)

func init() {
	rpcClient.RegisterGetPlayerCountHandler(getPlayerCount)
	clientMgr.RegisterGetPlayerHandler(GetPlayer)
}
