// package test

// import (
// 	"fmt"
// 	"os"
// 	"path"
// 	"testing"

// 	"www.github.com/dingsongjie/file-help/pkg/converter/aiconverter"
// 	"www.github.com/dingsongjie/file-help/pkg/log"
// )

// func TestToFastJpeg(t *testing.T) {
// 	log.Initialise()
// 	converter := aiconverter.NewConverter()
// 	mydir, _ := os.Getwd()
// 	t.Run("Ai-1 to jpeg", func(t *testing.T) {
// 		aiAbsolutePath := path.Join(mydir, "./assets/1.ai")
// 		outputAbsolutePath := path.Join(mydir, "./outputs/ai-1.jpeg")
// 		converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
// 	})
// 	t.Run("Ai-2 to jpeg", func(t *testing.T) {
// 		aiAbsolutePath := path.Join(mydir, "./assets/2.ai")
// 		outputAbsolutePath := path.Join(mydir, "./outputs/ai-2.jpeg")
// 		converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
// 	})

// 	t.Run("Ai-3 to jpeg", func(t *testing.T) {
// 		aiAbsolutePath := path.Join(mydir, "./assets/3.ai")
// 		outputAbsolutePath := path.Join(mydir, "./outputs/ai-3.jpeg")
// 		converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
// 	})

// }
// func TestToPrettyPdf(t *testing.T) {
// 	log.Initialise()
// 	converter := aiconverter.NewConverter()
// 	mydir, _ := os.Getwd()
// 	t.Run("Ai-1 to pdf", func(t *testing.T) {
// 		aiAbsolutePath := path.Join(mydir, "./assets/1.ai")
// 		outputAbsolutePath := path.Join(mydir, "./outputs/ai-1.pdf")
// 		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
// 	})
// 	t.Run("Ai-2 to pdf", func(t *testing.T) {
// 		aiAbsolutePath := path.Join(mydir, "./assets/2.ai")
// 		outputAbsolutePath := path.Join(mydir, "./outputs/ai-2.pdf")
// 		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
// 	})

// 	t.Run("Ai-3 to pdf", func(t *testing.T) {
// 		aiAbsolutePath := path.Join(mydir, "./assets/3.ai")
// 		outputAbsolutePath := path.Join(mydir, "./outputs/ai-3.pdf")
// 		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
// 	})
// }
// func BenchmarkToFastJpeg(b *testing.B) {
// 	log.Initialise()
// 	converter := aiconverter.NewConverter()
// 	mydir, _ := os.Getwd()

// 	for i := 0; i < b.N; i++ {
// 		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./assets/%d.ai", i+1))
// 		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./outputs/ai-%d.jpeg", i+1))
// 		converter.ToFastJpeg(aiAbsolutePath, outputAbsolutePath)
// 	}
// }
// func BenchmarkToPrettyPdf(b *testing.B) {
// 	log.Initialise()
// 	converter := aiconverter.NewConverter()
// 	mydir, _ := os.Getwd()

// 	for i := 0; i < b.N; i++ {
// 		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./assets/%d.ai", i+1))
// 		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./outputs/ai-%d.pdf", i+1))
// 		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
// 	}
// }
