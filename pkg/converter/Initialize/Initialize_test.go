package Initialize

import (
	"testing"

	"www.github.com/dingsongjie/file-help/pkg/converter"
	"www.github.com/dingsongjie/file-help/pkg/converter/aiconverter"
	"www.github.com/dingsongjie/file-help/pkg/converter/cdrconverter"
	"www.github.com/dingsongjie/file-help/pkg/converter/imagickconverter"
)

func TestRegisterConverters(t *testing.T) {
	RegisterConverters()
	err := "initialize not successd"
	if len(converter.Converters) != 3 {
		t.Errorf(err)
	}

	if _, ok := (converter.Converters[0]).(*aiconverter.AiConverter); ok {
		t.Errorf(err)
	}
	if _, ok := (converter.Converters[1]).(*cdrconverter.CdrConverter); ok {
		t.Errorf(err)
	}
	if _, ok := (converter.Converters[2]).(*imagickconverter.ImagickConverter); ok {
		t.Errorf(err)
	}
}

func TestDestoryConverters() {
	for i := range converter.Converters {
		converter.Converters[i].Destory()
	}
}
