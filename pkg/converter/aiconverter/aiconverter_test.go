package aiconverter

import (
	"fmt"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/MrSaints/go-ghostscript/ghostscript"
	"github.com/stretchr/testify/assert"
	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/log"
)

var assertAI3 string = "./test/assets/3.ai"
var assertAI2Jpg3 string = "./test/outputs/ai-3.jpeg"
var defaultDpi int = 72

func TestCanHandle(t *testing.T) {
	log.Initialise()
	newconverter := NewConverter()
	assert := assert.New(t)

	pdfToJpg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "pdf", TargetType: "jpg"})
	assert.True(pdfToJpg)
	pdfToJpeg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "pdf", TargetType: "jpeg"})
	assert.True(pdfToJpeg)
	aiToJpg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "ai", TargetType: "jpg"})
	assert.True(aiToJpg)
	aiToJpeg := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "ai", TargetType: "jpeg"})
	assert.True(aiToJpeg)
	aiTopdf := newconverter.CanHandle(converter.ConverterTypePair{SourceType: "ai", TargetType: "pdf"})
	assert.True(aiTopdf)
}

func TestToFastJpeg(t *testing.T) {
	log.Initialise()
	assert := assert.New(t)
	converter := NewConverter()
	mydir, _ := os.Getwd()

	t.Run("Ai-1 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/1.ai")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-1.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, 0)
		assert.Nil(err)
	})
	t.Run("Ai-2 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/2.ai")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-2.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.Nil(err)
	})

	t.Run("Ai-3 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, assertAI3)
		outputAbsolutePath := path.Join(mydir, assertAI2Jpg3)
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.Nil(err)
	})

	t.Run("convert first page of 1.pdf to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/1.pdf")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/pdf-1.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.Nil(err)
	})

	t.Run("convert first page of 2.pdf to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/2.pdf")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/pdf-2.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.Nil(err)
	})

	t.Run("convert first page of 3.pdf to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/3.pdf")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/pdf-3.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.Nil(err)
	})

	t.Run("gs is released when init faild", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/3.pdf1")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/pdf-3.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.NotNil(err)

		aiAbsolutePath = path.Join(mydir, "./test/assets/3.pdf")
		outputAbsolutePath = path.Join(mydir, "./test/outputs/pdf-3.jpeg")
		err = converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.Nil(err)
	})

	t.Run("not supported dpi", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/1.ai")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-1.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, 310)
		assert.NotNil(err)
	})
	// 如果 gs init 出错进程会崩溃，所以忽略文件找不到的情况
	// t.Run("wrong path ai", func(t *testing.T) {
	// 	aiAbsolutePath := path.Join(mydir, "./test/assets/4.ai")
	// 	outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-4.jpeg")
	// 	err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
	// 	assert.NotNil(err)
	// })
	converter.Destory()

	t.Run("run safely concurrently", func(t *testing.T) {
		wantedCount := 10
		var wg sync.WaitGroup
		wg.Add(wantedCount)
		for i := 0; i < wantedCount; i++ {
			go func() {
				converter := NewConverter()
				aiAbsolutePath := path.Join(mydir, assertAI3)
				outputAbsolutePath := path.Join(mydir, assertAI2Jpg3)
				err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
				assert.Nil(err)
				converter.Destory()
				wg.Done()
			}()
		}
		wg.Wait()
	})

	t.Run("get revision faild", func(t *testing.T) {
		converter := NewConverter()
		aiAbsolutePath := path.Join(mydir, "./test/assets/3.pdf")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/pdf-3.jpeg")
		gsGetRevision = func() (ghostscript.Revision, error) {
			return ghostscript.Revision{}, fmt.Errorf("get revision faild")
		}
		logGetRevisionFaild = func(err error, rev ghostscript.Revision) {
			assert.NotNil(err)
			assert.Equal("get revision faild", err.Error())
		}
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
		assert.Nil(err)
		gsGetRevision = ghostscript.GetRevision
	})

	// t.Run("new instance faild", func(t *testing.T) {
	// 	converter := NewConverter()
	// 	aiAbsolutePath := path.Join(mydir, "./test/assets/3.pdf")
	// 	outputAbsolutePath := path.Join(mydir, "./test/outputs/pdf-3.jpeg")
	// 	gsNewInstance = func() (*ghostscript.Ghostscript, error) {
	// 		return nil, fmt.Errorf("new instance faild")
	// 	}
	// 	logNewInstanceFaild = func(err error) {
	// 		assert.NotNil(err)
	// 		assert.Equal("new instance faild", err.Error())
	// 	}
	// 	err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath, defaultDpi)
	// 	assert.Nil(err)
	// 	gsNewInstance = ghostscript.NewInstance
	// })
}
func TestToPrettyPdf(t *testing.T) {
	log.Initialise()
	assert := assert.New(t)
	converter := NewConverter()
	mydir, _ := os.Getwd()
	t.Run("Ai-1 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/1.ai")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-1.pdf")
		err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})
	t.Run("Ai-2 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/2.ai")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-2.pdf")
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})

	t.Run("Ai-3 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, assertAI3)
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-3.pdf")
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})

	// 如果 gs init 出错进程会崩溃，所以忽略文件找不到的情况

	converter.Destory()

	t.Run("run safely concurrently", func(t *testing.T) {
		wantedCount := 10
		var wg sync.WaitGroup
		wg.Add(wantedCount)
		for i := 0; i < wantedCount; i++ {
			go func() {
				converter := NewConverter()
				aiAbsolutePath := path.Join(mydir, assertAI3)
				outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-3.pdf")
				err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
				assert.Nil(err)
				converter.Destory()
				wg.Done()
			}()
		}
		wg.Wait()
	})

	t.Run("gs is released when init faild", func(t *testing.T) {
		converter := NewConverter()
		aiAbsolutePath := path.Join(mydir, "./test/assets/assets/2.aif")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-2.pdf")
		err := converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.NotNil(err)

		aiAbsolutePath = path.Join(mydir, "./test/assets/2.ai")
		outputAbsolutePath = path.Join(mydir, "./test/outputs/ai-2.pdf")
		err = converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})
}
func BenchmarkToFastJpeg(b *testing.B) {
	log.Initialise()
	converter := NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/assets/%d.ai", i+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/outputs/ai-%d.jpeg", i+1))
		converter.ToFastJpeg(aiAbsolutePath, outputAbsolutePath, defaultDpi)
	}
	converter.Destory()
}
func BenchmarkToPrettyPdf(b *testing.B) {
	log.Initialise()
	converter := NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/assets/%d.ai", i+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/outputs/ai-%d.pdf", i+1))
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
	}
	converter.Destory()
}
