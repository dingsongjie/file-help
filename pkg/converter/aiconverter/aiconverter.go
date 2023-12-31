package aiconverter

import (
	"fmt"

	"github.com/MrSaints/go-ghostscript/ghostscript"
	"github.com/STRockefeller/go-linq"
	"github.com/dingsongjie/file-help/pkg/converter"
	"github.com/dingsongjie/file-help/pkg/log"
	"github.com/sasha-s/go-deadlock"
)

var (
	singletonMu         deadlock.Mutex = deadlock.Mutex{}
	gsCommandMu         deadlock.Mutex = deadlock.Mutex{}
	instance            *AiConverter
	gsGetRevision       = ghostscript.GetRevision
	gsNewInstance       = ghostscript.NewInstance
	logGetRevisionFaild = func(err error, rev ghostscript.Revision) {
		log.Logger.Sugar().Fatalf("Revision: %+v\n", rev)
	}
	logNewInstanceFaild = func(err error) {
		log.Logger.Sugar().Fatalf("Error: %+v\n", err)
	}
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
			rev, err := gsGetRevision()
			if err != nil {
				logGetRevisionFaild(err, rev)
			}

			r.internalGSInstance, err = gsNewInstance()
			if err != nil {
				logNewInstanceFaild(err)
			}
		}

	}
}
func (r *AiConverter) ToFastImage(inputFile string, outputFile string, dpi int) error {
	if dpi == 0 {
		dpi = 72
	}
	if dpi > 300 {
		return fmt.Errorf("dpi is not allowed to exceed 300")
	}
	return r.ToFastJpeg(inputFile, outputFile, dpi)

}
func (r *AiConverter) ToFastJpeg(inputFile string, outputFile string, dpi int) error {
	gsCommandMu.Lock()
	defer gsCommandMu.Unlock()
	r.initialise()
	defer func() {
		if r.internalGSInstance != nil {
			r.internalGSInstance.Destroy()
			r.internalGSInstance = nil
		}
	}()

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
		fmt.Sprint("-r", dpi),
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
		log.Logger.Error(err.Error())
		return fmt.Errorf("convert faild")
	}
	defer func() {
		r.internalGSInstance.Exit()
	}()
	return nil
}

func (r *AiConverter) ToPrettyPdf(inputFile string, outputFile string) error {
	gsCommandMu.Lock()
	defer gsCommandMu.Unlock()
	r.initialise()
	defer func() {
		r.internalGSInstance.Destroy()
		r.internalGSInstance = nil
	}()
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
		log.Logger.Error(err.Error())
		return fmt.Errorf("convert faild")
	}
	defer func() {
		r.internalGSInstance.Exit()
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
