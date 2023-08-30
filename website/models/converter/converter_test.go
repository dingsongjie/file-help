package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFisrtImageByGavingKeyResponseAddItem(t *testing.T) {
	assert := assert.New(t)
	t.Run("all item succeed", func(t *testing.T) {
		response := NewGetFisrtImageByGavingKeyResponse()
		response.AddItem(&ConvertByGavingKeyResponseItem{SourceKey: "testkey", IsSucceed: true})
		response.AddItem(&ConvertByGavingKeyResponseItem{SourceKey: "testkey2", IsSucceed: true})
		response.AddItem(&ConvertByGavingKeyResponseItem{SourceKey: "testkey3", IsSucceed: true})
		assert.True(response.isAllSucceed)
		assert.Equal(3, len(response.items))
	})

	t.Run("one faild", func(t *testing.T) {
		response := NewGetFisrtImageByGavingKeyResponse()
		response.AddItem(&ConvertByGavingKeyResponseItem{SourceKey: "testkey", IsSucceed: true})
		response.AddItem(&ConvertByGavingKeyResponseItem{SourceKey: "testkey3", IsSucceed: false, Message: "key not exist"})
		response.AddItem(&ConvertByGavingKeyResponseItem{SourceKey: "testkey2", IsSucceed: true})
		assert.False(response.isAllSucceed)
		assert.Equal(3, len(response.items))
		assert.Equal("key not exist", response.items[1].Message)
	})
}
