package goroutineMgr

import (
	"fmt"
	"runtime"
	"sort"
	"sync"

	"work.goproject.com/goutil/logUtil"
)

var (
	goroutineCountMap   = make(map[string]int)
	goroutineCountMutex sync.RWMutex
)

// 增加指定名称的goroutine的数量
// goroutineName:goroutine名称
func increaseCount(goroutineName string) {
	goroutineCountMutex.Lock()
	defer goroutineCountMutex.Unlock()

	newCount := 1
	if currCount, exists := goroutineCountMap[goroutineName]; exists {
		newCount = currCount + 1
	}

	goroutineCountMap[goroutineName] = newCount
}

// 减少指定名称的goroutine的数量
// goroutineName:goroutine名称
func decreaseCount(goroutineName string) {
	goroutineCountMutex.Lock()
	defer goroutineCountMutex.Unlock()

	newCount := -1
	if currCount, exists := goroutineCountMap[goroutineName]; exists {
		newCount = currCount - 1
	}

	if newCount <= 0 {
		delete(goroutineCountMap, goroutineName)
	} else {
		goroutineCountMap[goroutineName] = newCount
	}
}

// 获取指定名称的goroutine的数量
// goroutineName:goroutine名称
// 返回值:
// 对应数量
func getGoroutineCount(goroutineName string) int {
	goroutineCountMutex.RLock()
	defer goroutineCountMutex.RUnlock()

	if currCount, exists := goroutineCountMap[goroutineName]; exists {
		return currCount
	} else {
		return 0
	}
}

// 转化成字符串
func toString() string {
	goroutineCountMutex.RLock()
	defer goroutineCountMutex.RUnlock()

	keys := make([]string, 0, 16)
	for key, _ := range goroutineCountMap {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	str := fmt.Sprintf("Goroutine Info:(%s,%d)", "NumGoroutine", runtime.NumGoroutine())
	for _, key := range keys {
		str += fmt.Sprintf("(%s,%d)", key, goroutineCountMap[key])
	}

	return str
}

// 记录goroutine数量信息
func logGoroutineCountInfo() {
	logUtil.NormalLog(toString(), logUtil.Debug)
}

func Test() {
	logGoroutineCountInfo()
}
