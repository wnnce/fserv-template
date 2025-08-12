package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func DefaultRecoverHandler(ctx *fiber.Ctx, value any) {
	slog.ErrorContext(ctx.UserContext(), "panic recovered",
		slog.String("url", ctx.OriginalURL()),
		slog.Any("error", value),
		slog.String("debug", string(debug.Stack())),
	)
}
