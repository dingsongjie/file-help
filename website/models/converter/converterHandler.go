package converter

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/STRockefeller/go-linq"
	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/file"
	"www.github.com/dingsongjie/file-help/pkg/s3helper"
)

type GetFisrtImageByGavingKeyRequestHandler struct {
	itemHandler *GetFisrtImageByGavingKeyRequestItemHandler
}

type GetFisrtImageByGavingKeyRequestItemHandler struct {
	s3Helper    s3helper.S3Helper
	getFileSize func(filePath string) int64
}

var (
	downloadHttpFile = file.DownLoadAndReturnLocalPath
	osMakeDir        = os.MkdirAll
)

func NewGetFisrtImageByGavingKeyRequestHandler(endpoint, accessKey, secretKey, bucketName string) (*GetFisrtImageByGavingKeyRequestHandler, error) {
	s3helper := s3helper.NewS3Helper(endpoint, accessKey, secretKey, bucketName)
	itemHandler := GetFisrtImageByGavingKeyRequestItemHandler{s3Helper: s3helper, getFileSize: fileSize}
	handler := GetFisrtImageByGavingKeyRequestHandler{itemHandler: &itemHandler}
	return &handler, nil
}

func (r *GetFisrtImageByGavingKeyRequestHandler) Handle(request *ConvertByGavingKeyRequest) *ConvertByGavingKeyResponse {
	GetFisrtImageByGavingKeyResponse := NewGetFisrtImageByGavingKeyResponse()
	items := linq.Linq[ConvertByGavingKeyRequestItem]{}
	items.AddRange(request.Items)
	items.ForEach(func(gfibgki ConvertByGavingKeyRequestItem) {
		responseItem := r.itemHandler.Handle(&gfibgki)
		GetFisrtImageByGavingKeyResponse.AddItem(responseItem)
	})
	return GetFisrtImageByGavingKeyResponse
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) Handle(item *ConvertByGavingKeyRequestItem) *ConvertByGavingKeyResponseItem {
	fileSize, err := r.HandleCore(item)
	if err == nil {
		return &ConvertByGavingKeyResponseItem{SourceKey: item.SourceKey, IsSucceed: true, TargetFileSize: fileSize}
	}
	return &ConvertByGavingKeyResponseItem{SourceKey: item.SourceKey, IsSucceed: false, Message: err.Error()}
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) HandleCore(item *ConvertByGavingKeyRequestItem) (int64, error) {
	var fileSize int64 = 0

	fileHandler, err := r.downloadSourceFile(item.SourceKey)
	if err != nil {
		return fileSize, err
	}
	defer fileHandler.Destory()
	pair, err := r.validateAndGetFileConverterPair(&struct {
		SourceFileName string
		TargetFileName string
	}{SourceFileName: fileHandler.Path, TargetFileName: item.TargetKey})
	if err != nil {
		return fileSize, err
	}
	firstHandler := converter.Converters.FirstOrDefault(func(c converter.Converter) bool {
		return c.CanHandle(*pair)
	})
	if firstHandler == nil {
		return fileSize, fmt.Errorf("convert is not support, sourceType:%s,targetType:%s", pair.SourceType, pair.TargetType)
	}
	generateFilePath := path.Join(os.TempDir(), "generate", item.TargetKey)
	err = osMakeDir(path.Dir(generateFilePath), 0770)
	if err != nil {
		return fileSize, err
	}

	// 这里目前只有两种 先做简单判断
	if pair.TargetType == "pdf" {
		err = firstHandler.ToPrettyPdf(fileHandler.Path, generateFilePath)
	} else {
		err = firstHandler.ToFastImage(fileHandler.Path, generateFilePath, item.TargetFileDpi)
	}
	if err != nil {
		return 0, err
	}
	fileSize = r.getFileSize(generateFilePath)
	defer os.Remove(generateFilePath)
	err = r.uploadTargetFile(generateFilePath, item.TargetKey)
	if err != nil {
		return fileSize, err
	}
	return fileSize, nil
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) validateAndGetFileConverterPair(item *struct{ SourceFileName, TargetFileName string }) (*converter.ConverterTypePair, error) {
	sourceKeySplit := strings.Split(item.SourceFileName, ".")
	if len(sourceKeySplit) < 2 {
		return nil, fmt.Errorf("wrong sourceFileName, sourceFileName:%s", item.SourceFileName)
	}
	targetKeySplit := strings.Split(item.TargetFileName, ".")
	if len(targetKeySplit) < 2 {
		return nil, fmt.Errorf("wrong targetFileName, targetFileName:%s", item.TargetFileName)
	}
	return &converter.ConverterTypePair{SourceType: sourceKeySplit[len(sourceKeySplit)-1], TargetType: targetKeySplit[len(targetKeySplit)-1]}, nil
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) downloadSourceFile(key string) (*file.LocalFileHandle, error) {
	var (
		file *file.LocalFileHandle
		err  error
	)

	if isValidURL(key) {
		file, err = downloadHttpFile(key)
	} else {
		file, err = r.s3Helper.DownLoadAndReturnLocalPath(key)
	}

	if err != nil {
		return nil, err
	}
	return file, nil
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) uploadTargetFile(localPath, targetKey string) error {
	err := r.s3Helper.Upload(localPath, targetKey)
	if err != nil {
		return err
	}
	return nil
}

func isValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	return err == nil
}

func fileSize(filePath string) int64 {
	info, _ := os.Stat(filePath)
	return info.Size()
}
