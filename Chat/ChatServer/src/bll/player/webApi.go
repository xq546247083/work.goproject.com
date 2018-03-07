package player

import (
	"work.goproject.com/Chat/ChatServer/src/server_http"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func init() {
	server_http.RegisterHandler("/API/player/login", playerLoginHandler)
}

func playerLoginHandler(context *server_http.Context) *ServerResponseObject {
	responseObj := NewServerResponseObject()

	var gsPlayerObj *Player
	if exists, err := context.Unmarshal("player", &gsPlayerObj, true); err != nil || !exists {
		return responseObj.SetResultStatus(DataError)
	}

	UpdateFromGS(gsPlayerObj)

	return responseObj
}
