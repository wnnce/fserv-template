package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	_ "github.com/wnnce/fserv-template/biz/dal"
	_ "github.com/wnnce/fserv-template/biz/mw"
	"github.com/wnnce/fserv-template/biz/route"
	"github.com/wnnce/fserv-template/config"
	"github.com/wnnce/fserv-template/internal/constat"
	"github.com/wnnce/fserv-template/internal/middleware"
	"github.com/wnnce/fserv-template/logging"
	"github.com/wnnce/fserv-template/pkg/tool"
)

func initialize() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         config.ViperGet[string]("server.name", "fserv-template"),
		JSONEncoder:     sonic.Marshal,
		JSONDecoder:     sonic.Unmarshal,
		ErrorHandler:    ChainErrorHandler,
		StructValidator: tool.Validator(),
	})
	app.Use(recoverer.New(recoverer.Config{
		EnableStackTrace:  true,
		StackTraceHandler: middleware.DefaultRecoverHandler,
	}))
	app.Use(middleware.CorsMiddleware(middleware.DefaultCorsConfig))
	if config.ViperGet[string]("server.environment", "test") == "dev" {
		app.Use(pprof.New())
	}
	app.Use(middleware.TraceMiddleware())
	route.RegisterRouter(app)
	customRouter(app)
	return app
}

func bootstrap(cancel context.CancelFunc, app *fiber.App) {
	host := config.ViperGet[string]("server.host", "127.0.0.1")
	port := config.ViperGet[int]("server.port", 7000)
	address := host + ":" + strconv.Itoa(port)
	if err := app.Listen(address); err != nil {
		cancel()
	}
}

func main() {
	if err := config.LoadConfigWithFile("config.yaml", "./configs"); err != nil {
		panic(err)
	}
	logger, err := logging.NewLoggerWithContext(constat.ContextTraceKey, constat.ContextUserIDKey)
	if err != nil {
		panic(err)
	}
	slog.SetDefault(logger)
	ctx, cancel := context.WithCancel(context.Background())
	cleanup, err := config.DoReaderConfiguration(ctx)
	if err != nil {
		panic(err)
	}

	defer func() {
		cleanup()
		cancel()
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	app := initialize()
	go bootstrap(cancel, app)

	select {
	case <-exit:
		slog.Info("listen system exit signal, shutdown application!")
		if err := app.Shutdown(); err != nil {
			slog.Info("shutdown app error", slog.String("error", err.Error()))
		}
	case <-ctx.Done():
		slog.Info("context canceled application exit")
	}
}
