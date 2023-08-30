package converter

import (
	"fmt"
	"testing"

	"github.com/STRockefeller/go-linq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/s3helper"
)

var (
	testAiKey        string = "img/1.ai"
	testPsdKey       string = "img/2.psd"
	testS3Endpoint   string = "http://tests3.com"
	testS3AccessKey  string = "test-accessKey"
	testS3SecretKey  string = "test-secretKey"
	testS3BucketName string = "test-bucket"
	testTempJpeg     string = "tmp/1.jpeg"
)

type mockedS3Helper struct {
	mock.Mock
}

func (r *mockedS3Helper) DownLoadAndReturnLocalPath(fileKey string) (*s3helper.LocalFileHandle, error) {
	args := r.Called(fileKey)
	return (args.Get(0)).(*s3helper.LocalFileHandle), args.Error(1)
}

func (r *mockedS3Helper) Upload(localPath string, fileKey string) (err error) {
	args := r.Called(localPath, fileKey)
	return args.Error(0)
}

type mockedAiConverter struct {
	mock.Mock
}

func (r *mockedAiConverter) ToFastImage(inputFile string, outputFile string) error {
	args := r.Called(inputFile, outputFile)
	return args.Error(0)
}

func (r *mockedAiConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	return r.ToFastImage(inputFile, outputFile)
}

func (r *mockedAiConverter) Destory() {
	r.Called()
}

func (r *mockedAiConverter) CanHandle(pair converter.ConverterTypePair) bool {
	args := r.Called(pair)
	return args.Bool(0)
}

type mockedImagickConverter struct {
	mock.Mock
}

func (r *mockedImagickConverter) ToFastImage(inputFile string, outputFile string) error {
	args := r.Called(inputFile, outputFile)
	return args.Error(0)
}

func (r *mockedImagickConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	return r.ToFastImage(inputFile, outputFile)
}

func (r *mockedImagickConverter) Destory() {
	r.Called()
}

func (r *mockedImagickConverter) CanHandle(pair converter.ConverterTypePair) bool {
	args := r.Called(pair)
	return args.Bool(0)
}

func RegisterMockedConverters() {
	converter.Converters = make([]converter.Converter, 2)
	aiConverter := new(mockedAiConverter)
	aiConverter.On("ToFastImage", mock.Anything, mock.Anything).Return(nil)
	aiConverter.On("ToPrettyPdf", mock.Anything, mock.Anything).Return(nil)
	aiConverter.On("CanHandle", mock.Anything).Return(true)
	aiConverter.On("Destory").Return(nil)
	converter.Converters[0] = aiConverter
	imagickConverter := new(mockedImagickConverter)
	imagickConverter.On("ToFastImage", mock.Anything, mock.Anything).Return(nil)
	imagickConverter.On("ToPrettyPdf", mock.Anything, mock.Anything).Return(nil)
	imagickConverter.On("CanHandle", mock.Anything).Return(true)
	imagickConverter.On("Destory").Return(nil)
	converter.Converters[1] = imagickConverter
}

func RegisterMockedFailConvertConverters() {
	converter.Converters = make([]converter.Converter, 2)
	aiConverter := new(mockedAiConverter)
	aiConverter.On("ToFastImage", mock.Anything, mock.Anything).Return(fmt.Errorf("convert to image error"))
	aiConverter.On("ToPrettyPdf", mock.Anything, mock.Anything).Return(nil)
	aiConverter.On("CanHandle", mock.Anything).Return(true)
	aiConverter.On("Destory").Return(nil)
	converter.Converters[0] = aiConverter
	imagickConverter := new(mockedImagickConverter)
	imagickConverter.On("ToFastImage", mock.Anything, mock.Anything).Return(nil)
	imagickConverter.On("ToPrettyPdf", mock.Anything, mock.Anything).Return(fmt.Errorf("convert to pdf error"))
	imagickConverter.On("CanHandle", mock.Anything).Return(true)
	imagickConverter.On("Destory").Return(nil)
	converter.Converters[1] = imagickConverter
}
func RegisterNoConverters() {
	converter.Converters = []converter.Converter{}
}

func TestNewGetFisrtImageByGavingKeyRequestHandler(t *testing.T) {
	assert := assert.New(t)
	t.Run("success", func(t *testing.T) {
		handler, err := NewGetFisrtImageByGavingKeyRequestHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
		assert.Nil(err)
		assert.NotNil(handler.itemHandler.s3Helper)
	})
}

func NewMockedGetFisrtImageByGavingKeyRequestHandler() *GetFisrtImageByGavingKeyRequestHandler {
	handler, _ := NewGetFisrtImageByGavingKeyRequestHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
	mockedS3Helper := new(mockedS3Helper)
	mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&s3helper.LocalFileHandle{Path: testTempJpeg, IsDestory: false}, nil)
	mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(nil)
	handler.itemHandler.s3Helper = mockedS3Helper
	return handler
}

func NewDownloadFaildGetFisrtImageByGavingKeyRequestHandler() *GetFisrtImageByGavingKeyRequestHandler {
	handler, _ := NewGetFisrtImageByGavingKeyRequestHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
	mockedS3Helper := new(mockedS3Helper)
	mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&s3helper.LocalFileHandle{Path: testTempJpeg, IsDestory: false}, nil).Once()
	mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&s3helper.LocalFileHandle{Path: testTempJpeg, IsDestory: false}, fmt.Errorf("key not exist")).Once()
	mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(nil)
	handler.itemHandler.s3Helper = mockedS3Helper
	return handler
}

