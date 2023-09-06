package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTarHeplerPack(t *testing.T) {
	assert := assert.New(t)
	tempDir := t.TempDir()
	mydir, _ := os.Getwd()
	helper := NewTarHepler()
	files := []PackItem{
		{FilePath: path.Join(mydir, "assets/《个人防疫手册（第二版）》.pdf"),
			FileName:       path.Join("inner/", "《个人防疫手册（第二版）》.pdf"),
			LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
		{FilePath: path.Join(mydir, "assets/导出外贸ERP.ps1"),
			FileName:       path.Join("inner/", "导出外贸ERP.ps1"),
			LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
		{FilePath: path.Join(mydir, "assets/流程排错解决方案.xmind"),
			FileName:       "流程排错解决方案.xmind",
			LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
		{FilePath: path.Join(mydir, "assets/外销合同.txt"),
			FileName:       "inner22/外销合同.txt",
			LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
		{FilePath: path.Join(mydir, "assets/文档.docx"),
			FileName:       "inner22/文档.docx",
			LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
		{FilePath: path.Join(mydir, "assets/文档记录下载1690254313410.zip"),
			FileName:       "inner22/文档记录下载1690254313410.zip",
			LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
	}
	t.Run("sample tar", func(t *testing.T) {
		fileName := path.Join(tempDir, "001.tar")
		request := ExecuteContext{FileName: fileName, IsGziped: false}
		request.Items = append(request.Items, files...)

		err := helper.Pack(request)
		assert.Nil(err)
		_, err = os.Stat(fileName)
		assert.Nil(err)
		file, err := os.Open(fileName)
		assert.Nil(err)
		defer file.Close()
		tarReader := tar.NewReader(file)
		iterations := 0
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				// 已到达存档末尾
				break
			}
			currentFile := files[iterations]
			assert.Equal(currentFile.FileName, header.Name)
			assert.NotEqual(0, header.Size)
			assert.Equal(currentFile.LastModifyTime, time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC))
			iterations++
		}
	})

	t.Run("sample gzip", func(t *testing.T) {
		fileName := path.Join(tempDir, "001.tar.gz")
		request := ExecuteContext{FileName: fileName, IsGziped: true}
		request.Items = append(request.Items, files...)

		err := helper.Pack(request)
		assert.Nil(err)
		_, err = os.Stat(fileName)
		assert.Nil(err)
		file, err := os.Open(fileName)
		assert.Nil(err)
		defer file.Close()
		gzipReader, err := gzip.NewReader(file)
		assert.Nil(err)
		tarReader := tar.NewReader(gzipReader)
		iterations := 0
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				// 已到达存档末尾
				break
			}
			currentFile := files[iterations]
			assert.Equal(currentFile.FileName, header.Name)
			assert.NotEqual(0, header.Size)
			assert.Equal(currentFile.LastModifyTime, time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC))
			iterations++
		}
	})

	t.Run("wrong file path", func(t *testing.T) {
		files := []PackItem{
			{FilePath: path.Join(mydir, "assets/wrong.pdf"),
				FileName:       path.Join("inner/", "《个人防疫手册（第二版）》.pdf"),
				LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
			{FilePath: path.Join(mydir, "assets/导出外贸ERP.ps1"),
				FileName:       path.Join("inner/", "导出外贸ERP.ps1"),
				LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
			{FilePath: path.Join(mydir, "assets/流程排错解决方案.xmind"),
				FileName:       "流程排错解决方案.xmind",
				LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
			{FilePath: path.Join(mydir, "assets/外销合同.txt"),
				FileName:       "inner22/外销合同.txt",
				LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
			{FilePath: path.Join(mydir, "assets/文档.docx"),
				FileName:       "inner22/文档.docx",
				LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
			{FilePath: path.Join(mydir, "assets/文档记录下载1690254313410.zip"),
				FileName:       "inner22/文档记录下载1690254313410.zip",
				LastModifyTime: time.Date(2023, time.September, 1, 12, 0, 0, 0, time.UTC)},
		}
		fileName := path.Join(tempDir, "001.tar.gz")
		request := ExecuteContext{FileName: fileName, IsGziped: true}
		request.Items = append(request.Items, files...)

		err := helper.Pack(request)
		assert.NotNil(err)
		_, ok := err.(*os.PathError)
		assert.True(ok)
	})

}

func TestNewTarHepler(t *testing.T) {
	assert := assert.New(t)
	helper := NewTarHepler()
	assert.NotNil(helper)
}
