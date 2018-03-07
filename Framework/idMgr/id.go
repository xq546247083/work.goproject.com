/*
1、生成的每一个Id都将包含一个固定前缀值（对应游戏里面的服务器组Id），这样就可以让不同的服务器组生成的Id不重复，便于合服
2、为了生成尽可能多的不重复数字，所以使用int64来表示一个数字，其中：
0 000000000000000 0000000000000000000000000000 00000000000000000000
1：固定为0
2-16：共15位，表示固定前缀值。范围为[0, math.Pow(2, 15))
17-44：共28位，表示当前时间距离基础时间的秒数。范围为[0, math.Pow(2, 28))，约合8.5年，以2017-1-1 00:00:00为基准则可以持续到2025-07-01 00:00:00
45-64：共20位，表示自增种子。范围为[0, math.Pow(2, 20))，共1048576个数字
3、总体而言，此规则支持每秒生成1048576个不同的数字，并且在8.5年的时间范围内有效
*/

package idMgr

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	// 使用15位来表示固定前缀值，范围为[0, math.Pow(2, 15)=32768)
	con_Prefix_BitCount = 15

	// 使用28位来时间的秒数，可表示的范围为[0, math.Pow(2, 28))，约合8.5年；
	con_Time_BitCount = 28

	// 使用20位来表示自增种子，可表示的范围为[0, math.Pow(2, 20)=1048576)
	con_Seed_BitCount = 20
)

var (
	prefix_Min int64 = 0
	prefix_Max int64 = int64(math.Pow(2, con_Prefix_BitCount)) - 1
	baseTime         = time.Date(2017, time.January, 1, 0, 0, 0, 0, time.Local)
	seed       int64 = 0
	seed_Min   int64 = 0
	seed_Max   int64 = int64(math.Pow(2, con_Seed_BitCount)) - 1
	mutex      sync.Mutex
)

func generateSeed() int64 {
	mutex.Lock()
	defer mutex.Unlock()

	if seed >= seed_Max {
		seed = seed_Min
	} else {
		seed = seed + 1
	}

	return seed
}

// 生成新的Id
// prefix：固定前缀值，范围[0,32768)
// 返回值：
// 新的Id
// 错误对象
func GenerateNewId(prefix int64) (int64, error) {
	if prefix < prefix_Min || prefix > prefix_Max {
		return 0, fmt.Errorf("prefix溢出，有效范围为[%d,%d)", prefix_Min, prefix_Max)
	}

	tick := int64(time.Now().Sub(baseTime) / time.Second)
	seed := generateSeed()
	id := prefix<<48 + tick<<20 + seed

	return id, nil
}
