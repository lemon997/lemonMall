package service

import (
	"errors"
	"mime/multipart"

	"github.com/lemon997/lemonMall/internal/dao"

	"github.com/lemon997/lemonMall/global"

	"github.com/lemon997/lemonMall/common/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

// func (svc Service) UploadFile(fileType upload.FileType, file multipart.File,
// 	fileHeader *multipart.FileHeader) (*FileInfo, error) {
// 	//multipart.FileHeader是关于上传文件的属性
// 	fileName := upload.GetFileName(fileHeader.Filename)
// 	if !upload.CheckContainExt(fileType, fileName) {
// 		return nil, errors.New("file suffix is not supported.")
// 	}

// 	if upload.CheckExceedMaxSize(fileType, file) {
// 		return nil, errors.New("exceeded maximum file limit.")
// 	}

// 	uploadSavePath := upload.GetSavePath()

// 	if upload.CheckSavePath(uploadSavePath) {
// 		// 指定保存目录不存在则true
// 		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
// 			//目标目录权限777
// 			return nil, errors.New("failed to create save directory.")
// 		}
// 	}

// 	if upload.CheckNotPermission(uploadSavePath) {
// 		return nil, errors.New("insufficient file permissions.")
// 	}

// 	dst := uploadSavePath + "/" + fileName //最终保存的目录地址
// 	if err := upload.SaveFile(fileHeader, dst); err != nil {
// 		return nil, err
// 	}

// 	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName //URL地址
// 	return &FileInfo{
// 		Name:      fileName,
// 		AccessUrl: accessUrl,
// 	}, nil

// }

func (svc Service) UploadFile(customerID int64, fileHeader *multipart.FileHeader) (*FileInfo, error) {

	fileName := upload.GetFileName(fileHeader.Filename) //生成md5后的文件名

	if !upload.CheckContainExt(fileName) { //检查文件后缀是否过关
		return nil, errors.New("file suffix is not supported.")
	}

	if upload.CheckExceedMaxSize(fileHeader.Size) { //检查文件大小是否满足要求
		return nil, errors.New("exceeded maximum file limit.")
	}

	uploadSavePath := upload.GetSavePath() //检查上传文件的最终保存目录是否存在，不存在则创建
	if upload.CheckSavePath(uploadSavePath) {
		// 指定保存目录不存在则true
		if err := upload.CreateSavePath(uploadSavePath, 0400); err != nil {
			//目标目录权限400
			return nil, errors.New("failed to create save directory.")
		}
	}

	if upload.CheckNotPermission(uploadSavePath) { //检查保存目录权限
		return nil, errors.New("insufficient file permissions.")
	}

	dst := uploadSavePath + "/" + fileName //最终保存的目录地址
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName //URL地址

	//写入数据库
	db := dao.LoginMethod{}
	err := db.SetUrl(customerID, accessUrl)

	return &FileInfo{
		Name:      fileName,
		AccessUrl: accessUrl,
	}, err
}
