package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// CorsConfig defines the configuration for the CORS (Cross-Origin Resource Sharing) middleware.
type CorsConfig struct {
	UseOrigin           bool     // If true, the request's Origin header will be used as the allowed origin.
	AllowOrigin         string   // The allowed origin(s). Ignored if UseOrigin is true.
	AllowCredentials    bool     // Whether to allow credentials (cookies, authorization headers, etc.) in requests.
	AllowMethods        []string // A list of allowed HTTP methods for cross-origin requests.
	AllowHeaders        []string // A list of allowed HTTP headers in cross-origin requests.
	OptionMaxAge        int64    // The duration (in seconds) the preflight request can be cached.
	ReleaseOptionMethod bool     // If true, the middleware will respond directly to preflight OPTIONS requests.
}

// DefaultCorsConfig provides a default configuration for the CORS middleware.
var DefaultCorsConfig = CorsConfig{
	UseOrigin:           true,
	AllowOrigin:         "*",
	AllowCredentials:    true,
	AllowMethods:        []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete, fiber.MethodOptions},
	AllowHeaders:        []string{"*"},
	OptionMaxAge:        3600,
	ReleaseOptionMethod: true,
}

// corsConfigDefault fills in any missing values in the given CorsConfig with defaults
// from DefaultCorsConfig.
func corsConfigDefault(cfg *CorsConfig) {
	if cfg.AllowOrigin == "" {
		cfg.AllowOrigin = DefaultCorsConfig.AllowOrigin
	}
	if cfg.AllowMethods == nil || len(cfg.AllowMethods) == 0 {
		cfg.AllowMethods = DefaultCorsConfig.AllowMethods
	}
	if cfg.AllowHeaders == nil || len(cfg.AllowHeaders) == 0 {
		cfg.AllowHeaders = DefaultCorsConfig.AllowHeaders
	}
	if cfg.OptionMaxAge < 0 {
		cfg.OptionMaxAge = DefaultCorsConfig.OptionMaxAge
	}
}

// CorsMiddleware creates a Fiber middleware handler that sets the necessary
// CORS headers according to the provided CorsConfig.
//
// If the request method is OPTIONS and ReleaseOptionMethod is true,
// the middleware responds immediately with a 204 status and proper headers.
func CorsMiddleware(config CorsConfig) fiber.Handler {
	corsConfigDefault(&config)
	return func(ctx fiber.Ctx) error {
		if config.UseOrigin {
			ctx.Set(fiber.HeaderAccessControlAllowOrigin, ctx.Get(fiber.HeaderOrigin, ""))
		} else {
			ctx.Set(fiber.HeaderAccessControlAllowOrigin, config.AllowOrigin)
		}
		ctx.Set(fiber.HeaderAccessControlAllowCredentials, strconv.FormatBool(config.AllowCredentials))
		ctx.Set(fiber.HeaderAccessControlAllowMethods, strings.Join(config.AllowMethods, ","))
		ctx.Set(fiber.HeaderAccessControlAllowHeaders, strings.Join(config.AllowHeaders, ","))
		if ctx.Method() == fiber.MethodOptions {
			ctx.Set(fiber.HeaderAccessControlMaxAge, strconv.FormatInt(config.OptionMaxAge, 10))
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		return ctx.Next()
	}
}
