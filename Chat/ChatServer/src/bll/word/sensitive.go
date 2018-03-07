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
	sensitiveWordList []string = make([]string, 0, 1024)
	sensitiveDFAObj   *dfaUtil.DFAUtil
	sensitiveMutex    sync.RWMutex
)

func init() {
	if err := reloadSensitive(); err != nil {
		panic(fmt.Errorf("初始化敏感词列表失败，错误信息为：%s", err))
	}

	// 注册重新加载的方法
	reloadMgr.RegisterReloadFunc("word.reloadSensitive", reloadSensitive)
}

// 重新加载敏感词列表
func reloadSensitive() error {
	var tmpSensitiveWordList []*SensitiveWord
	if err := dal.GetAll(&tmpSensitiveWordList); err != nil {
		return err
	}

	sensitiveMutex.Lock()
	defer sensitiveMutex.Unlock()
	sensitiveWordList = make([]string, 0, 1024)
	for _, item := range tmpSensitiveWordList {
		sensitiveWordList = append(sensitiveWordList, item.Word)
	}

	// 构造DFAUtil对象
	sensitiveDFAObj = dfaUtil.NewDFAUtil(sensitiveWordList)

	debugUtil.Printf("sensitiveWordList count:%d,first:%s\n", len(sensitiveWordList), sensitiveWordList[0])

	return nil
}

// 处理屏蔽词汇
// 输入字符串
// 处理屏蔽词汇后的字符串
func HandleSensitiveWords(input string) string {
	sensitiveMutex.RLock()
	defer sensitiveMutex.RUnlock()
	if len(sensitiveWordList) == 0 {
		return input
	}

	return sensitiveDFAObj.HandleWord(input, '*')
}
