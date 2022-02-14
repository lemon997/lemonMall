package uploading

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"

	"github.com/lemon997/lemonMall/common/errcode"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	svc := service.New(c.Request.Context())
	ctx := c.Request.Context()

	file, err := c.FormFile("file")
	if err != nil {
		// global.Logger.Infof(ctx, "uploading.file.UploadFile, c.FormFile,err: ", err)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error())) //返回err
		return
	}

	customerId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	fileInfo, err := svc.UploadFile(customerId.(int64), file)
	if err != nil {
		global.Logger.Errorf(ctx, "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"status":          0,
		"msg":             "文件上传成功",
		"file_access_url": fileInfo.AccessUrl,
	})

}

// func (u Upload) UploadFile(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	response := app.NewResponse(c)

// 	file, fileHeader, err := c.Request.FormFile("file")
// 	if err != nil {
// 		global.Logger.Infof(ctx, "uploading.file.UploadFile,c.Request.FormFile,err: ", err)
// 		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error())) //返回err
// 		return
// 	}

// 	fileType := convert.StrTo(c.PostForm("type")).MustInt()
// 	if fileHeader == nil || fileType <= 0 {
// 		global.Logger.Infof(ctx, "uploading.file.UploadFile,c.PostForm, fileType: %d", fileType)
// 		response.ToErrorResponse(errcode.InvalidParams)
// 		return
// 	}

// 	svc := service.New(c.Request.Context())
// 	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
// 	if err != nil {
// 		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
// 		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
// 		return
// 	}

// 	response.ToResponse(gin.H{
// 		"file_access_url": fileInfo.AccessUrl,
// 	})

// }
