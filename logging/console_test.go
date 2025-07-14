package logging

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

func BenchmarkConsoleHandler(b *testing.B) {
	handler := newConsoleHandler(nil, io.Discard) // 不输出
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "message", 0)
	rec.AddAttrs(slog.String("user", "alice"), slog.Int("id", 42))
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(context.Background(), rec)
	}
}

func BenchmarkContextHandlerAndConsole(b *testing.B) {
	handler := newContextHandler(newConsoleHandler(nil, io.Discard), "traceId")
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "message", 0)
	rec.AddAttrs(slog.String("user", "alice"), slog.Int("id", 42))
	ctx := context.WithValue(context.Background(), "traceId", "231231")
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(ctx, rec)
	}
}
func BenchmarkContextHandlerAndJSON(b *testing.B) {
	handler := newContextHandler(slog.NewJSONHandler(io.Discard, nil), "traceId")
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "message", 0)
	rec.AddAttrs(slog.String("user", "alice"), slog.Int("id", 42))
	ctx := context.WithValue(context.Background(), "traceId", "231231")
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(ctx, rec)
	}
}
