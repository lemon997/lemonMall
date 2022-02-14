package upload

import (
	"io"
	// "io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/lemon997/lemonMall/common/util"
	"github.com/lemon997/lemonMall/global"
)

type FileType int

const TypeImage FileType = iota + 1

func GetFileName(name string) string {
	ext := GetFileExt(name) //识别待去除后缀
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

func GetFileExt(name string) string {
	return filepath.Ext(name)
}

func GetSavePath() string {
	//返回配置文档指定的保存目录
	return global.AppSetting.UploadSavePath
}

func CheckSavePath(dst string) bool {
	//检查文件是否存在，不存在则true
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

// func CheckContainExt(t FileType, name string) bool {
// 	//检查后缀是否在允许的范围
// 	ext := GetFileExt(name)
// 	ext = strings.ToUpper(ext) //将后缀转成大写
// 	switch t {
// 	case TypeImage:
// 		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
// 			//找到对应的后缀
// 			if strings.ToUpper(allowExt) == ext {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func CheckContainExt(name string) bool {
	//检查后缀是否在允许的范围
	ext := GetFileExt(name)
	ext = strings.ToUpper(ext) //将后缀转成大写
	for _, allowExt := range global.AppSetting.UploadImageAllowExts {
		//找到对应的后缀
		if strings.ToUpper(allowExt) == ext {
			return true
		}
	}
	return false
}

func CheckExceedMaxSize(size int64) bool {
	//检测文件大小是否超出最大大小限制
	//multipart是关于上传文件的官方包
	// content, _ := ioutil.ReadAll(f) //读取f文件,io.Reader实现File接口
	// size := len(content)
	// switch t {
	// case TypeImage:
	// 	if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
	// 		return true
	// 	}
	// }
	if size >= int64(global.AppSetting.UploadImageMaxSize*1<<20) {
		return true
	}
	return false
}

func CheckNotPermission(dst string) bool {
	//检查文件权限是否足够，有权限则放回false,没有权限返回true
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	//filemode是文件的模式和权限位置，比如mode.is.dir
	//创建上传文件时所需要的保存目录，目标目录存在则不进行操作，返回nil
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	//保存文件
	src, err := file.Open() //打开需要上传的文件

	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst) //打开准备写入的文件
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src) //复制进去
	return err
}
