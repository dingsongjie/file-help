// package aiconverter

// import (
// 	"sync"

// 	"github.com/mileworks/plm-files-preview/utils"
// 	"www.github.com/dingsongjie/file-help/pkg/converter"
// )

// var (
// 	singletonMu sync.Mutex = sync.Mutex{}
// 	instance    *CadConverter
// )

// type CadConverter struct {
// 	AllowedConverteTypeMaps []*converter.ConverterTypePair
// }

// func NewConverter() *CadConverter {
// 	if instance == nil {
// 		singletonMu.Lock()
// 		if instance == nil {
// 			instance = &CadConverter{}
// 			instance.AllowedConverteTypeMaps = make([]*converter.ConverterTypePair, 3)
// 			instance.AllowedConverteTypeMaps[0] = &converter.ConverterTypePair{SourceType: "dwg", TargetType: "pdf"}
// 			instance.AllowedConverteTypeMaps[1] = &converter.ConverterTypePair{SourceType: "ai", TargetType: "jpeg"}
// 		}
// 		singletonMu.Unlock()
// 	}
// 	return instance
// }

// func (r *CadConverter) ToFastImage(inputFile string, outputFile string) error {
// 	return r.ToFastJpeg(inputFile, outputFile)

// }
// func (r *CadConverter) ToFastJpeg(inputFile string, outputFile string) error {
// 	output := utils.co

// 	return nil
// }

// func (r *CadConverter) ToPrettyPdf(inputFile string, outputFile string) error {
// 	r.mu.Lock()
// 	r.initialise()
// 	args := []string{
// 		"gs", // This will be ignored
// 		"-q",
// 		"-dBATCH",
// 		"-sPageList=1",
// 		"-sDEVICE=pdfwrite",
// 		"-dAutoFilterColorImages=false",
// 		"-dAutoFilterGrayImages=false",
// 		"-dPassThroughMonoImages=false",
// 		"-dPassThroughJPEGImages=false",
// 		"-dPassThroughJPXImages=false",
// 		"-dDownsampleColorImages=false",
// 		"-dDownsampleGrayImages=false",
// 		"-dDownsampleMonoImages=false",
// 		"-dColorImageFilter=/FlateEncode",
// 		"-dNOPAUSE",
// 		"-dNumRenderingThreads=2",
// 		"-dDEVICEWIDTHPOINTS=150",
// 		"-dDEVICEHEIGHTPOINTS=150",
// 		"-dBandBufferSpace=200000000",
// 		"-sBandListStorage=memory",
// 		"-dNoOutputFonts",
// 		"-sOutputFile=" + outputFile,
// 		inputFile,
// 	}

// 	if err := r.internalGSInstance.Init(args); err != nil {
// 		panic(err)
// 	}
// 	r.mu.Unlock()
// 	defer func() {
// 		r.internalGSInstance.Exit()
// 		r.internalGSInstance.Destroy()
// 		r.internalGSInstance = nil
// 	}()
// 	return nil
// }
// func (r *CadConverter) Destory() {

// }
