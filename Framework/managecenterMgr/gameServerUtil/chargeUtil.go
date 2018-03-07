package gameServerUtil

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"

	"work.goproject.com/Framework/managecenterMgr"
	. "work.goproject.com/Framework/managecenterModel/chargeConfig"
	. "work.goproject.com/Framework/managecenterModel/partner"
	"work.goproject.com/goutil/logUtil"
	"work.goproject.com/goutil/webUtil"
)

type ChargeUtil struct{}

// 获取有序的充值配置列表
// partnerId：合作商Id
// isMonthCard：是否月卡
// 返回值：
// 充值配置列表
// 是否存在
// 错误对象
func (this *ChargeUtil) getSortedChargeConfigList(partnerId int32, isMonthCard bool) (chargeConfigList []*ChargeConfig, exists bool, err error) {
	// 判断合作商是否存在
	var partnerObj *Partner
	partnerObj, exists = managecenterMgr.GetPartner(partnerId)
	if !exists {
		return
	}

	// 反序列化充值配置
	tmpChargeConfigList := make([]*ChargeConfig, 0, 16)
	if err = json.Unmarshal([]byte(partnerObj.ChargeConfig), &tmpChargeConfigList); err != nil {
		return
	}

	// 根据isMonthCard进行筛选
	for _, item := range tmpChargeConfigList {
		if item.IsMonthCard == isMonthCard {
			chargeConfigList = append(chargeConfigList, item)
		}
	}

	// 判断是否有符合条件的数据
	if len(chargeConfigList) == 0 {
		exists = false
		return
	}

	// 按默认规则进行排序
	sort.Slice(chargeConfigList, func(i, j int) bool {
		return chargeConfigList[i].SortByChargePointAsc(chargeConfigList[j])
	})

	return
}

// 获取充值配置列表
// partnerId：合作商Id
// vipLv：玩家VIP等级
// isMonthCard：是否月卡
// filter：过滤器
// 返回值：
// 充值配置列表
// 是否存在
// 错误对象
func (this *ChargeUtil) GetChargeConfigList(partnerId int32, vipLv byte, isMonthCard bool,
	filter func(*ChargeConfig) bool) (chargeConfigList []*ChargeConfig, exists bool, err error) {

	// 获取排好序的充值配置列表
	tmpList := make([]*ChargeConfig, 0, 8)
	tmpList, exists, err = this.getSortedChargeConfigList(partnerId, isMonthCard)
	if err != nil || !exists {
		return
	}

	// 根据vip和filter进行筛选
	for _, item := range tmpList {
		if filter != nil {
			if vipLv >= item.VipLv && filter(item) {
				chargeConfigList = append(chargeConfigList, item)
			}
		} else {
			if vipLv >= item.VipLv {
				chargeConfigList = append(chargeConfigList, item)
			}
		}
	}

	// 判断是否有符合条件的数据
	if len(chargeConfigList) == 0 {
		exists = false
		return
	}

	return
}

// 获取充值配置项
// partnerId：合作商Id
// vipLv：玩家VIP等级
// isMonthCard：是否月卡
// money：充值金额
// 返回值：
// 充值配置项
// 是否存在
// 错误对象
func (this *ChargeUtil) GetChargeConfigItem(partnerId int32, vipLv byte, isMonthCard bool, money float64) (chargeConfigItem *ChargeConfig, exists bool, err error) {
	// 获取排好序的充值配置列表
	tmpList := make([]*ChargeConfig, 0, 8)
	tmpList, exists, err = this.getSortedChargeConfigList(partnerId, isMonthCard)
	if err != nil || !exists {
		return
	}

	// 获取满足充值金额和VIP条件的最后一条数据
	for _, item := range tmpList {
		if vipLv >= item.VipLv && money >= item.ChargePoint {
			chargeConfigItem = item
		} else {
			break
		}
	}

	// 如果没有符合条件的，则选择第一条配置
	if chargeConfigItem == nil {
		chargeConfigItem = tmpList[0]
	}

	return
}

