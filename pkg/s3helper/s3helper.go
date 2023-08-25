package s3helper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

var (
	sess client.ConfigProvider
	mu   sync.Mutex
)

type S3Helper struct {
	bucketName string
	session    *client.ConfigProvider
	downloader s3manageriface.DownloaderAPI
	uploader   s3manageriface.UploaderAPI
	tempPath   string
}

func NewS3Helper(endpoint, secretKey, accessKey, bucketName string) (*S3Helper, error) {
	sess := safelyCreateOrGetSession(endpoint, secretKey, accessKey, bucketName)
	instance := &S3Helper{
		bucketName: bucketName,
		session:    sess,
		downloader: s3manager.NewDownloader(*sess),
		uploader:   s3manager.NewUploader(*sess),
		tempPath:   os.TempDir(),
	}
	return instance, nil
}

func safelyCreateOrGetSession(endpoint, secretKey, accessKey, bucketName string) *client.ConfigProvider {
	if sess == nil {
		mu.Lock()
		if sess == nil {
			instance, err := session.NewSession(&aws.Config{
				Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
				Endpoint:         aws.String(endpoint),
				DisableSSL:       aws.Bool(false),
				S3ForcePathStyle: aws.Bool(true),
			})
			sess = session.Must(instance, err)
		}
	}
	return &sess
}

func (r *S3Helper) DownLoadAndReturnLocalPath(fileKey string) (*LocalFileHandle, error) {
	filename := path.Join(r.tempPath, fileKey)
	var (
		f   *os.File = nil
		err error
	)
	if !isFileExist(filename) {
		err = os.MkdirAll(path.Dir(filename), 0770)
		if err != nil {
			return nil, err
		}
		// 这里的err不可能发生，所以忽略错误返回
		f, _ = os.Create(filename)
		defer f.Close()
	} else {
		handle := LocalFileHandle{Path: filename, IsDestory: false}
		return &handle, nil
	}
	_, err = r.downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileKey),
	}, func(d *s3manager.Downloader) {
		d.PartSize = 1024 * 1024 * 32
	})
	if err != nil {
		return nil, err
	}
	handle := LocalFileHandle{Path: filename, IsDestory: false}
	return &handle, nil
}

func (r S3Helper) Upload(localPath string, fileKey string) (err error) {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileKey),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	return nil
}

func isFileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
