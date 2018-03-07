package src

// ChatServer的响应对象
type ServerResponseObject struct {
	// 响应结果的状态值
	*ResultStatus

	// 响应结果的数据
	Data interface{} `json:"Data,omitempty"`

	// 响应结果对应的请求的方法名称
	MethodName string
}

func (this *ServerResponseObject) SetResultStatus(rs *ResultStatus) *ServerResponseObject {
	this.ResultStatus = rs

	return this
}

func (this *ServerResponseObject) SetData(data interface{}) *ServerResponseObject {
	this.Data = data

	return this
}

func (this *ServerResponseObject) SetMethodName(methodName string) *ServerResponseObject {
	this.MethodName = methodName

	return this
}

func NewServerResponseObject() *ServerResponseObject {
	return &ServerResponseObject{
		ResultStatus: Success,
		Data:         nil,
		MethodName:   "",
	}
}
