package log

import (
	"sync"

	"www.github.com/dingsongjie/file-help/configs"

	"go.uber.org/zap"
)

var (
	Logger      *zap.Logger
	singletonMu sync.Mutex = sync.Mutex{}
)

func Initialise() *zap.Logger {
	if Logger == nil {
		singletonMu.Lock()
		var config zap.Config
		if Logger == nil {
			if configs.IsGinInDebug {
				config = zap.NewDevelopmentConfig()

			} else {
				config = zap.NewProductionConfig()
			}
			config.EncoderConfig.MessageKey = "Message"
			config.EncoderConfig.TimeKey = "Timestamp"
			config.EncoderConfig.LevelKey = "LogLevel"
			config.EncoderConfig.NameKey = "Category"
			config.EncoderConfig.StacktraceKey = "Exception"
			Logger, _ = config.Build()
		}
		singletonMu.Unlock()
	}
	return Logger
}
