package aiconverter

import (
	"fmt"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"www.github.com/dingsongjie/file-help/pkg/log"
)

var assertAI3 string = "./test/assets/3.ai"
var assertAI2Jpg3 string = "./test/outputs/ai-3.jpeg"

func TestToFastJpeg(t *testing.T) {
	log.Initialise()
	assert := assert.New(t)
	converter := NewConverter()
	mydir, _ := os.Getwd()

	t.Run("Ai-1 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/1.ai")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-1.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})
	t.Run("Ai-2 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./test/assets/2.ai")
		outputAbsolutePath := path.Join(mydir, "./test/outputs/ai-2.jpeg")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
	})

	t.Run("Ai-3 to jpeg", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, assertAI3)
		outputAbsolutePath := path.Join(mydir, assertAI2Jpg3)
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		assert.Nil(err)
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
				err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
				assert.Nil(err)
				converter.Destory()
				wg.Done()
			}()
		}
		wg.Wait()
	})

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
}
func BenchmarkToFastJpeg(b *testing.B) {
	log.Initialise()
	converter := NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/assets/%d.ai", i+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./test/outputs/ai-%d.jpeg", i+1))
		converter.ToFastJpeg(aiAbsolutePath, outputAbsolutePath)
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
