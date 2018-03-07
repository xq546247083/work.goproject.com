package server_gs

import (
	"encoding/json"

	"work.goproject.com/Chat/ChatServer/src/bll/player"
	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/zlibUtil"
)

// 处理客户端请求
// clientObj：对应的客户端对象
// message：请求内容字节数组(json格式)
// 返回值：无
func handleRequest(clientObj *client, message []byte) {
	// 先进行解压缩，再进行处理
	content, err := zlibUtil.Decompress(message)
	if err != nil {
		logUtil.ErrorLog("server_gs zlib解压缩错误，错误信息为：%s", err)
		return
	}

	// 更新信息
	increaseSize(len(content))

	// 将接收到的数据反序列化为GSPlayer对象
	playerObj := new(Player)
	if err := json.Unmarshal(content, playerObj); err != nil {
		logUtil.ErrorLog("server_gs 反序列化%s出错，错误信息为：%s", string(content), err)
		return
	}

	debugUtil.Printf("收到来自GS的数据：%v\n", playerObj)

	// 处理数据
	go player.UpdateFromGS(playerObj)
}
