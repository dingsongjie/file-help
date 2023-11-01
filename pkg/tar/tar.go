package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"time"
)

type TarHepler interface {
	Pack(request ExecuteContext) error
}

type ExecuteContext struct {
	FileName string
	IsGziped bool
	Items    []PackItem
}

type PackItem struct {
	FilePath, FileName string
	LastModifyTime     time.Time
}

func NewTarHepler() TarHepler {
	helper := new(DefaultTarHepler)
	return helper
}

type DefaultTarHepler struct {
}

func (r *DefaultTarHepler) Pack(request ExecuteContext) error {
	output, err := os.Create(request.FileName)
	if err != nil {
		return err
	}
	defer output.Close()
	var tarWriter *tar.Writer = nil
	if request.IsGziped {
		// 创建gzip压缩器
		gzWriter := gzip.NewWriter(output)
		defer gzWriter.Close()
		tarWriter = tar.NewWriter(gzWriter)
	} else {
		tarWriter = tar.NewWriter(output)
	}

	defer tarWriter.Close()
	for index := range request.Items {
		item := request.Items[index]
		fileInfo, err := os.Stat(item.FilePath)
		if err != nil {
			return err
		}
		// 创建tar文件头
		header := &tar.Header{
			Name: item.FileName,
			Size: fileInfo.Size(),
			//Mode:    int64(fileInfo.Mode()),
			ModTime: item.LastModifyTime,
			//Format:  tar.FormatGNU,
			// PAXRecords: map[string]string{
			// 	"filename": item.FileName, // 使用Unicode编码文件名
			// },
		}
		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(item.FilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		// 将文件内容复制到tar文件中
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}

	}
	return nil
}
