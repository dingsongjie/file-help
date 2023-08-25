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
		if Logger == nil {
			if configs.IsGinInDebug {
				Logger, _ = zap.NewDevelopmentConfig().Build()
			} else {
				Logger, _ = zap.NewProductionConfig().Build()
			}
		}
		singletonMu.Unlock()
	}
	return Logger
}
