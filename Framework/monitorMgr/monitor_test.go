package monitorMgr

import (
	"testing"
)

func TestDuplicate(t *testing.T) {
	SetParam("20.255.0.7", "Test", 5)

	content := "content"

	if isDuplicate(content) == true {
		t.Errorf("%s不应该重复却重复了", content)
	}
	addToMap(content)
	if isDuplicate(content) == false {
		t.Errorf("%s应该重复却没有重复", content)
	}
}
