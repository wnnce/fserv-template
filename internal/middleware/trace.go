package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/wnnce/fserv-template/internal/constat"
)

func TraceMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		traceID := uuid.New().String()
		c := context.WithValue(ctx.Context(), constat.ContextTraceKey, traceID)
		ctx.SetContext(c)
		beginTime := time.Now().UnixMilli()
		defer func() {
			latency := time.Now().UnixMilli() - beginTime
			slog.Info(fmt.Sprintf("[ %s ] >> %s  %s  LATENCY:%dms  STATUS:%d   %s",
				ctx.Method(),
				ctx.OriginalURL(),
				traceID,
				latency,
				ctx.Response().StatusCode(),
				ctx.IP(),
			))
		}()
		return ctx.Next()
	}
}
