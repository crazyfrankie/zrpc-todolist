package logs

import (
	"context"
	"os"
	"sync"
)

var (
	globalLogger FullLogger
	globalMutex  sync.RWMutex
)

func SetGlobalLogger(logger FullLogger) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	globalLogger = logger
}

func getGlobalLogger() FullLogger {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	if globalLogger == nil {
		globalLogger = NewLogger(os.Stderr)
	}
	return globalLogger
}

func Info(v ...any) {
	getGlobalLogger().Info(v...)
}

func Infof(format string, v ...any) {
	getGlobalLogger().Infof(format, v...)
}

func Debug(v ...any) {
	getGlobalLogger().Debug(v...)
}

func Debugf(format string, v ...any) {
	getGlobalLogger().Debugf(format, v...)
}

func Warn(v ...any) {
	getGlobalLogger().Warn(v...)
}

func Warnf(format string, v ...any) {
	getGlobalLogger().Warnf(format, v...)
}

func Error(v ...any) {
	getGlobalLogger().Error(v...)
}

func Errorf(format string, v ...any) {
	getGlobalLogger().Errorf(format, v...)
}

func Fatal(v ...any) {
	getGlobalLogger().Fatal(v...)
}

func Fatalf(format string, v ...any) {
	getGlobalLogger().Fatalf(format, v...)
}

func CtxTracef(ctx context.Context, format string, v ...any) {
	getGlobalLogger().CtxTracef(ctx, format, v...)
}

func CtxDebugf(ctx context.Context, format string, v ...any) {
	getGlobalLogger().CtxDebugf(ctx, format, v...)
}

func CtxInfof(ctx context.Context, format string, v ...any) {
	getGlobalLogger().CtxInfof(ctx, format, v...)
}

func CtxNoticef(ctx context.Context, format string, v ...any) {
	getGlobalLogger().CtxNoticef(ctx, format, v...)
}

func CtxWarnf(ctx context.Context, format string, v ...any) {
	getGlobalLogger().CtxWarnf(ctx, format, v...)
}

func CtxErrorf(ctx context.Context, format string, v ...any) {
	getGlobalLogger().CtxErrorf(ctx, format, v...)
}

func CtxFatalf(ctx context.Context, format string, v ...any) {
	getGlobalLogger().Fatalf(format, v...)
}
