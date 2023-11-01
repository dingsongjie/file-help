package file

import (
	"fmt"
	"io"
	"net/http"
	urlUtil "net/url"
	"os"
	"path"
	"path/filepath"
)

var (
	httpFileDownloadDir     = os.TempDir()
	httpGet                 = http.Get
	createFileAndReturnFile = CreateFileAndReturnFile
)

func DownLoadAndReturnLocalPath(url string) (*LocalFileHandle, error) {
	response, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 检查HTTP响应状态
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", response.Status)
	}
	// 从Content-Disposition标头中获取文件名
	filename := response.Header.Get("Content-Disposition")
	if filename == "" {
		u, _ := urlUtil.Parse(url)

		// 获取URL路径中的文件名
		filename = filepath.Base(u.Path)
	}
	if filename == "" || filename == "." {
		return nil, fmt.Errorf("Content-Disposition not found in http response and also we can not get filename from request url,url: %s", url)
	}
	filePath := path.Join(httpFileDownloadDir, filename)
	f, err := createFileAndReturnFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = io.Copy(f, response.Body)
	if err != nil {
		return nil, err
	}
	handle := LocalFileHandle{Path: filePath, IsDestory: false}
	return &handle, nil
}
