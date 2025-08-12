package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wnnce/fserv-template/internal/constat"
)

func TraceMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		traceID := uuid.New().String()
		c := context.WithValue(ctx.UserContext(), constat.ContextTraceKey, traceID)
		ctx.SetUserContext(c)
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
