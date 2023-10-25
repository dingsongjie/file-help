package aiconverter

import (
	"sync"

	"github.com/MrSaints/go-ghostscript/ghostscript"
	"github.com/STRockefeller/go-linq"

	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/log"
)

var (
	singletonMu sync.Mutex = sync.Mutex{}
	gsCommandMu sync.Mutex = sync.Mutex{}
	instance    *AiConverter
)

type AiConverter struct {
	internalGSInstance      *ghostscript.Ghostscript
	AllowedConverteTypeMaps linq.Linq[*converter.ConverterTypePair]
}

func NewConverter() *AiConverter {
	if instance == nil {
		singletonMu.Lock()
		if instance == nil {
			instance = &AiConverter{}
			instance.AllowedConverteTypeMaps = make([]*converter.ConverterTypePair, 5)
			instance.AllowedConverteTypeMaps[0] = &converter.ConverterTypePair{SourceType: "ai", TargetType: "jpg"}
			instance.AllowedConverteTypeMaps[1] = &converter.ConverterTypePair{SourceType: "ai", TargetType: "jpeg"}
			instance.AllowedConverteTypeMaps[2] = &converter.ConverterTypePair{SourceType: "ai", TargetType: "pdf"}
			instance.AllowedConverteTypeMaps[3] = &converter.ConverterTypePair{SourceType: "pdf", TargetType: "jpg"}
			instance.AllowedConverteTypeMaps[4] = &converter.ConverterTypePair{SourceType: "pdf", TargetType: "jpeg"}
		}
		singletonMu.Unlock()
		instance.initialise()
	}
	return instance
}

func (r *AiConverter) initialise() {
	if r.internalGSInstance == nil {
		singletonMu.Lock()
		defer singletonMu.Unlock()
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
}
func (r *AiConverter) ToFastImage(inputFile string, outputFile string) error {
	return r.ToFastJpeg(inputFile, outputFile)

}
func (r *AiConverter) ToFastJpeg(inputFile string, outputFile string) error {
	gsCommandMu.Lock()
	instance.initialise()
	defer gsCommandMu.Unlock()
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
		r.internalGSInstance.Destroy()
		r.internalGSInstance = nil
		return err
	}
	defer func() {
		r.internalGSInstance.Exit()
		r.internalGSInstance.Destroy()
		r.internalGSInstance = nil
	}()

	return nil
}

func (r *AiConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	gsCommandMu.Lock()
	r.initialise()
	defer gsCommandMu.Unlock()
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
		// "-dPDFA",
		"-r300",
		"-sOutputFile=" + outputFile,
		inputFile,
	}

	if err := r.internalGSInstance.Init(args); err != nil {
		r.internalGSInstance.Destroy()
		r.internalGSInstance = nil
		return err
	}
	defer func() {
		r.internalGSInstance.Exit()
		r.internalGSInstance.Destroy()
		r.internalGSInstance = nil
	}()
	return nil
}

func (r *AiConverter) CanHandle(pair converter.ConverterTypePair) bool {
	return r.AllowedConverteTypeMaps.Exists(func(ctp *converter.ConverterTypePair) bool {

		return ctp.SourceType == pair.SourceType && ctp.TargetType == pair.TargetType
	})
}
func (r *AiConverter) Destory() {
	log.Logger.Info("AiConverter destoryed")
}
