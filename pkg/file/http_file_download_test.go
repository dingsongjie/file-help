package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownLoadAndReturnLocalPath(t *testing.T) {
	httpFileDownloadDir = t.TempDir()
	assert := assert.New(t)
	filePath := path.Join(httpFileDownloadDir, "1.txt")
	f, _ := CreateFileAndReturnFile(filePath)
	defer f.Close()
	t.Run("download and return successful", func(t *testing.T) {

		httpGet = func(url string) (resp *http.Response, err error) {
			//Status     string // e.g. "200 OK"
			//StatusCode int    // e.g. 200
			//Proto      string // e.g. "HTTP/1.0"
			//ProtoMajor int    // e.g. 1
			//ProtoMinor int    // e.g. 0
			response := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
			}
			response.Header = http.Header{}
			response.Header.Add("Content-Disposition", "1.txt")
			f.Write([]byte{'A', 'B'})
			// f.Close()
			response.Body = f
			return &response, nil
		}
		localFile, err := DownLoadAndReturnLocalPath("test-url", "1.txt")
		assert.NotNil(localFile)
		assert.Equal(filePath, localFile.Path)
		assert.False(localFile.IsDestory)
		assert.Nil(err)
		localFile.Destory()
	})

	t.Run("http get faild", func(t *testing.T) {
		httpGet = func(url string) (resp *http.Response, err error) {
			return nil, fmt.Errorf("url not vaild")
		}
		localFile, err := DownLoadAndReturnLocalPath("test-url", "")
		assert.Nil(localFile)
		assert.NotNil(err)
		assert.Equal("url not vaild", err.Error())
	})
	t.Run("http response code not 200", func(t *testing.T) {
		httpGet = func(url string) (resp *http.Response, err error) {
			//Status     string // e.g. "200 OK"
			//StatusCode int    // e.g. 200
			//Proto      string // e.g. "HTTP/1.0"
			//ProtoMajor int    // e.g. 1
			//ProtoMinor int    // e.g. 0
			response := http.Response{
				Status:     "400 NOTFOUND",
				StatusCode: 400,
				Proto:      "HTTP/1.1",
			}
			response.Body = &os.File{}
			return &response, nil
		}
		localFile, err := DownLoadAndReturnLocalPath("test-url", "")
		assert.Nil(localFile)
		assert.NotNil(err)
		assert.Equal("HTTP request failed with status: 400 NOTFOUND", err.Error())
	})

	t.Run("Content-Disposition not in response header", func(t *testing.T) {
		httpGet = func(url string) (resp *http.Response, err error) {
			//Status     string // e.g. "200 OK"
			//StatusCode int    // e.g. 200
			//Proto      string // e.g. "HTTP/1.0"
			//ProtoMajor int    // e.g. 1
			//ProtoMinor int    // e.g. 0
			response := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
			}
			response.Header = http.Header{}
			f, _ := os.Open(filePath)
			f.Write([]byte{'A', 'B'})
			f.Close()
			response.Body = f
			return &response, nil
		}
		localFile, err := DownLoadAndReturnLocalPath("http://www.test-url.com", "")
		assert.Nil(localFile)
		assert.NotNil(err)
		assert.Equal("Content-Disposition not found in http response and also we can not get filename from request url,url: http://www.test-url.com", err.Error())
	})

	t.Run("create local file faild", func(t *testing.T) {
		httpGet = func(url string) (resp *http.Response, err error) {
			//Status     string // e.g. "200 OK"
			//StatusCode int    // e.g. 200
			//Proto      string // e.g. "HTTP/1.0"
			//ProtoMajor int    // e.g. 1
			//ProtoMinor int    // e.g. 0
			response := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
			}
			response.Header = http.Header{}
			response.Header.Add("Content-Disposition", "1.txt")
			f, _ := os.Open(filePath)
			f.Write([]byte{'A', 'B'})
			f.Close()
			response.Body = f
			return &response, nil
		}
		createFileAndReturnFile = func(fullPath string) (*os.File, error) {
			return nil, fmt.Errorf("create local file faild")
		}
		localFile, err := DownLoadAndReturnLocalPath("test-url", "")
		assert.Nil(localFile)
		assert.NotNil(err)
		assert.Equal("create local file faild", err.Error())
		createFileAndReturnFile = CreateFileAndReturnFile
	})

	t.Run("io copy faild", func(t *testing.T) {
		httpGet = func(url string) (resp *http.Response, err error) {
			//Status     string // e.g. "200 OK"
			//StatusCode int    // e.g. 200
			//Proto      string // e.g. "HTTP/1.0"
			//ProtoMajor int    // e.g. 1
			//ProtoMinor int    // e.g. 0
			response := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
			}
			response.Header = http.Header{}
			response.Header.Add("Content-Disposition", "1.txt")
			f.Write([]byte{'A', 'B'})
			// f.Close()
			response.Body = f
			return &response, nil
		}
		ioCopy = func(dst io.Writer, src io.Reader) (written int64, err error) {
			return int64(0), fmt.Errorf("copy faild")
		}
		localFile, err := DownLoadAndReturnLocalPath("test-url", "1.txt")
		assert.Nil(localFile)
		assert.NotNil(err)
		assert.Equal("copy faild", err.Error())
		ioCopy = io.Copy
	})
}
