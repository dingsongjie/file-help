package imginfo

import (
	"www.github.com/dingsongjie/file-help/configs"
	"www.github.com/dingsongjie/file-help/pkg/file/img"
	"www.github.com/dingsongjie/file-help/pkg/s3helper"
)

var newS3Helper = func() s3helper.S3Helper {
	helper := s3helper.NewS3Helper(configs.S3Endpoint, configs.S3AccessKey, configs.S3SecretKey, configs.S3BacketName)
	return helper
}

func GetImgInfo(request *GetImgInfoRequest) *GetImgInfoResponse {
	s3HelperInstance := newS3Helper()
	response := NewGetImgInfoResponse()
	for _, item := range request.Items {
		bytes, err := s3HelperInstance.DownLoadAndReturnBuffer(item.FileKey)
		if err != nil {
			response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: false, Message: err.Error()})
			continue
		}
		info, err := img.GetimgInfo(bytes)
		if err != nil {
			response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: false, Message: err.Error()})
			continue
		}
		response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: true, Width: info.Width, Height: info.Height})
	}
	return response
}
