package converter

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/STRockefeller/go-linq"
	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/s3helper"
)

type GetFisrtImageByGavingKeyRequestHandler struct {
	itemHandler *GetFisrtImageByGavingKeyRequestItemHandler
}

type GetFisrtImageByGavingKeyRequestItemHandler struct {
	s3Helper s3helper.S3Helper
}

func NewGetFisrtImageByGavingKeyRequestHandler(endpoint, accessKey, secretKey, bucketName string) (*GetFisrtImageByGavingKeyRequestHandler, error) {
	s3helper, err := s3helper.NewS3Helper(endpoint, accessKey, secretKey, bucketName)
	if err != nil {
		return nil, err
	}
	itemHandler := GetFisrtImageByGavingKeyRequestItemHandler{s3Helper: s3helper}
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
	pair, err := r.validateAndGetFileConverterPair(item)
	var fileSize int64 = 0
	if err != nil {
		return fileSize, err
	}
	firstHandler := converter.Converters.FirstOrDefault(func(c converter.Converter) bool {
		return c.CanHandle(*pair)
	})
	if firstHandler == nil {
		return fileSize, fmt.Errorf("convert is not support, sourceType:%s,targetType:%s", pair.SourceType, pair.TargetType)
	}
	fileHandler, err := r.downloadSourceFile(item.SourceKey)
	if err != nil {
		return fileSize, err
	}
	defer fileHandler.Destory()
	generateFilePath := path.Join(os.TempDir(), "generate", item.TargetKey)
	err = os.MkdirAll(path.Dir(generateFilePath), 0770)
	if err != nil {
		return fileSize, err
	}
	// 这里目前只有两种 先做简单判断
	if pair.TargetType == "pdf" {
		err = firstHandler.ToPrettyPdf(fileHandler.Path, generateFilePath)
	} else {
		err = firstHandler.ToFastImage(fileHandler.Path, generateFilePath)
	}
	info, _ := os.Stat(generateFilePath)
	fileSize = info.Size()
	defer os.Remove(generateFilePath)
	if err != nil {
		return 0, err
	}
	err = r.uploadTargetFile(generateFilePath, item.TargetKey)
	if err != nil {
		return fileSize, err
	}
	return fileSize, nil
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) validateAndGetFileConverterPair(item *ConvertByGavingKeyRequestItem) (*converter.ConverterTypePair, error) {
	sourceKeySplit := strings.Split(item.SourceKey, ".")
	if len(sourceKeySplit) < 2 {
		return nil, fmt.Errorf("wrong sourceKey, sourceKey:%s", item.SourceKey)
	}
	targetKeySplit := strings.Split(item.TargetKey, ".")
	if len(targetKeySplit) < 2 {
		return nil, fmt.Errorf("wrong targetKey, targetKey:%s", item.TargetKey)
	}
	return &converter.ConverterTypePair{SourceType: sourceKeySplit[len(sourceKeySplit)-1], TargetType: targetKeySplit[len(targetKeySplit)-1]}, nil
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) downloadSourceFile(key string) (*s3helper.LocalFileHandle, error) {
	file, err := r.s3Helper.DownLoadAndReturnLocalPath(key)
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
