package logging

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"path"
	"runtime"
	"strconv"
	"sync"

	"github.com/wnnce/fserv-template/config"
)

const (
	reset = "\033[0m"

	black        = "30"
	red          = "31"
	green        = "32"
	yellow       = "33"
	blue         = "34"
	magenta      = "35"
	cyan         = "36"
	lightGray    = "37"
	darkGray     = "90"
	lightRed     = "91"
	lightGreen   = "92"
	lightYellow  = "93"
	lightBlue    = "94"
	lightMagenta = "95"
	lightCyan    = "96"
	white        = "97"
)

type consoleHandler struct {
	slog.Handler
	pool  *sync.Pool
	wr    io.Writer
	mutex *sync.Mutex
}

func newConsoleHandler(parent slog.Handler, wr io.Writer) slog.Handler {
	return &consoleHandler{
		Handler: parent,
		pool: &sync.Pool{
			New: func() any {
				return &bytes.Buffer{}
			},
		},
		wr:    wr,
		mutex: &sync.Mutex{},
	}
}

func (self *consoleHandler) Handle(_ context.Context, record slog.Record) error {
	builder := self.pool.Get().(*bytes.Buffer)
	self.colorize(builder, lightGray, record.Time.Format("2006-01-02 15:04:05.000"))
	builder.WriteByte(' ')
	switch record.Level {
	case slog.LevelDebug:
		self.colorize(builder, lightGray, record.Level.String())
	case slog.LevelInfo:
		self.colorize(builder, cyan, record.Level.String())
	case slog.LevelWarn:
		self.colorize(builder, lightYellow, record.Level.String())
	case slog.LevelError:
		self.colorize(builder, lightRed, record.Level.String())
	default:
		builder.WriteString(record.Level.String())
	}
	builder.WriteByte(' ')
	if config.ViperGet[bool]("logger.source", false) {
		sr := source(&record)
		_, file := path.Split(sr.Function)
		self.colorize(builder, lightBlue, file+":"+strconv.Itoa(sr.Line))
		builder.WriteString(" - ")
	}
	self.colorize(builder, green, record.Message)
	record.Attrs(func(attr slog.Attr) bool {
		switch attr.Key {
		case slog.TimeKey, slog.LevelKey, slog.MessageKey:
			return true
		default:
			builder.WriteByte(' ')
			builder.WriteString(attr.String())
		}
		return true
	})
	builder.WriteByte('\n')
	self.mutex.Lock()
	defer func() {
		self.mutex.Unlock()
		builder.Reset()
		self.pool.Put(builder)
	}()
	_, err := self.wr.Write(builder.Bytes())
	return err
}

func source(record *slog.Record) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{record.PC})
	f, _ := fs.Next()
	return &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}

func (self *consoleHandler) colorize(buffer *bytes.Buffer, colorCode, value string) {
	buffer.WriteString("\033[")
	buffer.WriteString(colorCode)
	buffer.WriteByte('m')
	buffer.WriteString(value)
	buffer.WriteString(reset)
}
