package gameServerUtil

import (
	"encoding/json"
	"fmt"
	"testing"

	"work.goproject.com/Framework/managecenterMgr"
	. "work.goproject.com/Framework/managecenterModel/chargeConfig"
	"work.goproject.com/goutil/logUtil"
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

func TestGetChargeConfigList(t *testing.T) {
	var partnerId int32 = 1
	var vipLv byte = 1
	var isMonthCard bool = false
	var expectedCount = 1
	filter := func(item *ChargeConfig) bool {
		if item.IfFirstShow == 1 {
			return true
		} else {
			return false
		}
	}

	chargeConfigList, exists, err := ChargeUtilObj.GetChargeConfigList(partnerId, vipLv, isMonthCard, filter)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("1:应该存在充值配置，但是现在不存在")
	} else {
		bytes, err := json.Marshal(chargeConfigList)
		if err != nil {
			panic(err)
		}
		logUtil.InfoLog("ChargeConfigList:%s", string(bytes))
	}
	if len(chargeConfigList) != expectedCount {
		t.Error("1:ChargeConfigCount should be %d, but now %d.", expectedCount, len(chargeConfigList))
	}

	isMonthCard = false
	vipLv = 5
	expectedCount = 2
	chargeConfigList, exists, err = ChargeUtilObj.GetChargeConfigList(partnerId, vipLv, isMonthCard, filter)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("2:应该存在充值配置，但是现在不存在")
	} else {
		bytes, err := json.Marshal(chargeConfigList)
		if err != nil {
			panic(err)
		}
		logUtil.InfoLog("ChargeConfigList:%s", string(bytes))
	}
	if len(chargeConfigList) != expectedCount {
		t.Errorf("2:ChargeConfigCount should be %d, but now %d.", expectedCount, len(chargeConfigList))
	}

	isMonthCard = true
	vipLv = 5
	expectedCount = 0
	chargeConfigList, exists, err = ChargeUtilObj.GetChargeConfigList(partnerId, vipLv, isMonthCard, nil)
	if err != nil {
		panic(err)
	} else if exists {
		t.Error("3:不应该存在充值配置，但是现在却存在")
	} else {
		bytes, err := json.Marshal(chargeConfigList)
		if err != nil {
			panic(err)
		}
		logUtil.InfoLog("ChargeConfigList:%s", string(bytes))
	}
	if len(chargeConfigList) != expectedCount {
		t.Errorf("3:ChargeConfigCount should be %d, but now %d.", expectedCount, len(chargeConfigList))
	}

	isMonthCard = true
	vipLv = 10
	expectedCount = 1
	chargeConfigList, exists, err = ChargeUtilObj.GetChargeConfigList(partnerId, vipLv, isMonthCard, nil)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("4:应该存在充值配置，但是现在不存在")
	} else {
		bytes, err := json.Marshal(chargeConfigList)
		if err != nil {
			panic(err)
		}
		logUtil.InfoLog("ChargeConfigList:%s", string(bytes))
	}
	if len(chargeConfigList) != expectedCount {
		t.Errorf("4:ChargeConfigCount should be %d, but now %d.", expectedCount, len(chargeConfigList))
	}
}

func TestGetChargeConfigItem(t *testing.T) {
	var partnerId int32 = 1
	var vipLv byte = 1
	var isMonthCard bool = false
	var money float64 = 6.0
	var expectedChargePoint float64 = 6.0

	item, exists, err := ChargeUtilObj.GetChargeConfigItem(partnerId, vipLv, isMonthCard, money)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("1:应该存在充值配置，但是现在不存在")
	}
	if item.ChargePoint != expectedChargePoint {
		t.Errorf("1:期望得到%f, But now %f", expectedChargePoint, item.ChargePoint)
	}

	money = 5.0

	item, exists, err = ChargeUtilObj.GetChargeConfigItem(partnerId, vipLv, isMonthCard, money)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("2:应该存在充值配置，但是现在不存在")
	}
	if item.ChargePoint != expectedChargePoint {
		t.Errorf("2:期望得到%f, But now %f", expectedChargePoint, item.ChargePoint)
	}

	money = 30.0

	item, exists, err = ChargeUtilObj.GetChargeConfigItem(partnerId, vipLv, isMonthCard, money)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("2:应该存在充值配置，但是现在不存在")
	}
	if item.ChargePoint != expectedChargePoint {
		t.Errorf("2:期望得到%f, But now %f", expectedChargePoint, item.ChargePoint)
	}

	isMonthCard = true
	money = 50.0
	expectedChargePoint = 50.0

	item, exists, err = ChargeUtilObj.GetChargeConfigItem(partnerId, vipLv, isMonthCard, money)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("2:应该存在充值配置，但是现在不存在")
	}
	if item.ChargePoint != expectedChargePoint {
		t.Errorf("2:期望得到%f, But now %f", expectedChargePoint, item.ChargePoint)
	}
}

