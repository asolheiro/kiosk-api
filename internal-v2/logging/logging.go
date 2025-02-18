package logging

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapHandler struct {
	logger *zap.Logger
}

func NewZapHandler(logger *zap.Logger) *ZapHandler {
	return &ZapHandler{logger: logger}
}

func (h *ZapHandler) Handle(ctx context.Context, r slog.Record) error {
	var level zapcore.Level
	switch {
	case r.Level >= slog.LevelError:
		level = zapcore.ErrorLevel
	case r.Level >= slog.LevelWarn:
		level = zapcore.InfoLevel
	case r.Level >= slog.LevelInfo:
		level = zapcore.InfoLevel
	default:
		level = zapcore.DebugLevel
	}

	fields := make([]zap.Field, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields = append(fields, zapField(a))
		return true
	})

	h.logger.Log(level, r.Message, fields...)
	return nil
}

func (h *ZapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]zap.Field, 0, len(attrs))
	for _, attr := range attrs {
		fields = append(fields, zapField(attr))
	}
	return NewZapHandler(h.logger.With(fields...))
}

func (h *ZapHandler) WithGroup(name string) slog.Handler {
	return NewZapHandler(h.logger.Named(name))
}

func (h *ZapHandler) Enabled(ctx context.Context, level slog.Level) bool {
	var zapLevel zapcore.Level
	switch {
	case level >= slog.LevelError:
		zapLevel = zapcore.ErrorLevel
	case level >= slog.LevelWarn:
		zapLevel = zapcore.InfoLevel
	case level >= slog.LevelInfo:
		zapLevel = zapcore.InfoLevel
	default:
		zapLevel = zapcore.DebugLevel
	}
	return h.logger.Core().Enabled(zapLevel)
}

func zapField(attr slog.Attr) zap.Field {
	switch attr.Value.Kind() {
	case slog.KindString:
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindInt64:
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindFloat64:
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindBool:
		return zap.Bool(attr.Key, attr.Value.Bool())
	default:
		return zap.Any(attr.Key, attr.Value.Any())
	}
}
