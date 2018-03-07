package ipMgr

import (
	"testing"
)

func TestIsIpValid(t *testing.T) {
	ip := "10.255.0.7"
	if IsIpValid(ip) {
		t.Errorf("%s应该无效，但是现在却有效", ip)
	}

	ipList := []string{"10.255.0.7", "10.1.0.21"}
	Init(ipList)

	if IsIpValid(ip) == false {
		t.Errorf("%s应该有效，但是现在却无效", ip)
	}

	ipStr := "10.255.0.7,10.1.0.21;10.255.0.6|10.1.0.30||"
	InitString(ipStr)
	if IsIpValid(ip) == false {
		t.Errorf("%s应该有效，但是现在却无效", ip)
	}
}
