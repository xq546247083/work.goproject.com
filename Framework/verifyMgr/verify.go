package verifyMgr

import (
	"fmt"
	"sync"

	"work.goproject.com/Framework/managecenterMgr"
	"work.goproject.com/goutil/webUtil"
)

var (
	verifyUrlMap          = make(map[string]VerifyType)
	mutex                 sync.RWMutex
	con_SuccessStatusCode = 200
)

// 初始化需要验证的Url及对应的验证方式
func Init(url string, _type VerifyType) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := verifyUrlMap[url]; exists {
		panic(fmt.Errorf("%s已经存在，请检查", url))
	}

	verifyUrlMap[url] = _type
}

// 验证
func Verify() (errList []error) {
	errList1 := verifyInternal()
	errList2 := verifyServerGroup()

	if len(errList1) > 0 {
		errList = append(errList, errList1...)
	}
	if len(errList2) > 0 {
		errList = append(errList, errList2...)
	}

	return
}

// 验证内部的url
func verifyInternal() (errList []error) {
	// 验证单独指定的url
	mutex.RLock()
	defer mutex.RUnlock()

	for url, _type := range verifyUrlMap {
		switch _type {
		case Con_Get:
			if statusCode, _, err := webUtil.GetWebData2(url, nil, nil, nil); statusCode != con_SuccessStatusCode || err != nil {
				errList = append(errList, fmt.Errorf("access %s failed. StatusCode:%d, err:%s", url, statusCode, err))
			}
		case Con_Post:
			if statusCode, _, err := webUtil.PostByteData2(url, nil, nil, nil); statusCode != con_SuccessStatusCode || err != nil {
				errList = append(errList, fmt.Errorf("access %s failed. StatusCode:%d, err:%s", url, statusCode, err))
			}
		default:
			errList = append(errList, fmt.Errorf("the type of %s is wrong, not it's %d", url, _type))
		}
	}

	return
}

// 验证服务器组的url
func verifyServerGroup() (errList []error) {
	// 验证配置在ManageCenter中的所有ServerGroup
	serverGroupList := managecenterMgr.GetServerGroupList()
	chList := make([]chan bool, 0, len(serverGroupList))

	for _, item := range serverGroupList {
		url := item.GetGSCallbackUrl("Verify.ashx")
		ch := make(chan bool)
		chList = append(chList, ch)
		go func(_url string, _ch chan bool) {
			fmt.Printf("VerifyUrl:%s\n", _url)
			if statusCode, _, err := webUtil.PostByteData2(_url, nil, nil, nil); statusCode != con_SuccessStatusCode || err != nil {
				errList = append(errList, fmt.Errorf("access %s failed. StatusCode:%d, err:%s", _url, statusCode, err))
			}
			_ch <- true
		}(url, ch)
	}

	// 等待所有的ch返回
	for _, ch := range chList {
		<-ch
	}

	return
}
