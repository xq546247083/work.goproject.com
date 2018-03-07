package mysqlSync

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"work.goproject.com/goutil/fileUtil"
	"work.goproject.com/goutil/logUtil"
)

var (
	// 记录错误sql命令的文件名
	con_Error_FileName = "errorFile.txt"
)

// 定义处理错误命令的文件对象
type errorFile struct {
	// 错误文件
	file *os.File

	// 文件路径
	filePath string

	// 同步数据对象的唯一标识，用于进行重复判断
	identifier string
}

// 保存命令到错误文件
// command: sql命令
func (this *errorFile) SaveCommand(command string) {
	this.open()
	defer this.close()

	// 覆盖写入
	this.file.Seek(0, 0)

	// 写入命令
	_, err := this.file.WriteString(command)
	if err != nil {
		prefix := fmt.Sprintf("%s-%s", this.identifier, "errorFile.SaveCommand")
		err = fmt.Errorf("%s-Write %s to file failed:%s", prefix, command, err)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}

	// 清理残留数据
	this.file.Truncate(int64(len(command)))
}

// 读取文件中命令
func (this *errorFile) ReadCommand() string {
	this.open()
	defer this.close()

	this.file.Seek(0, 0)
	content, err := ioutil.ReadAll(this.file)
	if err != nil {
		prefix := fmt.Sprintf("%s-%s", this.identifier, "errorFile.ReadCommand")
		err = fmt.Errorf("%s-Read command failed:%s", prefix, err)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}
	return string(content)
}

// 打开文件
func (this *errorFile) open() {
	// 打开errorFile文件, 如果没有就创建
	var err error
	this.file, err = os.OpenFile(this.filePath, os.O_CREATE|os.O_RDWR, os.ModePerm|os.ModeTemporary)
	if err != nil {
		prefix := fmt.Sprintf("%s-%s", this.identifier, "errorFile.newErrorFile.os.OpenFile")
		err = fmt.Errorf("%s-Open File failed:%s", prefix, err)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}
}

// 关闭文件
func (this *errorFile) close() {
	this.file.Close()
}

// 删除文件
func (this *errorFile) Delete() {
	fileUtil.DeleteFile(this.filePath)
}

// 构造错误文件对象
// _dirPath:文件路径
// _identifier:唯一标识
func newErrorFile(_dirPath string, _identifier string) *errorFile {
	_filePath := filepath.Join(_dirPath, con_Error_FileName)
	return &errorFile{
		filePath:   _filePath,
		identifier: _identifier,
	}
}
