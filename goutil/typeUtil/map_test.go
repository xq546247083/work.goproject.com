package typeUtil

import "testing"

func TestMapToString(t *testing.T) {
	data := make(map[int]int)
	data[0] = 0
	data[1] = 1

	strVal, err := MapToString(data, ":", ",")
	if err != nil {
		t.Error(err)
		return
	}

	t.Errorf(strVal)
}
