package imginfo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"www.github.com/dingsongjie/file-help/pkg/file/img"
	"www.github.com/dingsongjie/file-help/pkg/mocks"
)

func TestGetImgInfo(t *testing.T) {
	t.Run("sample successful", func(t *testing.T) {
		assert := assert.New(t)
		request := GetImgInfoRequest{}
		request.Items = append(request.Items, GetImgInfoRequestItem{FileKey: "1.png"})
		request.Items = append(request.Items, GetImgInfoRequestItem{FileKey: "2.png"})
		queries := NewImgInfoQueries()
		mockedHelper := new(mocks.MockedS3Helper)
		mockedHelper.On("DownLoadAndReturnBuffer", mock.Anything).Return([]byte{'A', 'b'}, nil).Once()
		mockedHelper.On("DownLoadAndReturnBuffer", mock.Anything).Return([]byte{'c'}, nil).Once()
		queries.s3Helper = mockedHelper
		getImgInfo = func(buffer []byte) (*img.ImgInfo, error) {
			if len(buffer) == 2 {
				return &img.ImgInfo{Width: 100, Height: 200}, nil
			} else {
				return &img.ImgInfo{Width: 200, Height: 300}, nil
			}
		}
		response := queries.GetImgInfo(&request)
		assert.True(response.IsAllSucceed)
		assert.Equal(2, len(request.Items))

		assert.True(response.Items[0].IsSucceed)
		assert.Equal(100, response.Items[0].Width)
		assert.Equal(200, response.Items[0].Height)
		assert.Empty(response.Items[0].Message)

		assert.True(response.Items[1].IsSucceed)
		assert.Equal(200, response.Items[1].Width)
		assert.Equal(300, response.Items[1].Height)
		assert.Empty(response.Items[1].Message)
	})
	t.Run("one item faild", func(t *testing.T) {
		assert := assert.New(t)
		request := GetImgInfoRequest{}
		request.Items = append(request.Items, GetImgInfoRequestItem{FileKey: "1.png"})
		request.Items = append(request.Items, GetImgInfoRequestItem{FileKey: "2.png"})
		queries := NewImgInfoQueries()
		mockedHelper := new(mocks.MockedS3Helper)
		mockedHelper.On("DownLoadAndReturnBuffer", mock.Anything).Return([]byte{'A', 'b'}, nil).Once()
		mockedHelper.On("DownLoadAndReturnBuffer", mock.Anything).Return([]byte{'c'}, nil).Once()
		queries.s3Helper = mockedHelper
		getImgInfo = func(buffer []byte) (*img.ImgInfo, error) {
			if len(buffer) == 2 {
				return &img.ImgInfo{Width: 100, Height: 200}, nil
			} else {
				return nil, fmt.Errorf("invalid image format")
			}
		}
		response := queries.GetImgInfo(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(request.Items))

		assert.True(response.Items[0].IsSucceed)
		assert.Equal(100, response.Items[0].Width)
		assert.Equal(200, response.Items[0].Height)
		assert.Empty(response.Items[0].Message)

		assert.False(response.Items[1].IsSucceed)
		assert.Equal(0, response.Items[1].Width)
		assert.Equal(0, response.Items[1].Height)
		assert.Equal("invalid image format", response.Items[1].Message)
	})

	t.Run("download faild", func(t *testing.T) {
		assert := assert.New(t)
		request := GetImgInfoRequest{}
		request.Items = append(request.Items, GetImgInfoRequestItem{FileKey: "1.png"})
		request.Items = append(request.Items, GetImgInfoRequestItem{FileKey: "2.png"})
		queries := NewImgInfoQueries()
		mockedHelper := new(mocks.MockedS3Helper)
		mockedHelper.On("DownLoadAndReturnBuffer", mock.Anything).Return([]byte{'A', 'b'}, nil).Once()
		mockedHelper.On("DownLoadAndReturnBuffer", mock.Anything).Return(nil, fmt.Errorf("key not exist")).Once()
		queries.s3Helper = mockedHelper
		getImgInfo = func(buffer []byte) (*img.ImgInfo, error) {
			if len(buffer) == 2 {
				return &img.ImgInfo{Width: 100, Height: 200}, nil
			} else {
				return &img.ImgInfo{Width: 200, Height: 300}, nil
			}
		}
		response := queries.GetImgInfo(&request)
		assert.False(response.IsAllSucceed)
		assert.Equal(2, len(request.Items))

		assert.True(response.Items[0].IsSucceed)
		assert.Equal(100, response.Items[0].Width)
		assert.Equal(200, response.Items[0].Height)
		assert.Empty(response.Items[0].Message)

		assert.False(response.Items[1].IsSucceed)
		assert.Equal(0, response.Items[1].Width)
		assert.Equal(0, response.Items[1].Height)
		assert.Equal("key not exist", response.Items[1].Message)
	})
}
