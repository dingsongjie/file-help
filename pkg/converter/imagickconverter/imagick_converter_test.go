package imagickconverter

import (
	"fmt"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/log"
)

var assert1 string = "./test/assets/1.psd"

func TestCanHandle(t *testing.T) {
	log.Initialise()
	newconverter := NewConverter()
	assert := assert.New(t)

	pdfToJpg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "pdf", TargetType: "jpg"})
	assert.True(pdfToJpg)
	pdfToJpeg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "pdf", TargetType: "jpeg"})
	assert.True(pdfToJpeg)
	psdToJpg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "psd", TargetType: "jpg"})
	assert.True(psdToJpg)
	psdToJpeg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "psd", TargetType: "jpeg"})
	assert.True(psdToJpeg)
	psdTopdf := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "psd", TargetType: "pdf"})
	assert.True(psdTopdf)
}

func TestToFastImage(t *testing.T) {
	log.Initialise()
	converter := NewConverter()
	mydir, _ := os.Getwd()
	assert := assert.New(t)

	// t.Run("convert all pages of psd-1 to jpeg", func(t *testing.T) {
	// 	aiAbsolutePath := path.Join(mydir, "./test/assets/1.psd")
	// 	outputAbsolutePath := path.Join(mydir, "./test/outputs/convertall/psd-1.jpeg")
	// 	err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
	// 	assert.Nil(err)
	// })
	// t.Run("convert all pages of psd-2 to jpeg", func(t *testing.T) {
	// 	aiAbsolutePath := path.Join(mydir, "./test/assets/2.psd")
	// 	outputAbsolutePath := path.Join(mydir, "./test/outputs/convertall/psd-2.jpeg")
	// 	err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
	// 	assert.Nil(err)
	// })

	// t.Run("convert all pages of psd-3 to jpeg", func(t *testing.T) {
	// 	aiAbsolutePath := path.Join(mydir, "./test/assets/3.psd")
	// 	outputAbsolutePath := path.Join(mydir, "./test/outputs/convertall/psd-3.jpeg")
	// 	err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
	// 	assert.Nil(err)
	// })

	t.Run("convert first page of psd-1 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, assert1)
		outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-1.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})

	t.Run("convert first page of pdf to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/1.pdf")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/pdf-1.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})
	t.Run("convert first page of psd-2 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/2.psd")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-2.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})

	t.Run("worng path of psd to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/4.psd")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-4.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		assert.NotNil(err)
	})
	converter.Destory()

	t.Run("run safely concurrently", func(t *testing.T) {
		wantedCount := 10
		var wg sync.WaitGroup
		wg.Add(wantedCount)
		for i := 0; i < wantedCount; i++ {
			go func() {
				converter := NewConverter()
				aiAbsolutePath := path.Join(mydir, assert1)
				outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-1.jpeg")
				err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
				assert.Nil(err)
				converter.Destory()
				wg.Done()
			}()
		}
		wg.Wait()

	})
	converter.Destory()
}

func BenchmarkToFastImage(b *testing.B) {
	log.Initialise()
	converter := NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/assets/%d.psd", i%3+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/outputs/psd-%d.jpeg", i%3+1))
		converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
	}
	converter.Destory()
}

func TestToPrettyPdf(t *testing.T) {
	log.Initialise()
	converter := NewConverter()
	mydir, _ := os.Getwd()
	assert := assert.New(t)

	t.Run("convert first page of psd-1 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/1.psd")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-1.pdf")
		err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})
	t.Run("convert first page of psd-2 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/2.psd")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-2.pdf")
		err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})

	t.Run("worng path of psd to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/4.psd")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-4.pdf")
		err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.NotNil(err)
	})
	converter.Destory()

	t.Run("run safely concurrently", func(t *testing.T) {
		wantedCount := 10
		var wg sync.WaitGroup
		wg.Add(wantedCount)
		for i := 0; i < wantedCount; i++ {
			go func() {
				converter := NewConverter()
				aiAbsolutePath := path.Join(mydir, assert1)
				outputAbsolutePath := path.Join(mydir, "./test/outputs/convertfirst/psd-1.pdf")
				err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
				assert.Nil(err)
				converter.Destory()
				wg.Done()
			}()
		}
		wg.Wait()

	})
	converter.Destory()
}

func BenchmarkToPrettyPdf(b *testing.B) {
	log.Initialise()
	converter := NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/assets/%d.psd", i%2+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/outputs/psd-%d.jpeg", i%2+1))
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
	}
	converter.Destory()
}
