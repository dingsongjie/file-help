package Initialize

import (
	"github.com/dingsongjie/file-help/pkg/converter"
	"github.com/dingsongjie/file-help/pkg/converter/aiconverter"
	"github.com/dingsongjie/file-help/pkg/converter/imagickconverter"
)

func RegisterConverters() {
	converter.Converters = make([]converter.Converter, 2)
	converter.Converters[0] = aiconverter.NewConverter()
	// converter.Converters[1] = cdrconverter.NewConverter()
	converter.Converters[1] = imagickconverter.NewConverter()
}

func DestoryConverters() {
	for i := range converter.Converters {
		converter.Converters[i].Destory()
	}
}
