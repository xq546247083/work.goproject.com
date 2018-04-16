package typeUtil

import (
	"fmt"
	"reflect"
)

// 把一个集合转换成字符串
// data:slice类型的集合
// seprator:分隔符
// 返回值:
// result:转换后的字符串
// err:错误信息对象
func SliceToString(data interface{}, seprator string) (result string, err error) {
	if data == nil {
		return
	}

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		err = fmt.Errorf("目标类型不正确，只能是slice或array 当前类型是:%v", val.Kind().String())
		return
	}

	if val.Len() <= 0 {
		return
	}

	for i := 0; i < val.Len(); i++ {
		valItem := val.Index(i)
		result = result + fmt.Sprintf("%v", valItem.Interface()) + seprator
	}

	if val.Len() > 0 {
		result = result[:len(result)-1]
	}

	return
}
