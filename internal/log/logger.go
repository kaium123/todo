package log

import (
	"context"
	"go.uber.org/zap"
)

var TraceIDContextKey = "trid" // keep it short

// TraceID from context.
func TraceID(ctx context.Context) (traceID string) {
	traceID, _ = ctx.Value(TraceIDContextKey).(string)
	return
}

// WithTraceID sets given trace ID to given context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDContextKey, traceID) // nolint
}

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	conf := new(zap.Config)
	conf.DisableStacktrace = true
	conf.Encoding = "console"
	if err := conf.Level.UnmarshalText([]byte("debug")); err != nil {
		panic(err)
	}

	conf.EncoderConfig = zap.NewProductionEncoderConfig()
	conf.OutputPaths = []string{"stdout"}

	//zaplog, err := conf.Build()
	zaplog, err := conf.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	log := &Logger{}
	logger := zaplog
	log.Logger = logger
	return log
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (log *Logger) Named(s string) *Logger {
	named := new(Logger)
	named.Logger = log.Logger.Named(s)

	return named
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func (log *Logger) WithOptions(opts ...zap.Option) *Logger {
	n := new(Logger)
	n.Logger = log.Logger.WithOptions(opts...)
	return n
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (log *Logger) With(fields ...zap.Field) *Logger {
	n := new(Logger)
	n.Logger = log.Logger.With(fields...)
	return n
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	log.Logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = AddContextFields(ctx, fields...)
	}
	log.Logger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = AddContextFields(ctx, fields...)
	}
	log.Logger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = AddContextFields(ctx, fields...)
	}
	log.Logger.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (log *Logger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = AddContextFields(ctx, fields...)
	}
	log.Logger.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (log *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = AddContextFields(ctx, fields...)
	}
	log.Logger.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (log *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		fields = AddContextFields(ctx, fields...)
	}
	log.Logger.Fatal(msg, fields...)
}

// AddContextFields returns zap.Fields with request_id
// and trace_id if set in given context.
func AddContextFields(ctx context.Context, flds ...zap.Field) (
	all []zap.Field) {
	all = flds

	var traceID = TraceID(ctx)

	if traceID != "" {
		all = append(all, zap.String("trace_id", traceID))
	}

	return
}
