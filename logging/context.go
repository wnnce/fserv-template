package logging

import (
	"context"
	"log/slog"
	"time"
)

const (
	contextKey = "context"
)

type contextHandler struct {
	slog.Handler
	keys []string
}

func newContextHandler(parent slog.Handler, keys ...string) slog.Handler {
	return &contextHandler{
		keys:    keys,
		Handler: parent,
	}
}

func (self *contextHandler) Handle(ctx context.Context, record slog.Record) error {
	if self.keys == nil || len(self.keys) == 0 {
		return self.Handler.Handle(ctx, record)
	}
	attrs := make([]slog.Attr, 0, len(self.keys))
	for _, key := range self.keys {
		value := ctx.Value(key)
		if value == nil {
			continue
		}
		// 适配各种可能的value类型
		switch v := value.(type) {
		case string:
			attrs = append(attrs, slog.String(key, v))
		case int64:
			attrs = append(attrs, slog.Int64(key, v))
		case uint64:
			attrs = append(attrs, slog.Uint64(key, v))
		case bool:
			attrs = append(attrs, slog.Bool(key, v))
		case time.Duration:
			attrs = append(attrs, slog.Duration(key, v))
		case time.Time:
			attrs = append(attrs, slog.Time(key, v))
		default:
			attrs = append(attrs, slog.Any(key, v))
		}
	}
	record.AddAttrs(slog.Any(contextKey, slog.GroupValue(attrs...)))
	return self.Handler.Handle(ctx, record)
}
