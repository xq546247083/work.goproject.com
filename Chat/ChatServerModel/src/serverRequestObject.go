package src

// 请求对象
type ServerRequestObject struct {
	// 请求的唯一标识，是需要通过截取请求数据前4位得到并进行手动赋值的（暂时未使用）
	Id int32

	// 以下属性是由客户端直接传入的，可以直接反序列化直接得到的
	// 请求的模块名称
	ModuleName string

	// 请求的方法名称
	MethodName string

	// 请求的参数数组
	Parameters []interface{}
}
