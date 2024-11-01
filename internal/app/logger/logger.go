package logger

import "go.uber.org/zap"

var globalLogger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	globalLogger = l.Sugar()
}

func DebugKV(message string, kvs ...interface{}) {
	globalLogger.Debugw(message, kvs...)
}

func InfoKV(message string, kvs ...interface{}) {
	globalLogger.Infow(message, kvs...)
}

func WarnKV(message string, kvs ...interface{}) {
	globalLogger.Warnw(message, kvs...)
}

func ErrorKV(message string, kvs ...interface{}) {
	globalLogger.Errorw(message, kvs...)
}

func PanicKV(message string, kvs ...interface{}) {
	globalLogger.Panicw(message, kvs...)
}

func FatalKV(message string, kvs ...interface{}) {
	globalLogger.Fatalw(message, kvs...)
}
