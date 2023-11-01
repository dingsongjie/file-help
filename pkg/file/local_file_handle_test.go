package file

import (
	"fmt"
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
		file, err := CreateFileAndReturnFile(fullPath)
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

	t.Run("file exist but remove faild", func(t *testing.T) {
		dir := t.TempDir()
		assert := assert.New(t)
		fullPath := path.Join(dir, key)
		file, err := CreateFileAndReturnFile(fullPath)
		osRemove = func(name string) error {
			return os.ErrNotExist
		}
		assert.Nil(err)
		file.Close()
		handler := LocalFileHandle{Path: fullPath, IsDestory: false}
		assert.Nil(handler.Destory())
		assert.True(handler.IsDestory)
	})

	t.Run("remove other faild", func(t *testing.T) {
		dir := t.TempDir()
		assert := assert.New(t)
		fullPath := path.Join(dir, key)
		file, err := CreateFileAndReturnFile(fullPath)
		osRemove = func(name string) error {
			return fmt.Errorf("remove err")
		}
		assert.Nil(err)
		file.Close()
		handler := LocalFileHandle{Path: fullPath, IsDestory: false}
		assert.Error(handler.Destory())
		assert.False(handler.IsDestory)
		osRemove = os.Remove
	})

	t.Run("run safely concurrently ", func(t *testing.T) {
		dir := t.TempDir()
		assert := assert.New(t)
		fullPath := path.Join(dir, key)
		file, err := CreateFileAndReturnFile(fullPath)
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

func TestCreateFileAndReturnFile(t *testing.T) {
	dir := t.TempDir()
	key := "test/1.txt"
	assert := assert.New(t)
	t.Run("success", func(t *testing.T) {
		fullPath := path.Join(dir, key)
		file, err := CreateFileAndReturnFile(fullPath)
		assert.Nil(err)
		_, err = file.Stat()
		assert.Nil(err)
		file.Close()
	})

	t.Run("mkdir faild", func(t *testing.T) {
		fullPath := path.Join(dir, key)
		osMakeDir = func(name string, perm os.FileMode) error {
			return fmt.Errorf("mkdir faild")
		}
		_, err := CreateFileAndReturnFile(fullPath)
		assert.NotNil(err)
		assert.Equal("mkdir faild", err.Error())
		osMakeDir = os.Mkdir
	})

	t.Run("create file faild", func(t *testing.T) {
		fullPath := path.Join(dir, key)
		osCreate = func(name string) (*os.File, error) {
			return nil, fmt.Errorf("create file faild")
		}
		_, err := CreateFileAndReturnFile(fullPath)
		assert.NotNil(err)
		assert.Equal("create file faild", err.Error())
		osCreate = os.Create
	})
}