// 生成充值订单号
// url:生成订单号的服务器地址
// productId:产品Id
// partnerId:合作商Id
// serverId:服务器Id
// userId:平台用户Id
// playerId:玩家Id
// mac:mac
// idfa:idfa
// ip:ip
// imei:imei
// extra:extra
// isMonthCard:是否月卡
// 返回值:
// 订单号
// 错误对象
func (this *ChargeUtil) GenerateOrderId(url, productId string, partnerId, serverId int32,
	userId, playerId, mac, idfa, ip, imei, extra string,
	isMonthCard bool) (orderId string, err error) {

	if extra == "" {
		extra = "FromGameServer"
	}

	// 定义请求参数
	postDict := make(map[string]string)
	postDict["ProductId"] = productId
	postDict["PartnerId"] = fmt.Sprintf("%d", partnerId)
	postDict["ServerId"] = fmt.Sprintf("%d", serverId)
	postDict["UserId"] = userId
	postDict["PlayerId"] = playerId
	postDict["MAC"] = mac
	postDict["IDFA"] = idfa
	postDict["IP"] = ip
	postDict["IMEI"] = imei
	postDict["Extra"] = extra
	postDict["IsMonthCard"] = strconv.FormatBool(isMonthCard)

	// 连接充值服务器，以生成订单号
	var returnBytes []byte
	if returnBytes, err = webUtil.PostWebData(url, postDict, nil); err != nil {
		logUtil.ErrorLog(fmt.Sprintf("生成订单号错误:%s,错误信息为:%s", url, err))
		return
	} else {
		orderId = string(returnBytes)
	}

	if orderId == "" {
		err = fmt.Errorf("Order Is Empty")
	}

	return
}

// 计算充值获得的游戏点数
// partnerId：合作商Id
// vipLv：玩家VIP等级
// isMonthCard：是否月卡
// money：充值金额
// 返回值:
// 充值获得的游戏点数
// 是否存在
// 错误对象
func (this *ChargeUtil) CalcChargeGamePoint(partnerId int32, vipLv byte, isMonthCard bool, money float64) (chargeGamePoint int, exists bool, err error) {
	// 获取排好序的充值配置列表
	tmpList := make([]*ChargeConfig, 0, 8)
	tmpList, exists, err = this.getSortedChargeConfigList(partnerId, isMonthCard)
	if err != nil || !exists {
		return
	}

	var chargeConfigItem *ChargeConfig

	// 获取满足充值金额和VIP条件的最后一条数据
	for _, item := range tmpList {
		if vipLv >= item.VipLv && money >= item.ChargePoint {
			chargeConfigItem = item
		} else {
			break
		}
	}

	// 如果找不到对应的档位，则选择最低金额档位
	if chargeConfigItem == nil {
		chargeConfigItem = tmpList[0]
	}

	// 计算充值对应的ProductId，以及获得的游戏货币
	if money == chargeConfigItem.ChargePoint {
		chargeGamePoint = chargeConfigItem.GamePoint
	} else {
		chargeGamePoint = int(math.Ceil(money * chargeConfigItem.Ratio))
	}

	return
}

// 计算充值获得的所有游戏点数
// partnerId：合作商Id
// vipLv：玩家VIP等级
// isMonthCard：是否月卡
// money：充值金额
// activityMoney：活动金额
// isFirstCharge：是否首充
// 返回值:
// 充值获得的游戏点数
// 充值赠送获得的游戏内货币数量
// 充值活动获得的游戏内货币数量
// 总的元宝数量
// 是否存在
// 错误对象
func (this *ChargeUtil) CalcChargeAllGamePoint(partnerId int32, vipLv byte, isMonthCard bool,
	money, activityMoney float64, isFirstCharge bool) (
	chargeGamePoint int, giveGamePoint int, activityGamePoint int, totalGamePoint int,
	exists bool, err error) {

	// 获取排好序的充值配置列表
	tmpList := make([]*ChargeConfig, 0, 8)
	tmpList, exists, err = this.getSortedChargeConfigList(partnerId, isMonthCard)
	if err != nil || !exists {
		return
	}

	var chargeConfigItem *ChargeConfig

	// 获取满足充值金额和VIP条件的最后一条数据
	for _, item := range tmpList {
		if vipLv >= item.VipLv && money >= item.ChargePoint {
			chargeConfigItem = item
		} else {
			break
		}
	}

	// 如果找不到对应的档位，则选择最低金额档位
	if chargeConfigItem == nil {
		chargeConfigItem = tmpList[0]
	}

	// 计算充值对应的ProductId，以及获得的游戏货币
	if money == chargeConfigItem.ChargePoint {
		chargeGamePoint = chargeConfigItem.GamePoint
		if isFirstCharge {
			giveGamePoint = chargeConfigItem.FirstGiveGamePoint
		} else {
			giveGamePoint = chargeConfigItem.GiveGamePoint
		}
	} else {
		chargeGamePoint = int(math.Ceil(money * chargeConfigItem.Ratio))
		if isFirstCharge {
			giveGamePoint = chargeConfigItem.FirstGiveGamePoint
		} else {
			giveGamePoint = int(math.Ceil(money * chargeConfigItem.Ratio * chargeConfigItem.GiveRatio))
		}
	}

	activityGamePoint = int(math.Ceil(activityMoney * chargeConfigItem.Ratio))

	// 计算总和
	totalGamePoint = chargeGamePoint + giveGamePoint + activityGamePoint

	return
}

// ------------------类型定义和业务逻辑的分隔符-------------------------

var (
	ChargeUtilObj = new(ChargeUtil)
)
