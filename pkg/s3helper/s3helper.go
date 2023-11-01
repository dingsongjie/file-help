package s3helper

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"www.github.com/dingsongjie/file-help/pkg/file"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

var (
	sess client.ConfigProvider
	mu   sync.Mutex
)

type S3Helper interface {
	DownLoadAndReturnLocalPath(fileKey string) (*file.LocalFileHandle, error)
	Upload(localPath string, fileKey string) (err error)
	DownLoadAndReturnBuffer(fileKey string) ([]byte, error)
}

type S3DefaultHelper struct {
	bucketName string
	session    *client.ConfigProvider
	downloader s3manageriface.DownloaderAPI
	uploader   s3manageriface.UploaderAPI
	tempPath   string
}

func NewS3Helper(endpoint, accessKey, secretKey, bucketName string) S3Helper {
	sess := safelyCreateOrGetSession(endpoint, accessKey, secretKey, bucketName)
	instance := &S3DefaultHelper{
		bucketName: bucketName,
		session:    sess,
		downloader: s3manager.NewDownloader(*sess),
		uploader:   s3manager.NewUploader(*sess),
		tempPath:   os.TempDir(),
	}
	return instance
}

func safelyCreateOrGetSession(endpoint, accessKey, secretKey, bucketName string) *client.ConfigProvider {
	if sess == nil {
		mu.Lock()
		defaultRegion := "us-east-1"
		if sess == nil {
			instance, err := session.NewSession(&aws.Config{
				Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
				Endpoint:         aws.String(endpoint),
				DisableSSL:       aws.Bool(false),
				S3ForcePathStyle: aws.Bool(true),
				Region:           &defaultRegion,
			})
			sess = session.Must(instance, err)
		}
	}
	return &sess
}

func (r *S3DefaultHelper) DownLoadAndReturnLocalPath(fileKey string) (*file.LocalFileHandle, error) {
	filename := path.Join(r.tempPath, fileKey)
	var (
		f   *os.File = nil
		err error
	)
	if !file.IsFileExist(filename) {
		f, _ = file.CreateFileAndReturnFile(filename)
		defer f.Close()
	} else {
		handle := file.LocalFileHandle{Path: filename, IsDestory: false}
		return &handle, nil
	}
	_, err = r.downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileKey),
	}, func(d *s3manager.Downloader) {
		d.PartSize = 1024 * 1024 * 32
	})
	if err != nil {
		err2 := os.Remove(filename)
		fmt.Print(err2)
		return nil, err
	}
	handle := file.LocalFileHandle{Path: filename, IsDestory: false}
	return &handle, nil
}

func (r *S3DefaultHelper) DownLoadAndReturnBuffer(fileKey string) ([]byte, error) {
	buff := &aws.WriteAtBuffer{}
	_, err := r.downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (r S3DefaultHelper) Upload(localPath string, fileKey string) (err error) {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()
	contentType, reader := getContentType(f)
	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(fileKey),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	return nil
}

func getContentType(file *os.File) (string, *bytes.Reader) {
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	return fileType, fileBytes
}
