package tar

import (
	"errors"
	"testing"
	"time"

	"github.com/dingsongjie/file-help/pkg/file"
	"github.com/dingsongjie/file-help/pkg/mocks"
	"github.com/dingsongjie/file-help/pkg/tar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testS3Endpoint   string = "http://tests3.com"
	testS3AccessKey  string = "test-accessKey"
	testS3SecretKey  string = "test-secretKey"
	testS3BucketName string = "test-bucket"
	filePath1        string = "tmp/1.png"
	filePath2        string = "tmp/3.png"
	filePath3        string = "tmp/100.png"
	filePath4        string = "tmp/17.png"
	fileKey1         string = "/key/1.png"
	fileKey2         string = "/key/3.tar"
	fileKey3         string = "/tmp/100.png"
	fileKey4         string = "/tmp/17.png"
	fileName1        string = "/subpath/1.png"
	fileName2        string = "/100/3.tar"
	fileName3        string = "100.png"
	fileName4        string = "17.png"
)

type mockedTarHelper struct {
	mock.Mock
}

func (r *mockedTarHelper) Pack(request tar.ExecuteContext) error {
	args := r.Called(request)
	return args.Error(0)
}

func TestNewPackHandler(t *testing.T) {
	assert := assert.New(t)
	//tempDir := t.TempDir()

	t.Run("new handler", func(t *testing.T) {
		got, err := NewPackHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
		assert.NotNil(got)
		assert.Nil(err)
	})
}
func TestPackHandlerHandle(t *testing.T) {
	assert := assert.New(t)

	t.Run("success handle", func(t *testing.T) {
		handler, _ := NewPackHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
		mockedS3Helper := new(mocks.MockedS3Helper)
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath1, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath2, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath3, IsDestory: false}, nil).Once()
		//mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath4, IsDestory: false}, nil).Once()
		mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(nil)
		handler.s3Helper = mockedS3Helper
		mockedTarHelper := new(mockedTarHelper)
		mockedTarHelper.On("Pack", mock.Anything).Return(nil)
		handler.tarHelper = mockedTarHelper
		downloadHttpFile = func(url, fileName string) (*file.LocalFileHandle, error) {
			return &file.LocalFileHandle{Path: filePath4, IsDestory: false}, nil
		}

		request := PackRequest{FileKey: "/path1/1.tar", IsGziped: false}
		var items []PackRequestItem
		items = append(items, PackRequestItem{FileKey: fileKey1, FileName: fileName1, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey2, FileName: fileName2, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey3, FileName: fileName3, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: "http://www.baidu.com/17.png", FileName: fileName4, LastModifyTime: time.Now()})
		request.Items = &items
		response := handler.Handle(&request)
		assert.True(response.IsSuccessd)
		assert.Empty(response.Message)
		downloadHttpFile = file.DownLoadAndReturnLocalPath
	})

	t.Run("download faild handle", func(t *testing.T) {
		handler, _ := NewPackHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
		mockedS3Helper := new(mocks.MockedS3Helper)
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath1, IsDestory: false}, errors.New("download faild")).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath2, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath3, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath4, IsDestory: false}, nil).Once()
		mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(nil)
		handler.s3Helper = mockedS3Helper
		mockedTarHelper := new(mockedTarHelper)
		mockedTarHelper.On("Pack", mock.Anything).Return(nil)
		handler.tarHelper = mockedTarHelper
		request := PackRequest{FileKey: "/path1/1.tar", IsGziped: false}
		var items []PackRequestItem
		items = append(items, PackRequestItem{FileKey: fileKey1, FileName: fileName1, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey2, FileName: fileName2, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey3, FileName: fileName3, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey4, FileName: fileName4, LastModifyTime: time.Now()})
		request.Items = &items
		response := handler.Handle(&request)
		assert.False(response.IsSuccessd)
		assert.Equal("download faild", response.Message)
	})

	t.Run("upload faild handle", func(t *testing.T) {
		handler, _ := NewPackHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
		mockedS3Helper := new(mocks.MockedS3Helper)
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath1, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath2, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath3, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath4, IsDestory: false}, nil).Once()
		mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(errors.New("upload faild"))
		handler.s3Helper = mockedS3Helper
		mockedTarHelper := new(mockedTarHelper)
		mockedTarHelper.On("Pack", mock.Anything).Return(nil)
		handler.tarHelper = mockedTarHelper
		request := PackRequest{FileKey: "/path1/1.tar", IsGziped: false}
		var items []PackRequestItem
		items = append(items, PackRequestItem{FileKey: fileKey1, FileName: fileName1, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey2, FileName: fileName2, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey3, FileName: fileName3, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey4, FileName: fileName4, LastModifyTime: time.Now()})
		request.Items = &items
		response := handler.Handle(&request)
		assert.False(response.IsSuccessd)
		assert.Equal("upload faild", response.Message)
	})

	t.Run("pack faild", func(t *testing.T) {
		handler, _ := NewPackHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
		mockedS3Helper := new(mocks.MockedS3Helper)
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath1, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath2, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath3, IsDestory: false}, nil).Once()
		mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&file.LocalFileHandle{Path: filePath4, IsDestory: false}, nil).Once()
		mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(nil)
		handler.s3Helper = mockedS3Helper
		mockedTarHelper := new(mockedTarHelper)
		mockedTarHelper.On("Pack", mock.Anything).Return(errors.New("pack error"))
		handler.tarHelper = mockedTarHelper
		request := PackRequest{FileKey: "/path1/1.tar", IsGziped: false}
		var items []PackRequestItem
		items = append(items, PackRequestItem{FileKey: fileKey1, FileName: fileName1, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey2, FileName: fileName2, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: fileKey3, FileName: fileName3, LastModifyTime: time.Now()})
		items = append(items, PackRequestItem{FileKey: "wrong.key", FileName: fileName4, LastModifyTime: time.Now()})
		request.Items = &items
		response := handler.Handle(&request)
		assert.False(response.IsSuccessd)
		assert.Equal("pack error", response.Message)
	})
}

func BenchmarkPackHandlerHandle(b *testing.B) {
	handler, _ := NewPackHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
	mockedS3Helper := new(mocks.MockedS3Helper)
	mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Run(func(args mock.Arguments) {
		time.Sleep(300 * time.Microsecond)
	}).Return(&file.LocalFileHandle{Path: filePath1, IsDestory: false}, nil)
	mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		time.Sleep(400 * time.Microsecond)
	}).Return(nil)
	handler.s3Helper = mockedS3Helper
	mockedTarHelper := new(mockedTarHelper)
	mockedTarHelper.On("Pack", mock.Anything).Run(func(args mock.Arguments) {
		time.Sleep(10 * time.Microsecond)
	}).Return(nil)
	handler.tarHelper = mockedTarHelper

	request := PackRequest{FileKey: "/path1/1.tar", IsGziped: false}
	var items []PackRequestItem
	items = append(items, PackRequestItem{FileKey: fileKey1, FileName: fileName1, LastModifyTime: time.Now()})
	items = append(items, PackRequestItem{FileKey: fileKey2, FileName: fileName2, LastModifyTime: time.Now()})
	items = append(items, PackRequestItem{FileKey: fileKey3, FileName: fileName3, LastModifyTime: time.Now()})
	items = append(items, PackRequestItem{FileKey: fileKey4, FileName: fileName4, LastModifyTime: time.Now()})
	request.Items = &items
	handler.Handle(&request)
}
