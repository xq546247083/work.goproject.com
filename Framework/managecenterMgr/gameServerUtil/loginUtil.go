package gameServerUtil

import (
	"encoding/json"
	"fmt"
	"strings"

	. "work.goproject.com/Framework/managecenterModel/returnObject"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
)

// 登陆助手类
type LoginUtil struct{}

// 验证登陆信息
// url:验证地址
// partnerId:合作商Id
// userId:合作商用户Id
// loginInfo:登陆信息
// 返回值:
// 成功与否
// 错误对象
func (this *LoginUtil) CheckLoginInfo(url string, partnerId int32, userId, loginInfo string) (success bool, err error) {
	// 验证用户合法性
	loginItemList := strings.Split(loginInfo, "_")
	if len(loginItemList) != 2 {
		err = fmt.Errorf("CheckLoginInfo Failed. partnerId:%d, userId:%s, loginInfo:%s", partnerId, userId, loginInfo)
		return
	}

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["UserId"] = userId
	postDict["LoginKey"] = loginItemList[0]

	// 去LoginServer验证
	var returnBytes []byte
	if returnBytes, err = webUtil.PostWebData(url, postDict, nil); err != nil {
		err = fmt.Errorf("CheckLoginInfo Failed. partnerId:%d, userId:%s, loginInfo:%s, err:%s", partnerId, userId, loginInfo, err)
		return
	}

	// 解析返回值
	returnObj := new(ReturnObject)
	if err = json.Unmarshal(returnBytes, &returnObj); err != nil {
		err = fmt.Errorf("CheckLoginInfo Failed. partnerId:%d, userId:%s, loginInfo:%s, err:%s", partnerId, userId, loginInfo, err)
		return
	}

	// 判断返回状态是否为成功
	if returnObj.Code != 0 {
		logUtil.ErrorLog(fmt.Sprintf("CheckLoginInfo Failed. partnerId:%d, userId:%s, loginInfo:%s, Code:%d, Message:%s", partnerId, userId, loginInfo, returnObj.Code, returnObj.Message))
		return
	}

	success = true

	return
}

// ------------------类型定义和业务逻辑的分隔符-------------------------

var (
	LoginUtilObj = new(LoginUtil)
)
