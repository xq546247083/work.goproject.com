package player

import (
	"strconv"
	"time"

	"work.goproject.com/Chat/ChatServerCenter/src/server_http"
	. "work.goproject.com/Chat/ChatServerModel/src"
)

func init() {
	server_http.RegisterHandler("/API/player/silent", silentHandler)
	server_http.RegisterHandler("/API/player/unSilent", unSilentHandler)
	server_http.RegisterHandler("/API/player/getSilentInfo", getSilentInfoHandler)
}

func silentHandler(context *server_http.Context) *CenterResponseObject {
	responseObj := NewCenterResponseObject()

	// Get param
	id, exists := context.GetFormValue("Id")
	if !exists {
		return responseObj.SetResultStatus(ParamNotMatch)
	}
	minutes_str, exists := context.GetFormValue("Minutes")
	if !exists {
		return responseObj.SetResultStatus(ParamNotMatch)
	}

	// Check the param
	var minutes int
	var err error
	if minutes, err = strconv.Atoi(minutes_str); err != nil {
		return responseObj.SetResultStatus(ParamTypeError)
	}

	// get player
	playerObj, exists, err := get(id)
	if err != nil {
		return responseObj.SetResultStatus(DBError)
	} else if !exists {
		return responseObj.SetResultStatus(PlayerNotExists)
	}

	// update info
	err = updateSilentInfo(playerObj, time.Now().Add(time.Duration(minutes)*time.Minute))
	if err != nil {
		return responseObj.SetResultStatus(DBError)
	}

	return responseObj
}

func unSilentHandler(context *server_http.Context) *CenterResponseObject {
	responseObj := NewCenterResponseObject()

	// Get param
	id, exists := context.GetFormValue("Id")
	if !exists {
		return responseObj.SetResultStatus(ParamNotMatch)
	}

	// get player
	playerObj, exists, err := get(id)
	if err != nil {
		return responseObj.SetResultStatus(DBError)
	} else if !exists {
		return responseObj.SetResultStatus(PlayerNotExists)
	}

	// update info
	err = updateSilentInfo(playerObj, time.Now())
	if err != nil {
		return responseObj.SetResultStatus(DBError)
	}

	return responseObj
}

func getSilentInfoHandler(context *server_http.Context) *CenterResponseObject {
	responseObj := NewCenterResponseObject()

	// Get param
	id, exists := context.GetFormValue("Id")
	if !exists {
		return responseObj.SetResultStatus(ParamNotMatch)
	}

	// get player
	playerObj, exists, err := get(id)
	if err != nil {
		return responseObj.SetResultStatus(DBError)
	} else if !exists {
		return responseObj.SetResultStatus(PlayerNotExists)
	}

	isInSilent, leftMinutes := playerObj.IsInSilent()

	// 组装返回值
	data := make(map[string]interface{})
	data["IsInSilent"] = isInSilent
	data["LeftMinutes"] = leftMinutes

	return responseObj.SetData(data)
}
