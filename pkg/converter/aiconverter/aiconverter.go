package aiconverter

import (
	"sync"

	"github.com/MrSaints/go-ghostscript/ghostscript"

	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/log"
)

var (
	singletonMu sync.Mutex = sync.Mutex{}
	instance    *AiConverter
)

type AiConverter struct {
	internalGSInstance      *ghostscript.Ghostscript
	mu                      sync.Mutex
	AllowedConverteTypeMaps []*converter.ConverterTypePair
}

func NewConverter() *AiConverter {
	if instance == nil {
		singletonMu.Lock()
		if instance == nil {
			instance = &AiConverter{}
			instance.mu = sync.Mutex{}
			instance.AllowedConverteTypeMaps = make([]*converter.ConverterTypePair, 3)
			instance.AllowedConverteTypeMaps[0] = &converter.ConverterTypePair{SourceType: "ai", TargetType: "jpg"}
			instance.AllowedConverteTypeMaps[1] = &converter.ConverterTypePair{SourceType: "ai", TargetType: "jpeg"}
			instance.AllowedConverteTypeMaps[2] = &converter.ConverterTypePair{SourceType: "ai", TargetType: "pdf"}
		}
		singletonMu.Unlock()
	}
	return instance
}

func (r *AiConverter) initialise() {
	if r.internalGSInstance == nil {
		rev, err := ghostscript.GetRevision()
		if err != nil {
			log.Logger.Sugar().Fatalf("Revision: %+v\n", rev)
		}

		r.internalGSInstance, err = ghostscript.NewInstance()
		if err != nil {
			log.Logger.Sugar().Fatalf("Error: %+v\n", err)
		}
	}
}
func (r *AiConverter) ToFastImage(inputFile string, outputFile string) error {
	return r.ToFastJpeg(inputFile, outputFile)

}
func (r *AiConverter) ToFastJpeg(inputFile string, outputFile string) error {
	r.mu.Lock()
	r.initialise()
	args := []string{
		"gs", // This will be ignored
		"-q",
		"-dBATCH",
		"-sPageList=1",
		"-sDEVICE=jpeg",
		"-dAutoFilterColorImages=false",
		"-dAutoFilterGrayImages=false",
		"-dPassThroughMonoImages=false",
		"-dPassThroughJPEGImages=false",
		"-dPassThroughJPXImages=false",
		"-dDownsampleColorImages=false",
		"-dDownsampleGrayImages=false",
		"-dDownsampleMonoImages=false",
		"-dColorImageFilter=/FlateEncode",
		"-dNOPAUSE",
		"-dNumRenderingThreads=2",
		/* 		"-r100", */
		// "-g250x250",
		"-dDEVICEWIDTHPOINTS=150",
		"-dDEVICEHEIGHTPOINTS=150",
		"-dBandBufferSpace=200000000",
		"-sBandListStorage=memory",
		"-sOutputFile=" + outputFile,
		inputFile,
	}

	if err := r.internalGSInstance.Init(args); err != nil {
		panic(err)
	}
	r.mu.Unlock()
	defer func() {
		r.internalGSInstance.Exit()
		r.internalGSInstance.Destroy()
		r.internalGSInstance = nil
	}()
	return nil
}

func (r *AiConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	r.mu.Lock()
	r.initialise()
	args := []string{
		"gs", // This will be ignored
		"-q",
		"-dBATCH",
		"-sPageList=1",
		"-sDEVICE=pdfwrite",
		"-dAutoFilterColorImages=false",
		"-dAutoFilterGrayImages=false",
		"-dPassThroughMonoImages=false",
		"-dPassThroughJPEGImages=false",
		"-dPassThroughJPXImages=false",
		"-dDownsampleColorImages=false",
		"-dDownsampleGrayImages=false",
		"-dDownsampleMonoImages=false",
		"-dColorImageFilter=/FlateEncode",
		"-dNOPAUSE",
		"-dNumRenderingThreads=2",
		"-dDEVICEWIDTHPOINTS=150",
		"-dDEVICEHEIGHTPOINTS=150",
		"-dBandBufferSpace=200000000",
		"-sBandListStorage=memory",
		"-dNoOutputFonts",
		"-sOutputFile=" + outputFile,
		inputFile,
	}

	if err := r.internalGSInstance.Init(args); err != nil {
		panic(err)
	}
	r.mu.Unlock()
	defer func() {
		r.internalGSInstance.Exit()
		r.internalGSInstance.Destroy()
		r.internalGSInstance = nil
	}()
	return nil
}
func (r *AiConverter) Destory() {

}
