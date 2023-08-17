package test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"www.github.com/dingsongjie/file-help/pkg/converter/imagickconverter"
	"www.github.com/dingsongjie/file-help/pkg/log"
)

func TestConvertToPreview(t *testing.T) {
	log.Initialise()
	converter := imagickconverter.NewConverter()
	mydir, _ := os.Getwd()
	t.Run("convert all pages of psd-1 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/1.psd")
		outputAbsolutePath := path.Join(mydir, "./outputs/convertall/psd-1.jpeg")
		converter.ConvertToJpeg(aiAbsolutePath, outputAbsolutePath, false)
	})
	t.Run("convert all pages of psd-2 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/2.psd")
		outputAbsolutePath := path.Join(mydir, "./outputs/convertall/psd-2.jpeg")
		converter.ConvertToJpeg(aiAbsolutePath, outputAbsolutePath, false)
	})

	t.Run("convert all pages of psd-3 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/3.psd")
		outputAbsolutePath := path.Join(mydir, "./outputs/convertall/psd-3.jpeg")
		converter.ConvertToJpeg(aiAbsolutePath, outputAbsolutePath, false)
	})

	t.Run("convert first page of psd-1 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/1.psd")
		outputAbsolutePath := path.Join(mydir, "./outputs/convertfirst/psd-1.jpeg")
		converter.ConvertToJpeg(aiAbsolutePath, outputAbsolutePath, true)
	})
	t.Run("convert first page of psd-2 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/2.psd")
		outputAbsolutePath := path.Join(mydir, "./outputs/convertfirst/psd-2.jpeg")
		converter.ConvertToJpeg(aiAbsolutePath, outputAbsolutePath, true)
	})

	t.Run("convert first page of psd-3 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/3.psd")
		outputAbsolutePath := path.Join(mydir, "./outputs/convertfirst/psd-3.jpeg")
		converter.ConvertToJpeg(aiAbsolutePath, outputAbsolutePath, true)
	})
	converter.Destory()
}

func BenchmarkConvertToPreview(b *testing.B) {
	log.Initialise()
	converter := imagickconverter.NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./assets/%d.psd", i%3+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./outputs/psd-%d.jpeg", i%3+1))
		converter.ConvertToJpeg(aiAbsolutePath, outputAbsolutePath, true)
	}
	converter.Destory()
}
