package tar

import (
	"os"
	"path"

	"github.com/STRockefeller/go-linq"
	"github.com/google/uuid"
	fileHelper "www.github.com/dingsongjie/file-help/pkg/file"
	"www.github.com/dingsongjie/file-help/pkg/s3helper"
	"www.github.com/dingsongjie/file-help/pkg/tar"
	"www.github.com/dingsongjie/file-help/website/models"
)

type PackHandler struct {
	s3Helper  s3helper.S3Helper
	tarHelper tar.TarHepler
}

type PackHandlerInternalModel struct {
	file    *fileHelper.LocalFileHandle
	fileKey string
}

func NewPackHandler(endpoint, accessKey, secretKey, bucketName string) (*PackHandler, error) {
	s3helper := s3helper.NewS3Helper(endpoint, accessKey, secretKey, bucketName)
	tarHelper := tar.NewTarHepler()

	handler := PackHandler{s3Helper: s3helper, tarHelper: tarHelper}
	return &handler, nil
}

func (r *PackHandler) Handle(request *PackRequest) *models.CommandResponse {
	response := models.CommandResponse{IsSuccessd: true}
	files, err := r.concurrentDownload(request.Items)
	defer func() {
		if files != nil {
			for index := range *files {
				if file := (*files)[index].file; file != nil {
					file.Destory()
				}
			}
		}
	}()
	if err != nil {
		response.IsSuccessd = false
		response.Message = err.Error()
		return &response
	}
	tarPath := path.Join(os.TempDir(), uuid.New().String())
	context := tar.ExecuteContext{FileName: tarPath, IsGziped: request.IsGziped}
	linqItems := linq.Linq[PackRequestItem]{}
	linqItems.AddRange(*request.Items)
	for index := range *files {
		file := (*files)[index]
		currentItem := linqItems.FirstOrDefault(func(pi PackRequestItem) bool {
			return pi.FileKey == file.fileKey
		})
		//这里不可能匹配不到
		// if currentItem == (PackRequestItem{}) {
		// 	err := fmt.Sprintf("can not found item in requestItem by filekey, key:%s", file.fileKey)
		// 	log.Logger.Error(err)
		// 	response.IsSuccessd = false
		// 	response.Message = err
		// 	return &response
		// }
		context.Items = append(context.Items, tar.PackItem{FilePath: file.file.Path, FileName: currentItem.FileName,
			LastModifyTime: currentItem.LastModifyTime})
	}
	err = r.tarHelper.Pack(context)
	if err != nil {
		response.IsSuccessd = false
		response.Message = err.Error()
		return &response
	}
	defer os.Remove(tarPath)

	err = r.s3Helper.Upload(tarPath, request.FileKey)
	if err != nil {
		response.IsSuccessd = false
		response.Message = err.Error()
		return &response
	}

	return &response
}

func (r *PackHandler) concurrentDownload(items *[]PackRequestItem) (*[]*PackHandlerInternalModel, error) {
	wantedCount := len(*items)
	var (
		results []*PackHandlerInternalModel
	)
	channel := make(chan struct {
		file    *fileHelper.LocalFileHandle
		fileKey string
		err     error
	}, wantedCount)
	for _, item := range *items {
		go func(current PackRequestItem) {
			file, err := r.s3Helper.DownLoadAndReturnLocalPath(current.FileKey)
			if err != nil {
				channel <- struct {
					file    *fileHelper.LocalFileHandle
					fileKey string
					err     error
				}{file: nil, fileKey: current.FileKey, err: err}
			} else {
				channel <- struct {
					file    *fileHelper.LocalFileHandle
					fileKey string
					err     error
				}{file: file, fileKey: current.FileKey, err: nil}
			}
		}(item)
	}
	for i := 0; i < wantedCount; i++ {
		current := <-channel
		if current.err != nil {
			return nil, current.err
		}
		results = append(results, &PackHandlerInternalModel{file: current.file, fileKey: current.fileKey})
	}
	return &results, nil
}