func TestGenerateOrderId(t *testing.T) {
	var url string = "http://chargetest.hzgg.work.goproject.com/API/GenerateOrderId.ashx"
	var productId string = "xdzz_6"
	var partnerId int32 = 1
	var serverId int32 = 20008
	var userId string = "aa2ab1e9af3041cbaf4e9e25e71abc98"
	var playerId string = "01e31204-d3fd-4262-a4f1-e38fe0d93b6b"
	var mac string = ""
	var idfa string = "1AD4F263-5C77-4298-B826-9E9A02C01F62"
	var ip string = "117.139.247.210"
	var imei string = "e592d423-373f-47f7-a3cf-de3a4f3ee883"
	var extra string = "米大师自动补单,累充金额:420,已处理金额:360.00,补单金额:60.00"
	var isMonthCard bool = true
	orderId, err := ChargeUtilObj.GenerateOrderId(url, productId, partnerId, serverId, userId, playerId, mac, idfa, ip, imei, extra, isMonthCard)
	if err != nil {
		panic(err)
	}
	fmt.Printf("OrderId:%s\n", orderId)
}

func TestCalcChargeGamePoint(t *testing.T) {
	var partnerId int32 = 1
	var vipLv byte = 1
	var isMonthCard bool = false
	var money float64 = 6.0
	var expectedChargeGamePoint int = 60

	chargeGamePoint, exists, err := ChargeUtilObj.CalcChargeGamePoint(partnerId, vipLv, isMonthCard, money)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("1:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint {
		t.Errorf("1:期望得到%d, But now %d", expectedChargeGamePoint, chargeGamePoint)
	}

	money = 25.0
	expectedChargeGamePoint = 250
	chargeGamePoint, exists, err = ChargeUtilObj.CalcChargeGamePoint(partnerId, vipLv, isMonthCard, money)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("2:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint {
		t.Errorf("2:期望得到%d, But now %d", expectedChargeGamePoint, chargeGamePoint)
	}

	money = 50
	isMonthCard = true
	expectedChargeGamePoint = 500
	chargeGamePoint, exists, err = ChargeUtilObj.CalcChargeGamePoint(partnerId, vipLv, isMonthCard, money)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("3:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint {
		t.Errorf("3:期望得到%d, But now %d", expectedChargeGamePoint, chargeGamePoint)
	}
}