func NewUploadFaildGetFisrtImageByGavingKeyRequestHandler() *GetFisrtImageByGavingKeyRequestHandler {
	handler, _ := NewGetFisrtImageByGavingKeyRequestHandler(testS3Endpoint, testS3AccessKey, testS3SecretKey, testS3BucketName)
	mockedS3Helper := new(mockedS3Helper)
	mockedS3Helper.On("DownLoadAndReturnLocalPath", mock.Anything).Return(&s3helper.LocalFileHandle{Path: testTempJpeg, IsDestory: false}, nil)
	mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(fmt.Errorf("upload error")).Once()
	mockedS3Helper.On("Upload", mock.Anything, mock.Anything).Return(nil).Once()
	handler.itemHandler.s3Helper = mockedS3Helper
	return handler
}

func TestGetFisrtImageByGavingKeyRequestHandlerHandle(t *testing.T) {
	assert := assert.New(t)
	t.Run("convert to image success", func(t *testing.T) {
		RegisterMockedConverters()
		handler := NewMockedGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testAiKey, TargetKey: "img/1.jpeg"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.jpeg"})
		response := handler.Handle(&request)
		assert.True(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal(testAiKey, response.Items[0].SourceKey)
		assert.True(response.Items[0].IsSucceed)
		assert.Empty(response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.True(response.Items[1].IsSucceed)
		assert.Empty(response.Items[1].Message)
	})

	t.Run("convert to pdf success", func(t *testing.T) {
		RegisterMockedConverters()
		handler := NewMockedGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testAiKey, TargetKey: "img/1.pdf"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.pdf"})
		response := handler.Handle(&request)
		assert.True(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal(testAiKey, response.Items[0].SourceKey)
		assert.True(response.Items[0].IsSucceed)
		assert.Empty(response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.True(response.Items[1].IsSucceed)
		assert.Empty(response.Items[1].Message)
	})

	t.Run("download faild", func(t *testing.T) {
		RegisterMockedConverters()
		handler := NewDownloadFaildGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testAiKey, TargetKey: "img/1.jpeg"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.jpeg"})
		response := handler.Handle(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal(testAiKey, response.Items[0].SourceKey)
		assert.True(response.Items[0].IsSucceed)
		assert.Empty(response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.False(response.Items[1].IsSucceed)
		assert.Equal("key not exist", response.Items[1].Message)
	})

	t.Run("upload faild", func(t *testing.T) {
		RegisterMockedConverters()
		handler := NewUploadFaildGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testAiKey, TargetKey: "img/1.jpeg"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.jpeg"})
		response := handler.Handle(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal(testAiKey, response.Items[0].SourceKey)
		assert.False(response.Items[0].IsSucceed)
		assert.Equal("upload error", response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.True(response.Items[1].IsSucceed)
		assert.Empty(response.Items[1].Message)
	})

	t.Run("source fileKey wrong", func(t *testing.T) {
		RegisterMockedConverters()
		handler := NewDownloadFaildGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: "img/1ai", TargetKey: "img/1.jpeg"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.jpeg"})
		response := handler.Handle(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal("img/1ai", response.Items[0].SourceKey)
		assert.False(response.Items[0].IsSucceed)
		assert.Equal("wrong sourceKey, sourceKey:img/1ai", response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.True(response.Items[1].IsSucceed)
		assert.Empty(response.Items[1].Message)
	})

	t.Run("target fileKey wrong", func(t *testing.T) {
		RegisterMockedConverters()
		handler := NewDownloadFaildGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testAiKey, TargetKey: "img/1jpeg"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.jpeg"})
		response := handler.Handle(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal(testAiKey, response.Items[0].SourceKey)
		assert.False(response.Items[0].IsSucceed)
		assert.Equal("wrong targetKey, targetKey:img/1jpeg", response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.True(response.Items[1].IsSucceed)
		assert.Empty(response.Items[1].Message)
	})

	t.Run("converter not found", func(t *testing.T) {
		RegisterNoConverters()
		handler := NewDownloadFaildGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testAiKey, TargetKey: "img/1.jpeg"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.jpeg"})
		response := handler.Handle(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal(testAiKey, response.Items[0].SourceKey)
		assert.False(response.Items[0].IsSucceed)
		assert.Equal("convert is not support, sourceType:ai,targetType:jpeg", response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.False(response.Items[1].IsSucceed)
		assert.Equal("convert is not support, sourceType:psd,targetType:jpeg", response.Items[1].Message)
	})

	t.Run("convert faild", func(t *testing.T) {
		RegisterMockedFailConvertConverters()
		handler := NewMockedGetFisrtImageByGavingKeyRequestHandler()
		request := ConvertByGavingKeyRequest{Items: linq.Linq[ConvertByGavingKeyRequestItem]{}}
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testAiKey, TargetKey: "img/1.jpeg"})
		request.Items = append(request.Items, ConvertByGavingKeyRequestItem{SourceKey: testPsdKey, TargetKey: "img/2.jpeg"})
		response := handler.Handle(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(response.Items))
		assert.Equal(testAiKey, response.Items[0].SourceKey)
		assert.False(response.Items[0].IsSucceed)
		assert.Equal("convert to image error", response.Items[0].Message)
		assert.Equal(testPsdKey, response.Items[1].SourceKey)
		assert.False(response.Items[1].IsSucceed)
		assert.Equal("convert to image error", response.Items[1].Message)
	})
}
