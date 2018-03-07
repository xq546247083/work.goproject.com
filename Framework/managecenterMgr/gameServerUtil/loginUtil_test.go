package gameServerUtil

import (
	"testing"
)

func TestCheckLoginInfo(t *testing.T) {
	url := "http://loginsvrtest.hzgg.work.goproject.com/API/CheckDynamicLoginKey.ashx"
	partnerId := int32(1)
	userId := "aa2ab1e9af3041cbaf4e9e25e71abc98"
	loginInfo := "aa2ab1e9af3041cbaf4e9e25e71abc98"

	success, err := LoginUtilObj.CheckLoginInfo(url, partnerId, userId, loginInfo)
	if err == nil {
		t.Errorf("err should not be nil, but now is nil")
	}

	loginInfo = "aa2ab1e9af3041cbaf4e9e25e71abc98_1314579"
	success, err = LoginUtilObj.CheckLoginInfo(url, partnerId, userId, loginInfo)
	if err != nil {
		t.Errorf("err should not be nil, but now is nil")
	}
	if success == true {
		t.Errorf("success should be false, but now is true")
	}
}
