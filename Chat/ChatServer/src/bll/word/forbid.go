package word

import (
	"fmt"
	"sync"

	"work.goproject.com/Chat/ChatServer/src/dal"
	"work.goproject.com/Framework/reloadMgr"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/dfaUtil"
)

var (
	forbidWordList []string = make([]string, 0, 1024)
	forbidDFAObj   *dfaUtil.DFAUtil
	forbidMutex    sync.RWMutex
)

func init() {
	if err := reloadForbid(); err != nil {
		panic(fmt.Errorf("初始化屏蔽词列表失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadMgr.RegisterReloadFunc("word.reloadForbid", reloadForbid)
}

// 重新加载屏蔽词列表
func reloadForbid() error {
	var tmpForbidWordList []*ForbidWord
	if err := dal.GetAll(&tmpForbidWordList); err != nil {
		return err
	}

	forbidMutex.Lock()
	defer forbidMutex.Unlock()
	forbidWordList = make([]string, 0, 1024)
	for _, item := range tmpForbidWordList {
		forbidWordList = append(forbidWordList, item.Word)
	}

	// 构造DFAUtil对象
	forbidDFAObj = dfaUtil.NewDFAUtil(forbidWordList)

	debugUtil.Printf("forbidWordList:%v\n", forbidWordList)

	return nil
}

// 是否包含屏蔽词
func IfContainsForbidWords(input string) bool {
	forbidMutex.RLock()
	defer forbidMutex.RUnlock()
	if len(forbidWordList) == 0 {
		return false
	}

	return forbidDFAObj.IsMatch(input)
}
