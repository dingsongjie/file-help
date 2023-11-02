package imginfo

import (
	"www.github.com/dingsongjie/file-help/configs"
	"www.github.com/dingsongjie/file-help/pkg/file/img"
	"www.github.com/dingsongjie/file-help/pkg/s3helper"
)

type ImgInfoQueries struct {
	s3Helper s3helper.S3Helper
}

var (
	getImgInfo = img.GetImgInfo
)

func NewImgInfoQueries() *ImgInfoQueries {
	queries := ImgInfoQueries{}
	queries.s3Helper = s3helper.NewS3Helper(configs.S3Endpoint, configs.S3AccessKey, configs.S3SecretKey, configs.S3BacketName)
	return &queries
}

func (r *ImgInfoQueries) GetImgInfo(request *GetImgInfoRequest) *GetImgInfoResponse {
	s3HelperInstance := r.s3Helper
	response := NewGetImgInfoResponse()
	for _, item := range request.Items {
		bytes, err := s3HelperInstance.DownLoadAndReturnBuffer(item.FileKey)
		if err != nil {
			response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: false, Message: err.Error()})
			continue
		}
		info, err := getImgInfo(bytes)
		if err != nil {
			response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: false, Message: err.Error()})
			continue
		}
		response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: true, Width: info.Width, Height: info.Height})
	}
	return response
}
