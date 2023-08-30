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

func (r GetFisrtImageByGavingKeyRequestHandler) Handle(request *ConvertByGavingKeyRequest) *ConvertByGavingKeyResponse {
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
	err := r.HandleCore(item)
	if err == nil {
		return &ConvertByGavingKeyResponseItem{SourceKey: item.sourceKey, IsSucceed: true}
	}
	return &ConvertByGavingKeyResponseItem{SourceKey: item.sourceKey, IsSucceed: false, Message: err.Error()}
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) HandleCore(item *ConvertByGavingKeyRequestItem) error {
	pair, err := r.validateAndGetFileConverterPair(item)
	if err != nil {
		return err
	}
	firstHandler := converter.Converters.FirstOrDefault(func(c converter.Converter) bool {
		return c.CanHandle(*pair)
	})
	if firstHandler == nil {
		return fmt.Errorf("convert is not support, sourceType:%s,targetType:%s", pair.SourceType, pair.TargetType)
	}
	fileHandler, err := r.downloadSourceFile(item.sourceKey)
	if err != nil {
		return err
	}
	defer fileHandler.Destory()
	generateFilePath := path.Join(os.TempDir(), "generate", item.targetKey)
	// 这里目前只有两种 先做简单判断
	if pair.TargetType == "pdf" {
		err := firstHandler.ToPrettyPdf(fileHandler.Path, generateFilePath)
		if err != nil {
			return err
		}
	} else {
		err := firstHandler.ToFastImage(fileHandler.Path, generateFilePath)
		if err != nil {
			return err
		}
	}
	err = r.uploadTargetFile(generateFilePath, item.targetKey)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetFisrtImageByGavingKeyRequestItemHandler) validateAndGetFileConverterPair(item *ConvertByGavingKeyRequestItem) (*converter.ConverterTypePair, error) {
	sourceKeySplit := strings.Split(item.sourceKey, ".")
	if len(sourceKeySplit) != 2 {
		return nil, fmt.Errorf("wrong sourceKey, sourceKey:%s", item.sourceKey)
	}
	targetKeySplit := strings.Split(item.targetKey, ".")
	if len(targetKeySplit) != 2 {
		return nil, fmt.Errorf("wrong targetKey, targetKey:%s", item.targetKey)
	}
	return &converter.ConverterTypePair{SourceType: sourceKeySplit[1], TargetType: targetKeySplit[1]}, nil
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