func TestCalcChargeAllGamePoint(t *testing.T) {
	var partnerId int32 = 1
	var vipLv byte = 1
	var isMonthCard bool = false
	var money float64 = 6.0
	var activityMoney float64 = 0.0
	var isFirstCharge bool = false
	var expectedChargeGamePoint int = 60
	var expectedGiveGamePoint int = 6
	var expectedActivityGamePoint int = 0
	var expectedTotalGamePoint int = 66

	chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint, exists, err := ChargeUtilObj.CalcChargeAllGamePoint(partnerId, vipLv, isMonthCard, money, activityMoney, isFirstCharge)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("1:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint || giveGamePoint != expectedGiveGamePoint ||
		activityGamePoint != expectedActivityGamePoint || totalGamePoint != expectedTotalGamePoint {
		t.Errorf("1:期望得到%d,%d,%d,%d, But now %d,%d,%d,%d",
			expectedChargeGamePoint, expectedGiveGamePoint, expectedActivityGamePoint, expectedTotalGamePoint,
			chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint)
	}

	isFirstCharge = true
	expectedChargeGamePoint = 60
	expectedGiveGamePoint = 60
	expectedActivityGamePoint = 0
	expectedTotalGamePoint = 120

	chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint, exists, err = ChargeUtilObj.CalcChargeAllGamePoint(partnerId, vipLv, isMonthCard, money, activityMoney, isFirstCharge)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("2:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint || giveGamePoint != expectedGiveGamePoint ||
		activityGamePoint != expectedActivityGamePoint || totalGamePoint != expectedTotalGamePoint {
		t.Errorf("2:期望得到%d,%d,%d,%d, But now %d,%d,%d,%d",
			expectedChargeGamePoint, expectedGiveGamePoint, expectedActivityGamePoint, expectedTotalGamePoint,
			chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint)
	}

	isFirstCharge = false
	activityMoney = 0.5
	expectedChargeGamePoint = 60
	expectedGiveGamePoint = 6
	expectedActivityGamePoint = 5
	expectedTotalGamePoint = 71

	chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint, exists, err = ChargeUtilObj.CalcChargeAllGamePoint(partnerId, vipLv, isMonthCard, money, activityMoney, isFirstCharge)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("3:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint || giveGamePoint != expectedGiveGamePoint ||
		activityGamePoint != expectedActivityGamePoint || totalGamePoint != expectedTotalGamePoint {
		t.Errorf("3:期望得到%d,%d,%d,%d, But now %d,%d,%d,%d",
			expectedChargeGamePoint, expectedGiveGamePoint, expectedActivityGamePoint, expectedTotalGamePoint,
			chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint)
	}

	money = 10.0
	activityMoney = 1.0
	isFirstCharge = false
	expectedChargeGamePoint = 100
	expectedGiveGamePoint = 10
	expectedActivityGamePoint = 10
	expectedTotalGamePoint = 120

	chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint, exists, err = ChargeUtilObj.CalcChargeAllGamePoint(partnerId, vipLv, isMonthCard, money, activityMoney, isFirstCharge)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("4:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint || giveGamePoint != expectedGiveGamePoint ||
		activityGamePoint != expectedActivityGamePoint || totalGamePoint != expectedTotalGamePoint {
		t.Errorf("4:期望得到%d,%d,%d,%d, But now %d,%d,%d,%d",
			expectedChargeGamePoint, expectedGiveGamePoint, expectedActivityGamePoint, expectedTotalGamePoint,
			chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint)
	}

	money = 10.0
	activityMoney = 1.0
	isFirstCharge = true
	expectedChargeGamePoint = 100
	expectedGiveGamePoint = 60
	expectedActivityGamePoint = 10
	expectedTotalGamePoint = 170

	chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint, exists, err = ChargeUtilObj.CalcChargeAllGamePoint(partnerId, vipLv, isMonthCard, money, activityMoney, isFirstCharge)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("5:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint || giveGamePoint != expectedGiveGamePoint ||
		activityGamePoint != expectedActivityGamePoint || totalGamePoint != expectedTotalGamePoint {
		t.Errorf("5:期望得到%d,%d,%d,%d, But now %d,%d,%d,%d",
			expectedChargeGamePoint, expectedGiveGamePoint, expectedActivityGamePoint, expectedTotalGamePoint,
			chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint)
	}

	money = 50.0
	activityMoney = 0
	isFirstCharge = false
	isMonthCard = true
	vipLv = 10
	expectedChargeGamePoint = 500
	expectedGiveGamePoint = 55
	expectedActivityGamePoint = 0
	expectedTotalGamePoint = 555

	chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint, exists, err = ChargeUtilObj.CalcChargeAllGamePoint(partnerId, vipLv, isMonthCard, money, activityMoney, isFirstCharge)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("6:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint || giveGamePoint != expectedGiveGamePoint ||
		activityGamePoint != expectedActivityGamePoint || totalGamePoint != expectedTotalGamePoint {
		t.Errorf("6:期望得到%d,%d,%d,%d, But now %d,%d,%d,%d",
			expectedChargeGamePoint, expectedGiveGamePoint, expectedActivityGamePoint, expectedTotalGamePoint,
			chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint)
	}

	money = 50.0
	activityMoney = 0
	isFirstCharge = false
	isMonthCard = true
	vipLv = 1
	expectedChargeGamePoint = 500
	expectedGiveGamePoint = 55
	expectedActivityGamePoint = 0
	expectedTotalGamePoint = 555

	chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint, exists, err = ChargeUtilObj.CalcChargeAllGamePoint(partnerId, vipLv, isMonthCard, money, activityMoney, isFirstCharge)
	if err != nil {
		panic(err)
	} else if !exists {
		t.Error("7:应该存在充值配置，但是现在不存在")
	}
	if chargeGamePoint != expectedChargeGamePoint || giveGamePoint != expectedGiveGamePoint ||
		activityGamePoint != expectedActivityGamePoint || totalGamePoint != expectedTotalGamePoint {
		t.Errorf("7:期望得到%d,%d,%d,%d, But now %d,%d,%d,%d",
			expectedChargeGamePoint, expectedGiveGamePoint, expectedActivityGamePoint, expectedTotalGamePoint,
			chargeGamePoint, giveGamePoint, activityGamePoint, totalGamePoint)
	}
}
