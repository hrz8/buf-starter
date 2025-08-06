package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/hrz8/altalune"
)

type SlogLogger struct {
	*slog.Logger
	jsonLogger *SlogLogger
}

var _ altalune.Logger = (*SlogLogger)(nil)

type logHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *logHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = color.BlueString(level)
	case slog.LevelInfo:
		level = color.GreenString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05.000]")
	if len(fields) > 0 {
		h.l.Println(timeStr, level, r.Message, color.WhiteString(string(b)))
	} else {
		h.l.Println(timeStr, level, r.Message)
	}

	return nil
}

func newLogHandler(w io.Writer, opts *slog.HandlerOptions) *logHandler {
	return &logHandler{
		Handler: slog.NewJSONHandler(w, opts),
		l:       log.New(w, "", 0),
	}
}

func New(lvl string) *SlogLogger {
	var level slog.Level
	switch lvl {
	case "debug", "DEBUG", "verbose", "VERBOSE":
		level = slog.LevelDebug
	case "info", "INFO":
		level = slog.LevelInfo
	case "warn", "WARN", "warning", "WARNING":
		level = slog.LevelWarn
	case "err", "ERR", "error", "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelWarn
	}

	consoleLog := slog.New(newLogHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	jsonLog := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

	jsonLogger := &SlogLogger{
		Logger:     jsonLog,
		jsonLogger: nil,
	}

	return &SlogLogger{
		Logger:     consoleLog,
		jsonLogger: jsonLogger,
	}
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

func (l *SlogLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Logger.DebugContext(ctx, msg, args...)
}

func (l *SlogLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Logger.InfoContext(ctx, msg, args...)
}

func (l *SlogLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Logger.WarnContext(ctx, msg, args...)
}

func (l *SlogLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.Logger.ErrorContext(ctx, msg, args...)
}

func (l *SlogLogger) JSON() altalune.Logger {
	if l.jsonLogger == nil {
		return l
	}
	return l.jsonLogger
}
