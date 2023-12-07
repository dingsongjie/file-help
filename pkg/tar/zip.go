package tar

import (
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"time"
)

type ZipHepler interface {
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

func NewZipHepler() ZipHepler {
	helper := new(DefaultZipHepler)
	return helper
}

type DefaultZipHepler struct {
}

func (r *DefaultZipHepler) Pack(request ExecuteContext) error {
	output, err := os.Create(request.FileName)
	if err != nil {
		return err
	}
	defer output.Close()
	var zipWriter *zip.Writer = nil
	if request.IsGziped {
		// 创建gzip压缩器
		gzWriter := gzip.NewWriter(output)
		defer gzWriter.Close()
		zipWriter = zip.NewWriter(gzWriter)
	} else {
		zipWriter = zip.NewWriter(output)
	}

	defer zipWriter.Close()
	for index := range request.Items {
		item := request.Items[index]
		fileInfo, err := os.Stat(item.FilePath)
		if err != nil {
			return err
		}
		// 创建tar文件头
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		header.Name = item.FileName
		header.Modified = item.LastModifyTime

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		file, err := os.Open(item.FilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		// 将文件内容复制到tar文件中
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

	}
	return nil
}
