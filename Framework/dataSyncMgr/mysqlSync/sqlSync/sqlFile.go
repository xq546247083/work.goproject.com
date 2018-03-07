package sqlSync

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"sort"

	"work.goproject.com/goutil/fileUtil"
	"work.goproject.com/goutil/logUtil"
)

const (
	// 第一个文件名
	con_Default_FileName = "00000000"

	// 文件名后缀
	con_FileName_Suffix = "data"
)

// 同步数据对象(用于往文件中写入sql语句)
type SqlFile struct {
	// 存放同步数据的文件夹路径
	dirPath string

	// 同步数据对象的唯一标识，用于进行重复判断
	identifier string

	// 保存数据的大文件对象
	bigFileObj *fileUtil.BigFile

	// 数据同步对象
	mutex sync.Mutex
}

// 将数据写入同步数据对象
// data:待写入的数据
func (this *SqlFile) Write(data string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 写入数据
	err := Write(this.bigFileObj, data)
	if err != nil {
		prefix := fmt.Sprintf("%s-%s", this.identifier, "SqlFile.write.bigFileObj.WriteMessage")
		err = fmt.Errorf("%s-Write message to big file object failed:%s", prefix, err)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}
}

// 获取大文件对象的文件绝对路径
func (this *SqlFile) FileFullName() string {
	return filepath.Join(this.dirPath, this.bigFileObj.FileName())
}

// 当前读写的文件名
func (this *SqlFile) FileName() string {
	return this.bigFileObj.FileName()
}

// 创建同步数据对象
// _dirPath:目录的路径
// _identifier:当前数据的唯一标识（可以使用数据库表名）
// _maxFileSize:每个大文件的最大写入值（单位：Byte）
// 返回值:
// 同步数据对象
func NewSqlFile(dirPath, identifier, fileName string, maxFileSize int) *SqlFile {
	result := &SqlFile{
		dirPath:    dirPath,
		identifier: identifier,
	}

	// 初始化大文件对象
	if fileName == "" {
		fileName = con_Default_FileName
	}
	bigFileObj, err := fileUtil.NewBigFileWithNewFileNameFunc2(dirPath, "", fileName, maxFileSize, NewFileName)
	if err != nil {
		prefix := fmt.Sprintf("%s-%s", result.identifier, "SqlFile.newSqlFile.fileUtil.NewBigFileWithNewFileNameFunc")
		err = fmt.Errorf("%s-Create big file object failed:%s", prefix, err)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}
	result.bigFileObj = bigFileObj

	return result
}

// 根据当前文件名生成下一个sql文件名
// prefix:文件名前缀
// path:当前文件的路径
// 返回值:
// string:下一个文件的完整路径
func NewFileName(prefix, path string) string {
	fullName := filepath.Base(path)
	curFileName := strings.Split(fullName, ".")[0]
	curFileId, err := strconv.Atoi(curFileName)
	if err != nil {
		err = fmt.Errorf("%s-Convert newFileName:%s to int failed:%s", prefix, curFileName, err)
		logUtil.ErrorLog(err.Error())
		panic(err)
	}

	newFileId := curFileId + 1
	newFileName := fmt.Sprintf("%08d", newFileId)

	// 加上文件后缀
	newFileName = fmt.Sprintf("%s.%s", newFileName, con_FileName_Suffix)

	return newFileName
}

// 获取文件夹下所有的sql文件
// dirPath:指定要获取的文件夹路径
// 返回值:
// []string:sql文件列表
func GetDataFileList(dirPath string) []string {
	// 获取当前目录中所有的数据文件列表
	fileList, err := fileUtil.GetFileList2(dirPath, "", con_FileName_Suffix)
	if err != nil {
		if os.IsNotExist(err) {
		} else {
			err = fmt.Errorf("%s/*.%s-Get file list failed:%s", dirPath, con_FileName_Suffix, err)
			logUtil.ErrorLog(err.Error())
			panic(err)
		}
	}

	// 如果文件数量大于1，则进行排序，以便于后续处理
	if len(fileList) > 1 {
		sort.Strings(fileList)
	}

	return fileList
}
