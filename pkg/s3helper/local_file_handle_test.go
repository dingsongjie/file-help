package s3helper

import (
	"os"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDestory(t *testing.T) {
	key := "test/1.txt"
	t.Run("file exist", func(t *testing.T) {
		dir := t.TempDir()
		assert := assert.New(t)
		fullPath := path.Join(dir, key)
		file, err := createFileAndReturnFile(fullPath)
		assert.Nil(err)
		file.Close()
		handler := LocalFileHandle{Path: fullPath, IsDestory: false}
		handler.Destory()

		assert.True(handler.IsDestory)
		_, err = os.Stat(fullPath)
		assert.ErrorIs(err, os.ErrNotExist)
	})
	t.Run("file not exist", func(t *testing.T) {
		dir := t.TempDir()
		assert := assert.New(t)
		fullPath := path.Join(dir, key)
		handler := LocalFileHandle{Path: fullPath, IsDestory: false}
		handler.Destory()

		assert.True(handler.IsDestory)
	})

	t.Run("run safely concurrently ", func(t *testing.T) {
		dir := t.TempDir()
		assert := assert.New(t)
		fullPath := path.Join(dir, key)
		file, err := createFileAndReturnFile(fullPath)
		file.Close()
		assert.Nil(err)
		handler := LocalFileHandle{Path: fullPath, IsDestory: false}
		wantedCount := 100

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func() {
				err := handler.Destory()
				assert.Nil(err)
				wg.Done()
			}()
		}
		wg.Wait()
		assert.True(handler.IsDestory)
	})
}
