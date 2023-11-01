package imginfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGetImgInfoResponse(t *testing.T) {
	assert := assert.New(t)
	response := NewGetImgInfoResponse()
	assert.True(response.IsAllSucceed)
}

func TestAddItem(t *testing.T) {
	assert := assert.New(t)

	t.Run("sample success", func(t *testing.T) {
		response := NewGetImgInfoResponse()
		response.AddItem(&GetImgInfoItemResponse{FileKey: "test.json", Width: 40, Height: 50, IsSucceed: true})
		assert.True(response.IsAllSucceed)
	})
	t.Run("sample faild", func(t *testing.T) {
		response := NewGetImgInfoResponse()
		response.AddItem(&GetImgInfoItemResponse{FileKey: "test.json", Width: 40, Height: 50, IsSucceed: true})
		response.AddItem(&GetImgInfoItemResponse{FileKey: "test2.json", Width: 40, Height: 50, IsSucceed: false})
		assert.False(response.IsAllSucceed)
	})
}
