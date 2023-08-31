package s3helper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedDownloader struct {
	mock.Mock
}

func (r *mockedDownloader) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	w.WriteAt([]byte{'A', 'B'}, 0)
	args := r.Called(w, input, options)
	return int64(args.Int(0)), args.Error(1)
}
func (r *mockedDownloader) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	args := r.Called(ctx, w, input, options)
	return int64(args.Int(0)), args.Error(1)
}

type mockedUploader struct {
	mock.Mock
}

func (r *mockedUploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	args := r.Called(input)
	var arg1 *s3manager.UploadOutput = nil
	if arg1Tmp, ok := args.Get(0).(*s3manager.UploadOutput); ok {
		arg1 = arg1Tmp
	}
	return arg1, args.Error(1)
}

func (r *mockedUploader) UploadWithContext(ctx aws.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	args := r.Called(ctx, input, options)
	output := args.Get(0).(*s3manager.UploadOutput)
	return output, args.Error(1)
}

func NewMockedS3Helper(downloader s3manageriface.DownloaderAPI, uploader s3manageriface.UploaderAPI, tempPath string) (S3Helper, error) {
	sess := safelyCreateOrGetSession("https://test.com", "test", "test", "test-bucket")
	bucketName := "test-bucket"
	instance := &S3DefaultHelper{
		bucketName: bucketName,
		session:    sess,
		downloader: downloader,
		uploader:   uploader,
		tempPath:   tempPath,
	}
	return instance, nil
}

func createFileAndReturnFile(fullPath string) (*os.File, error) {
	err := os.MkdirAll(path.Dir(fullPath), 0770)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func TestNewS3Helper(t *testing.T) {
	assert := assert.New(t)
	inter, err := NewS3Helper("https://test.com", "test", "test", "test-bucket")
	assert.Nil(err)
	instance, ok := inter.(*S3DefaultHelper)
	assert.True(ok)
	assert.Equal("test-bucket", instance.bucketName)
	assert.NotNil(instance.downloader)
	assert.NotNil(instance.uploader)
	assert.Equal(os.TempDir(), instance.tempPath)
	assert.NotNil(instance.session)
}

func TestDownLoadAndReturnLocalPath(t *testing.T) {

	dir := t.TempDir()
	assert := assert.New(t)
	key := "test/1.txt"
	fullPath := path.Join(dir, key)

	t.Run("file exist", func(t *testing.T) {
		downloader := new(mockedDownloader)
		downloader.On("Download", mock.Anything, mock.Anything, mock.Anything).Return(2, nil)
		file, err := createFileAndReturnFile(fullPath)
		assert.Nil(err)
		file.Write([]byte{'A', 'B'})
		file.Close()
		mockedS3Helper, _ := NewMockedS3Helper(downloader, nil, dir)
		handle, _ := mockedS3Helper.DownLoadAndReturnLocalPath(key)
		assert.False(handle.IsDestory)
		assert.Equal(fullPath, handle.Path)
		buffer, err := os.ReadFile(handle.Path)
		assert.Nil(err)
		assert.Equal(2, len(buffer))
		assert.Equal('A', int32(buffer[0]))
		assert.Equal('B', int32(buffer[1]))
		handle.Destory()
	})

	t.Run("file not exist", func(t *testing.T) {
		downloader := new(mockedDownloader)
		downloader.On("Download", mock.Anything, mock.Anything, mock.Anything).Return(2, nil)
		key := "test/2.txt"
		fullPath := path.Join(dir, key)
		mockedS3Helper, _ := NewMockedS3Helper(downloader, nil, dir)
		handle, _ := mockedS3Helper.DownLoadAndReturnLocalPath(key)
		assert.False(handle.IsDestory)
		assert.Equal(fullPath, handle.Path)
		buffer, err := os.ReadFile(handle.Path)
		assert.Nil(err)
		assert.Equal(2, len(buffer))
		assert.Equal('A', int32(buffer[0]))
		assert.Equal('B', int32(buffer[1]))
		handle.Destory()
	})
	t.Run("download err", func(t *testing.T) {
		downloader := new(mockedDownloader)
		downloader.On("Download", mock.Anything, mock.Anything, mock.Anything).Return(0, fmt.Errorf("test error"))
		key := "test/3.txt"
		mockedS3Helper, _ := NewMockedS3Helper(downloader, nil, dir)
		_, err := mockedS3Helper.DownLoadAndReturnLocalPath(key)
		assert.NotNil(err)
	})
	t.Run("custom part size", func(t *testing.T) {
		downloader := new(mockedDownloader)
		mockedS3Helper, _ := NewMockedS3Helper(downloader, nil, dir)
		downloader.On("Download", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			if option, ok := args[2].([]func(d *s3manager.Downloader)); ok {
				option[0](&s3manager.Downloader{})
			} else {
				fmt.Print("11")
			}

		}).Return(2, nil)
		key := "test/1.txt"
		handle, _ := mockedS3Helper.DownLoadAndReturnLocalPath(key)
		assert.False(handle.IsDestory)
		assert.Equal(fullPath, handle.Path)
		buffer, err := os.ReadFile(handle.Path)
		assert.Nil(err)
		assert.Equal(2, len(buffer))
		assert.Equal('A', int32(buffer[0]))
		assert.Equal('B', int32(buffer[1]))
		handle.Destory()
	})

	//only windows
	// t.Run("wrong dir", func(t *testing.T) {
	// 	downloader := new(mockedDownloader)
	// 	dir = os.TempDir() + "........./......./"
	// 	err := os.MkdirAll(path.Dir(fullPath), 0770)
	// 	assert.Nil(err)
	// 	file, err := os.Create(fullPath)
	// 	assert.Nil(err)
	// 	file.Write([]byte{'A', 'B'})
	// 	file.Close()
	// 	mockedS3Helper, _ := NewMockedS3Helper(downloader, nil, dir)
	// 	_, err = mockedS3Helper.DownLoadAndReturnLocalPath(key)
	// 	assert.NotNil(err)
	// })

}

