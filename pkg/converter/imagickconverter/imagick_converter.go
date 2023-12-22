package imagickconverter

import (
	"fmt"
	"sync"

	"github.com/STRockefeller/go-linq"
	"github.com/dingsongjie/file-help/pkg/converter"
	"github.com/dingsongjie/file-help/pkg/log"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	instance    *ImagickConverter
	singletonMu sync.Mutex = sync.Mutex{}
)

type ImagickConverter struct {
	// internalGSInstance      *imagick.MagickWand
	mu                      sync.Mutex
	AllowedConverteTypeMaps linq.Linq[*converter.ConverterTypePair]
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

func (r *ImagickConverter) convertToPdf(inputFile string, outputFile string) error {
	// if firstPage {
	// 	inputFile = inputFile + "[0]"
	// }
	_, err := imagick.ConvertImageCommand([]string{"convert", inputFile, outputFile})
	if err != nil {
		log.Logger.Error(err.Error())
		return fmt.Errorf("convert faild")
	}
	return nil
}

func (r *ImagickConverter) ToFastImage(inputFile string, outputFile string, dpi int) error {
	if dpi == 0 {
		dpi = 72
	}
	if dpi > 300 {
		return fmt.Errorf("dpi is not allowed to exceed 300")
	}
	return r.convertToJpeg(inputFile, outputFile, true, dpi)
}

func (r *ImagickConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	return r.convertToPdf(inputFile, outputFile)
}

func (r *ImagickConverter) convertToJpeg(inputFile string, outputFile string, firstPage bool, dpi int) error {
	if firstPage {
		inputFile = inputFile + "[0]"
	}
	_, err := imagick.ConvertImageCommand([]string{"convert", "-density", fmt.Sprint(dpi), "-units", "pixelsperinch", inputFile, outputFile})
	//_, err := imagick.ConvertImageCommand([]string{"convert", inputFile, outputFile})
	if err != nil {
		log.Logger.Error(err.Error())
		return fmt.Errorf("convert faild")
	}
	return nil
}

func (r *ImagickConverter) CanHandle(pair converter.ConverterTypePair) bool {
	return r.AllowedConverteTypeMaps.Exists(func(ctp *converter.ConverterTypePair) bool {

		return ctp.SourceType == pair.SourceType && ctp.TargetType == pair.TargetType
	})
}

func (r *ImagickConverter) Destory() {
	imagick.Terminate()
}
