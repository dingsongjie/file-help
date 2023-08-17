package Initialize

import (
	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/converter/aiconverter"
	"www.github.com/dingsongjie/file-help/pkg/converter/cdrconverter"
	"www.github.com/dingsongjie/file-help/pkg/converter/imagickconverter"
)

var (
	Converters []converter.Converter
)

func RegisterConverters() {
	Converters = make([]converter.Converter, 3)
	Converters[0] = aiconverter.NewConverter()
	Converters[1] = cdrconverter.NewConverter()
	Converters[2] = imagickconverter.NewConverter()
}

func DestoryConverters() {
	for i := range Converters {
		Converters[i].Destory()
	}
}
