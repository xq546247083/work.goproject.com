package verifyMgr

import (
	"fmt"
	"testing"

	"work.goproject.com/Framework/managecenterMgr"
)

func init() {
	var url string = "http://managecenter.dzz.work.goproject.com/API"
	var groupType string = "Mix"
	var ip string = ""
	var groupId int32 = 0
	var ifOpen bool = false
	var interval int = 5
	managecenterMgr.Start(url, groupType, ip, groupId, ifOpen, interval)
}

func TestVerify(t *testing.T) {
	if errList := Verify(); len(errList) > 0 {
		t.Errorf("there should be no error, but now has")
		for _, err := range errList {
			fmt.Println(err)
		}
	}

	Init("http://www.baidu.com", Con_Get)
	if errList := Verify(); len(errList) > 0 {
		t.Errorf("there should be no error, but now has")
		for _, err := range errList {
			fmt.Println(err)
		}
	}
}
