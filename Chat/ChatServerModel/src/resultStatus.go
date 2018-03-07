package src

// 返回结果状态
type ResultStatus struct {
	// 状态值(成功是0，非成功以负数来表示)
	Code int

	// 状态信息
	Message string
}

func newResultStatus(code int, message string) *ResultStatus {
	return &ResultStatus{
		Code:    code,
		Message: message,
	}
}

// 定义所有的响应结果的状态枚举值
var (
	Success               = newResultStatus(0, "成功")
	DataError             = newResultStatus(-1, "数据错误")
	DBError               = newResultStatus(-2, "数据库错误")
	MethodNotDefined      = newResultStatus(-3, "方法未定义")
	ParamIsEmpty          = newResultStatus(-4, "参数为空")
	ParamNotMatch         = newResultStatus(-5, "参数不匹配")
	ParamTypeError        = newResultStatus(-6, "参数类型错误")
	OnlySupportPOST       = newResultStatus(-7, "只支持POST")
	APINotDefined         = newResultStatus(-8, "API未定义")
	APIParamError         = newResultStatus(-9, "API参数错误")
	InvalidIP             = newResultStatus(-10, "IP无效")
	PlayerNotExists       = newResultStatus(-11, "玩家不存在")
	NoAvailableServer     = newResultStatus(-12, "没有可用的服务器")
	ClientDataError       = newResultStatus(-13, "客户端数据错误")
	TokenInvalid          = newResultStatus(-14, "令牌无效")
	ChannelNotDefined     = newResultStatus(-15, "聊天频道未定义")
	NoTargetMethod        = newResultStatus(-16, "找不到目标方法")
	ParamInValid          = newResultStatus(-17, "参数无效")
	NoLogin               = newResultStatus(-18, "尚未登陆")
	NotInUnion            = newResultStatus(-19, "不在公会中")
	NotInShimen           = newResultStatus(-20, "不在师门中")
	NotFoundTarget        = newResultStatus(-21, "未找到目标玩家")
	PlayerNotExist        = newResultStatus(-22, "玩家不存在")
	ServerGroupNotExist   = newResultStatus(-23, "服务器组不存在")
	NotInTeam             = newResultStatus(-24, "不在队伍中")
	LoginOnAnotherDevice  = newResultStatus(-25, "在另一台设备上登录")
	CantSendMessageToSelf = newResultStatus(-26, "不能给自己发消息")
	ResourceNotEnough     = newResultStatus(-27, "资源不足")
	NetworkError          = newResultStatus(-28, "网络错误")
	ContainForbiddenWord  = newResultStatus(-29, "含有屏蔽词语")
	SendMessageTooFast    = newResultStatus(-30, "发送消息太快")
	LvIsNotEnough         = newResultStatus(-31, "等级不足，系统未开放")
	RepeatTooMuch         = newResultStatus(-32, "重复次数太多")
	CantCrossServerTalk   = newResultStatus(-33, "不能跨服私聊")
	InSilent              = newResultStatus(-33, "已被禁言")
	NotInCountry          = newResultStatus(-34, "不在国家中")
)
