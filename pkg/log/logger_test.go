package log

import (
	"testing"

	"github.com/dingsongjie/file-help/configs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitialise(t *testing.T) {
	assert := assert.New(t)
	t.Run("when in debug", func(t *testing.T) {
		configs.IsGinInDebug = true
		Initialise()
		level := Logger.Level()
		assert.Equal(zap.DebugLevel, level)
		//清空Logger
		Logger = nil
	})

	t.Run("when not in debug", func(t *testing.T) {
		configs.IsGinInDebug = false
		Initialise()
		level := Logger.Level()
		assert.Equal(zap.InfoLevel, level)
	})

}