func TestUpload(t *testing.T) {
	key := "test/uploader.txt"
	dir := t.TempDir()
	fullPath := path.Join(dir, key)

	assert := assert.New(t)

	t.Run("upload exist file", func(t *testing.T) {
		uploader := new(mockedUploader)
		mockedS3Helper, _ := NewMockedS3Helper(nil, uploader, dir)
		instance, ok := mockedS3Helper.(*S3DefaultHelper)
		assert.True(ok)
		err := os.MkdirAll(path.Dir(fullPath), 0770)
		assert.Nil(err)
		file, err := os.Create(fullPath)
		assert.Nil(err)
		file.Write([]byte{'E', 'F'})
		file.Close()
		uploader.On("Upload", mock.Anything).
			Return(&s3manager.UploadOutput{}, nil).
			Run(func(args mock.Arguments) {
				input := args.Get(0).(*s3manager.UploadInput)
				assert.Equal(instance.bucketName, *input.Bucket)
				assert.Equal(key, *input.Key)
				buffer := make([]byte, 100)
				input.Body.Read(buffer)

				assert.Equal('E', int32(buffer[0]))
				assert.Equal('F', int32(buffer[1]))
				fmt.Print(input)
			})

		err = mockedS3Helper.Upload(fullPath, key)
		assert.Nil(err)

		//清空测试文件
		os.Remove(fullPath)
	})
	t.Run("upload not exist file", func(t *testing.T) {
		uploader := new(mockedUploader)
		mockedS3Helper, _ := NewMockedS3Helper(nil, uploader, dir)
		instance, ok := mockedS3Helper.(*S3DefaultHelper)
		assert.True(ok)
		uploader.On("Upload", mock.Anything).
			Return(&s3manager.UploadOutput{}, nil).
			Run(func(args mock.Arguments) {
				input := args.Get(0).(*s3manager.UploadInput)
				assert.Equal(instance.bucketName, input.Bucket)
				assert.Equal(key, input.Key)
				buffer := make([]byte, 100)
				input.Body.Read(buffer)

				assert.Equal('E', int32(buffer[0]))
				assert.Equal('F', int32(buffer[1]))
				fmt.Print(input)
			})

		err := mockedS3Helper.Upload(fullPath, key)
		assert.True(errors.Is(err, os.ErrNotExist))
	})

	t.Run("upload error", func(t *testing.T) {
		uploader := new(mockedUploader)
		mockedS3Helper, _ := NewMockedS3Helper(nil, uploader, dir)
		uploader.On("Upload", mock.Anything).
			Return(nil, fmt.Errorf("test error"))
		err := os.MkdirAll(path.Dir(fullPath), 0770)
		assert.Nil(err)
		file, err := os.Create(fullPath)
		assert.Nil(err)
		file.Write([]byte{'E', 'F'})
		file.Close()
		err = mockedS3Helper.Upload(fullPath, key)
		assert.NotNil(err)
		//清空测试文件
		os.Remove(fullPath)
	})
}
