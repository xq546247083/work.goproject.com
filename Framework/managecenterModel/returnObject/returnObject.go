package returnObject

type ReturnObject struct {
	// 返回的状态值；0：成功；非0：失败（根据实际情况进行定义）
	Code int32

	// 返回的失败描述信息
	Message string

	// 返回的数据
	Data interface{}
}
