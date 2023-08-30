package converter

import "github.com/STRockefeller/go-linq"

type ConverterTypePair struct {
	SourceType string
	TargetType string
}

type Converter interface {
	ToFastImage(inputFile string, outputFile string) error
	ToPrettyPdf(inputFile string, outputFile string) error
	CanHandle(pair ConverterTypePair) bool
	Destory()
}

var (
	Converters linq.Linq[Converter]
)
