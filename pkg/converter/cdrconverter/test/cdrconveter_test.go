package test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"www.github.com/dingsongjie/file-help/pkg/converter/cdrconverter"
	"www.github.com/dingsongjie/file-help/pkg/log"
)

func TestToFastPng(t *testing.T) {
	log.Initialise()
	converter := cdrconverter.NewConverter()
	mydir, _ := os.Getwd()
	errorFormat := "convert error : %s"
	t.Run("cdr-1 to png", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/1.cdr")
		outputAbsolutePath := path.Join(mydir, "./outputs/cdr-1.png")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		if err != nil {
			t.Errorf(errorFormat, err)
		}
	})
	t.Run("cdr-2 to png", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/2.cdr")
		outputAbsolutePath := path.Join(mydir, "./outputs/cdr-2.png")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		if err != nil {
			t.Errorf(errorFormat, err)
		}
	})

	t.Run("cdr-3 to png", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/3.cdr")
		outputAbsolutePath := path.Join(mydir, "./outputs/cdr-3.png")
		err := converter.ToFastImage(aiAbsolutePath, outputAbsolutePath)
		if err != nil {
			t.Errorf(errorFormat, err)
		}
	})

	converter.Destory()
}
func TestToPrettyPdf(t *testing.T) {
	log.Initialise()
	converter := cdrconverter.NewConverter()
	mydir, _ := os.Getwd()
	t.Run("cdr-1 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/1.cdr")
		outputAbsolutePath := path.Join(mydir, "./outputs/cdr-1.pdf")
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
	})
	t.Run("cdr-2 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/2.cdr")
		outputAbsolutePath := path.Join(mydir, "./outputs/cdr-2.pdf")
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
	})

	t.Run("cdr-3 to pdf", func(t *testing.T) {
		aiAbsolutePath := path.Join(mydir, "./assets/3.cdr")
		outputAbsolutePath := path.Join(mydir, "./outputs/cdr-3.pdf")
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
	})
	converter.Destory()
}
func BenchmarkToFastPng(b *testing.B) {
	log.Initialise()
	converter := cdrconverter.NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./assets/%d.cdr", i+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./outputs/cdr-%d.png", i+1))
		converter.ToFastPng(aiAbsolutePath, outputAbsolutePath)
	}
	converter.Destory()
}
func BenchmarkToPrettyPdf(b *testing.B) {
	log.Initialise()
	converter := cdrconverter.NewConverter()
	mydir, _ := os.Getwd()

	for i := 0; i < b.N; i++ {
		aiAbsolutePath := path.Join(mydir, fmt.Sprintf("./assets/%d.cdr", i+1))
		outputAbsolutePath := path.Join(mydir, fmt.Sprintf("./outputs/cdr-%d.pdf", i+1))
		converter.ToPrettyPdf(aiAbsolutePath, outputAbsolutePath)
	}
	converter.Destory()
}
