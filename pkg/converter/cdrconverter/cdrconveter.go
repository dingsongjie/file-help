package cdrconverter

import (
	"fmt"
	"sync"

	"github.com/dingsongjie/file-help/pkg/converter"
	"github.com/dingsongjie/file-help/pkg/log"
	"github.com/galihrivanto/go-inkscape"
)

var (
	instance    *CdrConverter
	singletonMu sync.Mutex = sync.Mutex{}
)

type CdrConverter struct {
	internalGSInstance      *inkscape.Proxy
	mu                      sync.Mutex
	AllowedConverteTypeMaps []*converter.ConverterTypePair
}

func NewConverter() *CdrConverter {
	if instance == nil {
		singletonMu.Lock()
		if instance == nil {
			instance = &CdrConverter{}
			instance.mu = sync.Mutex{}
			instance.AllowedConverteTypeMaps = make([]*converter.ConverterTypePair, 2)
			instance.AllowedConverteTypeMaps[0] = &converter.ConverterTypePair{SourceType: "cdr", TargetType: "png"}
			instance.AllowedConverteTypeMaps[1] = &converter.ConverterTypePair{SourceType: "cdr", TargetType: "pdf"}
			instance.initialise()
		}
		singletonMu.Unlock()
	}
	return instance
}
func (r *CdrConverter) initialise() {
	if r.internalGSInstance == nil {
		r.mu.Lock()
		if r.internalGSInstance == nil {
			r.internalGSInstance = inkscape.NewProxy(inkscape.Verbose(false))
			err := r.internalGSInstance.Run()
			if err != nil {
				log.Logger.Sugar().Fatalf("new proxy faild")
			}
		}
		r.mu.Unlock()
	}
}

func (r *CdrConverter) ToFastImage(inputFile string, outputFile string) error {
	return r.ToFastPng(inputFile, outputFile)
}

func (r *CdrConverter) ToFastPng(inputFile string, outputFile string) error {
	_, err := r.internalGSInstance.RawCommands(
		fmt.Sprintf("export-filename:%s", outputFile),
		fmt.Sprintf("file-open:%s", inputFile),
		"export-do",
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *CdrConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	return r.ToFastPng(inputFile, outputFile)
}
func (r *CdrConverter) Destory() {
	if r.internalGSInstance != nil {
		r.mu.Lock()
		if r.internalGSInstance != nil {
			r.internalGSInstance.Close()
		}
		r.mu.Unlock()
	}
}
