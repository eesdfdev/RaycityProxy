package Logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()
	sugar = logger.Sugar()
}
func Infof(message string, fields ...interface{}) {
	sugar.Infof(message, fields...)
}
func Debugf(message string, fields ...interface{}) {
	sugar.Debugf(message, fields...)
}
func Errorf(message string, fields ...interface{}) {
	sugar.Errorf(message, fields...)
}
