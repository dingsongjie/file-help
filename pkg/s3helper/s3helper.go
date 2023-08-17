package s3helper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"www.github.com/dingsongjie/file-help/configs"
)

var (
	instance *S3Helper
	mu       sync.Mutex
	tempPath string = "/tmp"
)

type S3Helper struct {
	endPoint   string
	accessKey  string
	secretKey  string
	bucketName string
	session    *session.Session
}

func NewS3Helper() (*S3Helper, error) {
	if instance == nil {
		mu.Lock()
		if instance == nil {
			instance = &S3Helper{
				endPoint:   configs.S3ServiceUrl,
				secretKey:  configs.S3SecretKey,
				accessKey:  configs.S3AccessKey,
				bucketName: configs.S3BacketName,
			}
			sess, err := session.NewSession(&aws.Config{
				Credentials:      credentials.NewStaticCredentials(instance.accessKey, instance.secretKey, ""),
				Endpoint:         aws.String(instance.endPoint),
				DisableSSL:       aws.Bool(false),
				S3ForcePathStyle: aws.Bool(true),
			})
			if err != nil {
				return nil, err
			}
			instance.session = session.Must(sess, err)
		}
	}
	return instance, nil
}

func (r S3Helper) DownLoadAndReturnLocalPath(fileKey string) (string, error) {
	filename := path.Join(tempPath, fileKey)
	var (
		f   *os.File = nil
		err error
	)
	if !isFileExist(filename) {
		f, err = os.Create(filename)
		if err != nil {
			return "", err
		}
	} else {
		return filename, nil
	}

	downloader := s3manager.NewDownloader(r.session)
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileKey),
	}, func(d *s3manager.Downloader) {
		d.PartSize = 1024 * 1024 * 25
	})
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (r S3Helper) Upload(localPath string, fileKey string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", localPath, err)
	}

	downloader := s3manager.NewUploader(r.session)
	_, err = downloader.Upload(&s3manager.UploadInput{
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
