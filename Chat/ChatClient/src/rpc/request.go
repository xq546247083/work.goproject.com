package rpc

import (
	"encoding/json"
	"fmt"
	"work.goproject.com/goutil/logUtil"
)

func request(requestMap map[string]interface{}) {
	message, err := json.Marshal(requestMap)
	if err != nil {
		logUtil.Log(fmt.Sprintf("序列化请求数据%v出错", requestMap), logUtil.Error, true)
	} else {
		clientObj.sendByteMessage(message)
	}
}
