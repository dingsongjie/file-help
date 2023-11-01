package imginfo

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"www.github.com/dingsongjie/file-help/pkg/file"
	"www.github.com/dingsongjie/file-help/pkg/s3helper"
)

type mockedS3Helper struct {
	mock.Mock
}

func (r *mockedS3Helper) DownLoadAndReturnLocalPath(fileKey string) (*file.LocalFileHandle, error) {
	args := r.Called(fileKey)
	return (args.Get(0)).(*file.LocalFileHandle), args.Error(1)
}

func (r *mockedS3Helper) DownLoadAndReturnBuffer(fileKey string) ([]byte, error) {
	args := r.Called(fileKey)
	return (args.Get(0)).([]byte), args.Error(1)
}

func (r *mockedS3Helper) Upload(localPath string, fileKey string) (err error) {
	args := r.Called(localPath, fileKey)
	return args.Error(0)
}
func TestGetImgInfo(t testing.T) {
	// s3HelperInstance := s3helper.NewS3Helper(configs.S3Endpoint, configs.S3AccessKey, configs.S3SecretKey, configs.S3BacketName)
	// response := NewGetImgInfoResponse()
	// for _, item := range *request.Itemtesting.t
	// 	bytes, err := s3HelperInstance.DownLoadAndReturnBuffer(item.FileKey)
	// 	if err != nil {
	// 		response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: false, Message: err.Error()})
	// 		continue
	// 	}
	// 	info, err := img.GetimgInfo(bytes)
	// 	if err != nil {
	// 		response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: false, Message: err.Error()})
	// 		continue
	// 	}
	// 	response.AddItem(&GetImgInfoItemResponse{FileKey: item.FileKey, IsSucceed: true, Width: info.Width, Height: info.Height})
	// }
	request := GetImgInfoRequest{}
	request.Items = append(request.Items, GetImgInfoRequestItem{FileKey: "1.png"})
	newS3Helper = func() s3helper.S3Helper {
		return &mockedS3Helper{}
	}
	response := GetImgInfo(&request)
}
