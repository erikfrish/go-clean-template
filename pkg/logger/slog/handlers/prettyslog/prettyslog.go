package prettyslog

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/fatih/color"
)

const prettyTimeLayout = "[2006-01-02 15:04:05.000]"

func SetupPrettySlog(out io.Writer, level slog.Level) *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}
	handler := opts.NewPrettyHandler(out)
	return slog.New(handler)
}

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	h := &PrettyHandler{
		opts:    opts,
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}
	return h
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format(prettyTimeLayout)
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
