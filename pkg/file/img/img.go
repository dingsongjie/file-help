package img

import (
	"bytes"
	"image"
	_ "image/gif"  // 支持GIF格式
	_ "image/jpeg" // 支持JPEG格式
	_ "image/png"  // 支持PNG格式
)

var imgDecode = image.Decode

type ImgInfo struct {
	Width, Height int
}

func GetImgInfo(buffer []byte) (*ImgInfo, error) {
	img, _, err := imgDecode(bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}

	// 获取图像的宽度和高度
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	return &ImgInfo{Width: width, Height: height}, nil
}
