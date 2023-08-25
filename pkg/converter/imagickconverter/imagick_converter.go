package imagickconverter

import (
	"fmt"
	"sync"

	"gopkg.in/gographics/imagick.v3/imagick"
	"www.github.com/dingsongjie/file-help/pkg/converter"
)

var (
	instance    *ImagickConverter
	singletonMu sync.Mutex = sync.Mutex{}
)

type ImagickConverter struct {
	internalGSInstance      *imagick.MagickWand
	mu                      sync.Mutex
	AllowedConverteTypeMaps []*converter.ConverterTypePair
}

func NewConverter() *ImagickConverter {
	if instance == nil {
		singletonMu.Lock()
		if instance == nil {
			instance = &ImagickConverter{}
			instance.mu = sync.Mutex{}
			instance.AllowedConverteTypeMaps = make([]*converter.ConverterTypePair, 3)
			instance.AllowedConverteTypeMaps[0] = &converter.ConverterTypePair{SourceType: "psd", TargetType: "jpeg"}
			instance.AllowedConverteTypeMaps[1] = &converter.ConverterTypePair{SourceType: "psd", TargetType: "jpg"}
			instance.AllowedConverteTypeMaps[2] = &converter.ConverterTypePair{SourceType: "psd", TargetType: "pdf"}
			imagick.Initialize()
		}
		singletonMu.Unlock()
	}
	return instance
}

func (r *ImagickConverter) convert(inputFile string, outputFile string, firstPage bool) error {
	if firstPage {
		inputFile = inputFile + "[0]"
	}
	result, err := imagick.ConvertImageCommand([]string{"convert", inputFile, outputFile})
	fmt.Print(result)
	if err != nil {
		return err
	}
	return nil
}

func (r *ImagickConverter) ToFastImage(inputFile string, outputFile string) error {
	return r.ConvertToJpeg(inputFile, outputFile, true)
}

func (r *ImagickConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	return r.convert(inputFile, outputFile, false)
}

func (r *ImagickConverter) ConvertToJpeg(inputFile string, outputFile string, firstPage bool) error {
	return r.convert(inputFile, outputFile, firstPage)
}

func (r *ImagickConverter) Destory() {
	if r.internalGSInstance != nil {
		r.mu.Lock()
		if r.internalGSInstance != nil {
			r.internalGSInstance.Destroy()
		}
		r.mu.Unlock()
	}
	imagick.Terminate()
}
