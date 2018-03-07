package gameServerUtil

import (
	"time"
	// "encoding/json"
	"fmt"
	"testing"

	"work.goproject.com/Framework/managecenterMgr"
)

func init() {
	var url string = "http://localhost:27135/API"
	var groupType string = "Mix"
	var ip string = "120.92.9.243"
	var groupId int32 = 9501
	var ifOpen bool = false
	var interval int = 0
	managecenterMgr.Start(url, groupType, ip, groupId, ifOpen, interval)
}

func TestCheckNewGameVersion(t *testing.T) {
	var partnerId int32 = 1
	var serverId int32 = 9501
	var gameVersionId int32 = 100

	gameVersionUrl, exists := ManageCenterUtilObj.CheckNewGameVersion(partnerId, serverId, gameVersionId)
	if exists {
		t.Errorf("there should be no new gameversion")
	}
	fmt.Printf("gameVersionUrl:%s\n", gameVersionUrl)

	gameVersionId = 98
	gameVersionUrl, exists = ManageCenterUtilObj.CheckNewGameVersion(partnerId, serverId, gameVersionId)
	if !exists {
		t.Errorf("there should be new gameversion")
	}
	fmt.Printf("gameVersionUrl:%s\n", gameVersionUrl)
}

func TestCheckNewResourceVersion(t *testing.T) {
	var partnerId int32 = 1
	var serverId int32 = 9501
	var gameVersionId int32 = 997
	var resourceVersionName string = "1472626972_790918539caccb9aac82f16584bb8284"

	name, url, exists := ManageCenterUtilObj.CheckNewResourceVersion(partnerId, serverId, gameVersionId, resourceVersionName)
	if !exists {
		t.Errorf("there should be new resource")
	}

	fmt.Printf("name:%s\n", name)
	fmt.Printf("url:%s\n", url)

	resourceVersionName = "1472632478_a9d19a19dfda7fe0d38afd2edc2db61b"
	name, url, exists = ManageCenterUtilObj.CheckNewResourceVersion(partnerId, serverId, gameVersionId, resourceVersionName)
	if exists {
		t.Errorf("there should be no new resource")
	}
}

func TestCheckMaintainStatus(t *testing.T) {
	var serverGroupId int32 = 9501

	maintainMessage, isMaintaining := ManageCenterUtilObj.CheckMaintainStatus(serverGroupId)
	if isMaintaining {
		t.Errorf("the servergroup should not in maintaining")
	}

	serverGroupObj, exists := managecenterMgr.GetServerGroupItem(serverGroupId)
	if !exists {
		t.Errorf("there should be servergroup")
	}

	serverGroupObj.GroupState = 2
	maintainMessage, isMaintaining = ManageCenterUtilObj.CheckMaintainStatus(serverGroupId)
	if !isMaintaining {
		t.Errorf("the servergroup should be in maintaining")
	}

	fmt.Printf("maintainMessage:%s\n", maintainMessage)

	serverGroupObj.GroupState = 1
	serverGroupObj.MaintainBeginTimeTick = time.Now().Unix()
	serverGroupObj.MaintainMinutes = 10

	maintainMessage, isMaintaining = ManageCenterUtilObj.CheckMaintainStatus(serverGroupId)
	if !isMaintaining {
		t.Errorf("the servergroup should be in maintaining")
	}

	fmt.Printf("maintainMessage:%s\n", maintainMessage)
}

func TestGetServerOpenDateInfo(t *testing.T) {
	var serverGroupId int32 = 9501

	serverGroupObj, exists := managecenterMgr.GetServerGroupItem(serverGroupId)
	if !exists {
		t.Errorf("there should be servergroup")
	}

	serverGroupObj.OpenTimeTick = time.Now().Unix()
	openDate, openDays := ManageCenterUtilObj.GetServerOpenDateInfo(serverGroupId)
	if openDays != 1 {
		t.Errorf("it should open for 1 day.")
	}
	fmt.Printf("openDate:%v\n", openDate)
	fmt.Printf("openDays:%d\n", openDays)
}

func TestGetPartnerServerPairString(t *testing.T) {
	var serverGroupId int32 = 9501

	pairString := ManageCenterUtilObj.GetPartnerServerPairString(serverGroupId)
	fmt.Printf("PairString:%s\n", pairString)
}

func TestIfValidPartnerServerPair(t *testing.T) {
	var partnerId int32 = 19
	var serverId int32 = 9501
	var parnterServerPairString string = "19_9501|46_9501|1004_9501|48_9501|319_9501|18_9501|3_9501|44_9501|11_9501|15_9501|313_9501|41_9501|43_9501|31_9501|9_9501|20_9501|38_9501|25_9501|2_9501|17_9501|504_9501|12_9501|27_9501|5_9501|42_9501|39_9501|1_9501|29_9501|40_9501|8_9501|23_9501|34_9501|1001_9501|37_9501|299_9501|26_9501|502_9501|1003_9501|10_9501|1000_9501|7_9501|35_9501|45_9501|16_9501|501_9501|14_9501|36_9501|49_9501|503_9501|13_9501|32_9501|383_9501|500_9501|999_9501|1002_9501|6_9501|24_9501|21_9501|22_9501|28_9501|33_9501|300_9501|47_9501|312_9501|4_9501|30_9501|"

	valid := ManageCenterUtilObj.IfValidPartnerServerPair(partnerId, serverId, parnterServerPairString)
	if !valid {
		t.Errorf("it should be valid")
	}

	parnterServerPairString = "19_9502|"
	valid = ManageCenterUtilObj.IfValidPartnerServerPair(partnerId, serverId, parnterServerPairString)
	if valid {
		t.Errorf("it should not be valid")
	}
}
