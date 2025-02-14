package logs

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Warn(message string, args ...any)
	Info(message string, args ...any)
	Error(message string, args ...any)
	Debug(message string, args ...any)
	Log(message string, args ...any)
}

type defaultHandler struct {
	opts   *slog.HandlerOptions
	logger *log.Logger
	attrs  []slog.Attr
}

func newDefaultHandler(w io.Writer, opts *slog.HandlerOptions) *defaultHandler {
	return &defaultHandler{
		opts:   opts,
		logger: log.New(w, "", log.LstdFlags),
	}
}

func (h *defaultHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level.Level() >= h.opts.Level.Level()
}

func (h *defaultHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	level := r.Level.String()
	msg := r.Message

	var service string
	attrs := make([]string, 0, r.NumAttrs()+len(h.attrs))

	// Add stored attributes
	for _, attr := range h.attrs {
		if attr.Key == "service" {
			service = attr.Value.String()
		} else {
			attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any()))
		}
	}

	// Add record attributes
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "service" {
			service = a.Value.String()
		} else {
			attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value.Any()))
		}
		return true
	})

	var output string
	if service != "" {
		output = fmt.Sprintf("%s [%s]: %s", level, service, msg)
	} else {
		output = fmt.Sprintf("%s %s", level, msg)
	}

	if len(attrs) > 0 {
		output += " " + strings.Join(attrs, " ")
	}

	h.logger.Print(output)
	return nil
}

func (h *defaultHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

func (h *defaultHandler) WithGroup(name string) slog.Handler {
	return h
}

type logger struct {
	*slog.Logger
}

func (l logger) Log(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func NewSlogLogger(serviceName string) logger {
	debugMode := os.Getenv("LOG_LEVEL") == "debug"
	useJSON := os.Getenv("LOG_TYPE") == "json"
	logLevel := slog.LevelInfo

	if debugMode {
		logLevel = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}

	var handler slog.Handler
	if useJSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = newDefaultHandler(os.Stdout, opts)
	}

	l := slog.New(handler).With(
		slog.String("service", serviceName),
	)
	return logger{l}
}

var DefaultLogger = NewSlogLogger("DEFAULT")
