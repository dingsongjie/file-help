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
	ioCopy                  = io.Copy
)

func DownLoadAndReturnLocalPath(url, fileName string) (*LocalFileHandle, error) {
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
	if fileName == "" {
		fileName = response.Header.Get("Content-Disposition")
		if fileName == "" {
			u, _ := urlUtil.Parse(url)

			// 获取URL路径中的文件名
			fileName = filepath.Base(u.Path)
		}
		if fileName == "" || fileName == "." {
			return nil, fmt.Errorf("Content-Disposition not found in http response and also we can not get filename from request url,url: %s", url)
		}
	}

	filePath := path.Join(httpFileDownloadDir, fileName)
	f, err := createFileAndReturnFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = ioCopy(f, response.Body)
	if err != nil {
		return nil, err
	}
	handle := LocalFileHandle{Path: filePath, IsDestory: false}
	return &handle, nil
}
