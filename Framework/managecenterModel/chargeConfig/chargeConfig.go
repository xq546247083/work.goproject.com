package chargeConfig

// 充值配置
type ChargeConfig struct {
	// 产品Id
	ProductId string

	// 充值点数（以元为单位；如果是人民币是整数，但如果是美元，或者其它货币可能为小数）
	ChargePoint float64

	// 游戏内货币点数（元宝/钻石等，必定是整数）
	GamePoint int

	// 充值金额与游戏内货币的兑换比率
	Ratio float64

	// 赠送的游戏内货币点数（必定是整数）
	GiveGamePoint int

	// 赠送的比率
	GiveRatio float64

	// 首充时赠送的游戏内货币点数
	FirstGiveGamePoint int

	// 所需的vip等级
	VipLv byte

	// 首充时是否显示
	IfFirstShow byte

	// 第二次（及以后）充值时是否显示
	IfSecondShow byte

	// 是否为月卡
	IsMonthCard bool
}

// 按照充值金额进行升序排序
// target:另一个充值配置对象
// 是否是小于
func (this *ChargeConfig) SortByChargePointAsc(target *ChargeConfig) bool {
	return this.ChargePoint < target.ChargePoint
}
