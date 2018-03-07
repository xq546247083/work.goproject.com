package sqlSync

import (
	"os"

	"encoding/binary"

	"errors"

	"work.goproject.com/goutil/fileUtil"
	"work.goproject.com/goutil/intAndBytesUtil"
)

const (
	// 头部字节长度
	con_Header_Length = 4
)

var (
	// 字节的大小端顺序
	byteOrder = binary.LittleEndian
)

// 按照指定方式读取文本内容
// fileObj:大文件对象
// data:待写入的数据
// 返回值:
// error:写入是否存在异常
func Write(fileObj *fileUtil.BigFile, data string) error {
	// 获得数据内容的长度
	dataLength := len(data)

	// 将长度转化为字节数组
	header := intAndBytesUtil.Int32ToBytes(int32(dataLength), byteOrder)

	// 将头部与内容组合在一起
	message := append(header, data...)

	// 写入数据
	return fileObj.WriteMessage(message)
}

// 从文件读取一条数据
// fileObj:文件对象
// 返回值:
// result:读取到的字符串
// err:错误信息
func Read(fileObj *os.File) (result string, readLen int64, err error) {
	// 1. 读取头部内容
	header := make([]byte, 4)
	var n int
	n, err = fileObj.Read(header)
	if err != nil {
		return
	}

	if n < con_Header_Length {
		err = errors.New("can not read 4 byte for read len")
		readLen = int64(n)
		return
	}

	dataLength := intAndBytesUtil.BytesToInt32(header, byteOrder)

	// 2. 读取指定长度的内容
	data := make([]byte, dataLength)
	n, err = fileObj.Read(data)
	if err != nil {
		return
	}

	readLen = int64(len(header) + int(dataLength))
	result = string(data)
	return
}
