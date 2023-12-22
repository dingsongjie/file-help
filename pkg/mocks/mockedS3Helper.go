package mocks

import (
	"github.com/dingsongjie/file-help/pkg/file"
	"github.com/stretchr/testify/mock"
)

type MockedS3Helper struct {
	mock.Mock
}

func (r *MockedS3Helper) DownLoadAndReturnLocalPath(fileKey string) (*file.LocalFileHandle, error) {
	args := r.Called(fileKey)
	return (args.Get(0)).(*file.LocalFileHandle), args.Error(1)
}

func (r *MockedS3Helper) DownLoadAndReturnBuffer(fileKey string) ([]byte, error) {
	args := r.Called(fileKey)
	if instance, ok := args.Get(0).([]byte); ok {
		return instance, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (r *MockedS3Helper) Upload(localPath string, fileKey string) (err error) {
	args := r.Called(localPath, fileKey)
	return args.Error(0)
}
