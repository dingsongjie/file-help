package log

import (
	"www.github.com/dingsongjie/file-help/configs"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func Initialise() *zap.Logger {
	if configs.IsGinInDebug {
		Logger, _ = zap.NewDevelopmentConfig().Build()
	} else {
		Logger, _ = zap.NewProductionConfig().Build()
	}
	return Logger
}
