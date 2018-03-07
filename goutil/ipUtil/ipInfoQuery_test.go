package ipUtil

import (
	"testing"
)

func TestQueryIpInfo(t *testing.T) {
	info, err := QueryIpInfo2("165.161.111.1", "xh1", "1", 5)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Errorf("当前获取到的数据为:%v", info)
}
