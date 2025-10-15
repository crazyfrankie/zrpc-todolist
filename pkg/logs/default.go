package logs

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type BaseLogger struct {
	stdlog     *log.Logger
	level      Level
	enableCall bool
	callerSkip int
	prefix     []any
}

func NewLogger(output io.Writer) FullLogger {
	return &BaseLogger{
		stdlog: log.New(output, "", log.LstdFlags|log.Lmicroseconds),
		level:  LevelDebug,
		prefix: make([]any, 0),
	}
}

func (b *BaseLogger) With(kv ...any) {
	kvs := make([]interface{}, 0, len(b.prefix)+len(kv))
	kvs = append(kvs, kv...)
	kvs = append(kvs, b.prefix...)

	b.prefix = kvs
}

func (b *BaseLogger) WithCaller() {
	b.enableCall = true
}

func (b *BaseLogger) WithCallerSkip(skip int) {
	b.callerSkip += skip
}

func (b *BaseLogger) SetOutput(w io.Writer) {
	b.stdlog.SetOutput(w)
}

func (b *BaseLogger) SetLevel(lv Level) {
	b.level = lv
}

func (b *BaseLogger) buildPrefix() string {
	if len(b.prefix) == 0 {
		return ""
	}
	var parts []string
	for i := 0; i < len(b.prefix); i += 2 {
		if i+1 < len(b.prefix) {
			parts = append(parts, fmt.Sprintf("%v=%v", b.prefix[i], b.prefix[i+1]))
		}
	}
	return "[" + strings.Join(parts, " ") + "] "
}

func (b *BaseLogger) getCaller() string {
	if !b.enableCall {
		return ""
	}
	_, file, line, ok := runtime.Caller(4 + b.callerSkip)
	if !ok {
		return ""
	}
	if idx := strings.LastIndexByte(file, '/'); idx != -1 {
		file = file[idx+1:]
	}
	return fmt.Sprintf("%s:%d ", file, line)
}

func (b *BaseLogger) logf(lv Level, format *string, v ...any) {
	if b.level > lv {
		return
	}
	msg := b.getCaller() + lv.String() + b.buildPrefix()
	if format != nil {
		msg += fmt.Sprintf(*format, v...)
	} else {
		msg += fmt.Sprint(v...)
	}
	b.stdlog.Output(4, msg)
	if lv == LevelFatal {
		os.Exit(1)
	}
}

func (b *BaseLogger) logfCtx(ctx context.Context, lv Level, format *string, v ...any) {
	if b.level > lv {
		return
	}
	msg := b.getCaller() + lv.String() + b.buildPrefix()
	if traceID := ctx.Value("trace_id"); traceID != nil {
		msg += fmt.Sprintf("[trace_id=%v] ", traceID)
	}
	if format != nil {
		msg += fmt.Sprintf(*format, v...)
	} else {
		msg += fmt.Sprint(v...)
	}
	b.stdlog.Output(4, msg)
	if lv == LevelFatal {
		os.Exit(1)
	}
}

func (b *BaseLogger) Trace(v ...any) {
	b.logf(LevelTrace, nil, v...)
}

func (b *BaseLogger) Debug(v ...any) {
	b.logf(LevelDebug, nil, v...)
}

func (b *BaseLogger) Info(v ...any) {
	b.logf(LevelInfo, nil, v...)
}

func (b *BaseLogger) Notice(v ...any) {
	b.logf(LevelNotice, nil, v...)
}

func (b *BaseLogger) Warn(v ...any) {
	b.logf(LevelWarn, nil, v...)
}

func (b *BaseLogger) Error(v ...any) {
	b.logf(LevelError, nil, v...)
}

func (b *BaseLogger) Fatal(v ...any) {
	b.logf(LevelFatal, nil, v...)
}

func (b *BaseLogger) Tracef(format string, v ...any) {
	b.logf(LevelTrace, &format, v...)
}

func (b *BaseLogger) Debugf(format string, v ...any) {
	b.logf(LevelDebug, &format, v...)
}

func (b *BaseLogger) Infof(format string, v ...any) {
	b.logf(LevelInfo, &format, v...)
}

func (b *BaseLogger) Noticef(format string, v ...any) {
	b.logf(LevelNotice, &format, v...)
}

func (b *BaseLogger) Warnf(format string, v ...any) {
	b.logf(LevelWarn, &format, v...)
}

func (b *BaseLogger) Errorf(format string, v ...any) {
	b.logf(LevelError, &format, v...)
}

func (b *BaseLogger) Fatalf(format string, v ...any) {
	b.logf(LevelFatal, &format, v...)
}

func (b *BaseLogger) CtxTracef(ctx context.Context, format string, v ...any) {
	b.logfCtx(ctx, LevelTrace, &format, v...)
}

func (b *BaseLogger) CtxDebugf(ctx context.Context, format string, v ...any) {
	b.logfCtx(ctx, LevelDebug, &format, v...)
}

func (b *BaseLogger) CtxInfof(ctx context.Context, format string, v ...any) {
	b.logfCtx(ctx, LevelInfo, &format, v...)
}

func (b *BaseLogger) CtxNoticef(ctx context.Context, format string, v ...any) {
	b.logfCtx(ctx, LevelNotice, &format, v...)
}

func (b *BaseLogger) CtxWarnf(ctx context.Context, format string, v ...any) {
	b.logfCtx(ctx, LevelWarn, &format, v...)
}

func (b *BaseLogger) CtxErrorf(ctx context.Context, format string, v ...any) {
	b.logfCtx(ctx, LevelError, &format, v...)
}

func (b *BaseLogger) CtxFatalf(ctx context.Context, format string, v ...any) {
	b.logfCtx(ctx, LevelFatal, &format, v...)
}
