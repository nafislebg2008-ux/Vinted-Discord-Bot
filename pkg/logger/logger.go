package logger
import "go.uber.org/zap"
var Logger *zap.Logger
func Init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger = logger
}
func Info(msg string, fields ...zap.Field) {
	if Logger == nil {
		Init()
	}
	Logger.Info(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	if Logger == nil {
		Init()
	}
	Logger.Error(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	if Logger == nil {
		Init()
	}
	Logger.Debug(msg, fields...)
}
