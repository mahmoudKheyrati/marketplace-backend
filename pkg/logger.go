package pkg

import (
	"go.uber.org/zap"
	"sync"
)

var Log *zap.SugaredLogger
var once sync.Once

func NewLogger() *zap.SugaredLogger {
	once.Do(func() {
		logger, _ := zap.NewProduction()
		Log = logger.Sugar()
	})
	return Log
}
