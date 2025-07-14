package logging

import (
	"io"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	"github.com/wnnce/fserv-template/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggerLevelMap = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

func NewLogger() (*slog.Logger, error) {
	env := config.ViperGet[string]("server.environment", "dev")
	var handler slog.Handler
	levelStr := strings.ToUpper(config.ViperGet[string]("logger.level", "INFO"))
	level, ok := loggerLevelMap[levelStr]
	if !ok {
		level = slog.LevelInfo
	}
	options := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.ViperGet[bool]("logger.source", false),
	}
	if env == "dev" {
		handler = newConsoleHandler(slog.NewTextHandler(io.Discard, options), os.Stdout)
	} else {
		fileWrite, err := newLumberjackLogger()
		if err != nil {
			return nil, err
		}
		handler = slog.NewJSONHandler(io.MultiWriter(os.Stdout, fileWrite), options)
	}
	return slog.New(handler), nil
}

func NewLoggerWithContext(keys ...string) (*slog.Logger, error) {
	logger, err := NewLogger()
	if err != nil {
		return nil, err
	}
	handler := newContextHandler(logger.Handler(), keys...)
	return slog.New(handler), nil
}

func newLumberjackLogger() (io.Writer, error) {
	fileDir := config.ViperGet[string]("logger.file-path", "./logs/")
	if err := os.MkdirAll(fileDir, 0o777); err != nil {
		return nil, err
	}
	filename := time.Now().Format("2006-01-02") + ".log"
	fullFilename := path.Join(fileDir, filename)
	if _, err := os.Stat(fullFilename); err != nil {
		if _, err := os.Create(fullFilename); err != nil {
			return nil, err
		}
	}
	return &lumberjack.Logger{
		Filename:   fullFilename,
		MaxSize:    config.ViperGet[int]("logger.max-size", 50),
		MaxBackups: config.ViperGet[int]("logger.max-backups", 5),
		MaxAge:     config.ViperGet[int]("logger.max-age", 10),
		Compress:   config.ViperGet[bool]("logger.compress", false),
	}, nil
}
