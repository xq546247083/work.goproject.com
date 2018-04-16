package typeUtil

import "testing"

func TestSliceToString(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	if result, _ := SliceToString(data, ","); result != "1,2,3,4,5" {
		t.Error("转换错误:" + result)
		return
	}
}
func TestSliceToString2(t *testing.T) {
	data := []int{1}

	if result, _ := SliceToString(data, ","); result != "1,2,3,4,5" {
		t.Error("转换错误:" + result)
		return
	}
}

func TestSliceToString3(t *testing.T) {
	var data []int = nil
	if result, _ := SliceToString(data, ","); result != "" {
		t.Error("转换错误:" + result)
		return
	}
}
