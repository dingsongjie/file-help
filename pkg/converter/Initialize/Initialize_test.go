package Initialize

import (
	"testing"

	"github.com/dingsongjie/file-help/pkg/converter"
	"github.com/dingsongjie/file-help/pkg/converter/aiconverter"
	"github.com/dingsongjie/file-help/pkg/converter/imagickconverter"
	"github.com/dingsongjie/file-help/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRegisterConverters(t *testing.T) {
	assert := assert.New(t)
	RegisterConverters()
	assert.Equal(2, len(converter.Converters))
	_, ok := (converter.Converters[0]).(*aiconverter.AiConverter)
	assert.True(ok)
	_, ok = (converter.Converters[1]).(*imagickconverter.ImagickConverter)
	assert.True(ok)
}

func TestDestoryConverters(t *testing.T) {
	log.Initialise()
	RegisterConverters()
	DestoryConverters()
}
