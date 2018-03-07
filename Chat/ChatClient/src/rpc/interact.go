package rpc

import (
	"fmt"
	"work.goproject.com/goutil/logUtil"
	"strconv"
	"strings"
)

func Interaction(ch chan int) {
	// 处理内部未处理的异常，以免导致主线程退出，从而导致系统崩溃
	defer func() {
		if r := recover(); r != nil {
			logUtil.Log(fmt.Sprintf("通过recover捕捉到的未处理异常：%v", r), logUtil.Error, true)
		}

		ch <- 1
	}()

	// 源源不断地从键盘输入数据
	fmt.Println("请输入要发送的信息。如果要退出则输入q")
	for {
		var input string
		n, err := fmt.Scan(&input)
		if err != nil {
			fmt.Printf("input error:%s\n", err)
			break
		}

		if n == 0 {
			fmt.Println("你输入的数据为空，请重新输入。如果要退出则输入q")
			continue
		}

		// 如果选择退出
		if input == "q" {
			clientObj.conn.Close()
			break
		}

		// 发送数据
		param := assembleMessageParam(input)
		if param != nil {
			request(param)
		}
	}
}

func assembleMessageParam(message string) map[string]interface{} {
	methodName := ""
	parameters := make([]interface{}, 0, 4)

	// 根据message来判断是世界频道，还是公会频道
	msgList := strings.Split(message, ":")
	fmt.Printf("message:%s, msgList:%v\n", message, msgList)
	if len(msgList) == 1 && msgList[0] != "Logout" && msgList[0] != "GetHistoryInfo" {
		fmt.Println("输入参数错误，请重新输入")
		return nil
	}

	methodName = msgList[0]
	switch methodName {
	case "SendMessage":
		switch {
		case msgList[1] == "Private":
			parameters = append(parameters, msgList[1], msgList[2], "voice", msgList[3])
		case msgList[1] != "Private":
			parameters = append(parameters, msgList[1], msgList[2], "voice", "")
		}
	case "GetHistory":
		messageId, _ := strconv.Atoi(msgList[2])
		count, _ := strconv.Atoi(msgList[3])
		switch {
		case msgList[1] == "Private":
			parameters = append(parameters, msgList[1], messageId, count, msgList[4])
		case msgList[1] != "Private":
			parameters = append(parameters, msgList[1], messageId, count, "")
		}
	case "DeletePrivateHistory":
		parameters = append(parameters, msgList[1])
	}

	requestMap := make(map[string]interface{})
	requestMap["MethodName"] = methodName
	requestMap["Parameters"] = parameters

	return requestMap
}
