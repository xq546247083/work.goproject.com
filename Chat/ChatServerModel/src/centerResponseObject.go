package src

// ChatServerCenter的响应对象
type CenterResponseObject struct {
	// 响应结果的状态值
	*ResultStatus

	// 响应结果的数据
	Data interface{}
}

func (this *CenterResponseObject) SetResultStatus(rs *ResultStatus) *CenterResponseObject {
	this.ResultStatus = rs

	return this
}

func (this *CenterResponseObject) SetData(data interface{}) *CenterResponseObject {
	this.Data = data

	return this
}

func NewCenterResponseObject() *CenterResponseObject {
	return &CenterResponseObject{
		ResultStatus: Success,
		Data:         nil,
	}
}
