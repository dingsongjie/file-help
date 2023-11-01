package img

import (
	"fmt"
	"image"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetimgInfo(t *testing.T) {
	assert := assert.New(t)
	mydir, _ := os.Getwd()
	t.Run("test png seccessful", func(t *testing.T) {
		imgPath := path.Join(mydir, "./assets/1.png")
		file, _ := os.Open(imgPath)
		info, err := GetimgInfo(file)
		assert.Nil(err)
		assert.Equal(264, info.Width)
		assert.Equal(161, info.Height)
	})
	t.Run("test jpg seccessful", func(t *testing.T) {
		imgPath := path.Join(mydir, "./assets/1.jpg")
		file, _ := os.Open(imgPath)
		info, err := GetimgInfo(file)
		assert.Nil(err)
		assert.Equal(2268, info.Width)
		assert.Equal(2409, info.Height)
	})

	t.Run("img decode faild", func(t *testing.T) {
		// imgPath := path.Join(mydir, "./assets/wrong.jpg")
		// //
		// file, _ := os.Open(nil)
		imgDecode = func(r io.Reader) (image.Image, string, error) {
			return nil, "", fmt.Errorf("sample error")
		}
		_, err := GetimgInfo(nil)

		assert.NotNil(err)
		assert.Equal("sample error", err.Error())
	})
}
