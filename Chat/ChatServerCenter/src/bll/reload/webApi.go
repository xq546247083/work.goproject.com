package reload

import (
	"work.goproject.com/Chat/ChatServerCenter/src/server_http"
	"work.goproject.com/Chat/ChatServerCenter/src/server_tcp"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/Framework/reloadMgr"
)

func init() {
	server_http.RegisterHandler("/API/config/reload", reloadHandler)
}

func reloadHandler(context *server_http.Context) *CenterResponseObject {
	responseObj := NewCenterResponseObject()

	if errList := reloadMgr.Reload(); errList != nil {
		return responseObj.SetResultStatus(DataError)
	}

	// 推送给ChatServer
	server_tcp.ForwardObjectChannel <- NewForwardObject_Reload()

	return responseObj
}
