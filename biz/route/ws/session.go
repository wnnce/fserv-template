package ws

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/utils/v2"
)

// session.go
//
// Author:      cola
// Description: TODO: Describe this file
// Created:     2025/7/12 21:51

// SessionContext defines the minimal interface for a WebSocket session
// that handlers can use to interact with the underlying connection.
type SessionContext interface {
	// Shutdown gracefully closes the session and cleans up resources.
	Shutdown()

	// Write sends a message of the given type through the WebSocket.
	Write(messageType int, message []byte) error

	// WriteTextMessage sends a text message through the WebSocket.
	WriteTextMessage(message []byte) error

	// WriteBinaryMessage sends a binary message through the WebSocket.
	WriteBinaryMessage(data []byte) error

	// WriteTextMessageWithJSON marshals the given value to JSON
	// and sends it as a text message.
	WriteTextMessageWithJSON(message any) error

	// Websocket returns the raw *websocket.Conn for low-level operations.
	Websocket() *websocket.Conn

	// Context returns the session's context, which is canceled on shutdown.
	Context() context.Context
}

// WebsocketSession manages a single WebSocket connection, providing
// thread-safe read/write, lifecycle control, and pluggable handlers.
type WebsocketSession struct {
	conn    *websocket.Conn    // underlying WebSocket connection
	ctx     context.Context    // context controlling session lifecycle
	cancel  context.CancelFunc // function to cancel the context
	once    sync.Once          // ensures Shutdown runs only once
	handler WebsocketHandler   // message/event handler
	mutex   sync.Mutex         // protects concurrent writes
}

// NewWebsocketSession creates a new WebsocketSession with its own
// cancellable context derived from the provided parent context.
func NewWebsocketSession(ctx context.Context, conn *websocket.Conn, handler WebsocketHandler) *WebsocketSession {
	child, cancel := context.WithCancel(ctx)
	return &WebsocketSession{
		conn:    conn,
		ctx:     child,
		cancel:  cancel,
		handler: handler,
	}
}

// Shutdown closes the session exactly once, sending a close frame,
// closing the connection, canceling the context, and invoking the
// handler's OnClose callback.
func (self *WebsocketSession) Shutdown() {
	self.once.Do(func() {
		// Acquire lock to ensure outstanding writes complete first
		self.mutex.Lock()
		self.sendCloseMessage()
		_ = self.conn.Close()
		self.mutex.Unlock()

		self.cancel()
		if self.handler != nil {
			self.handler.OnClose(self)
		}
	})
}

// ReadLoop begins the read loop for the session, dispatching incoming
// messages to the handler. It recovers panics and ensures
// Shutdown is called.
func (self *WebsocketSession) ReadLoop() {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("websocket session panic exit", slog.Any("error", err))
		} else {
			slog.Info("websocket session runner exit")
		}
		self.Shutdown()
	}()
	self.conn.SetPongHandler(func(appData string) error {
		if self.handler != nil {
			return self.handler.OnPong(self, appData)
		}
		return nil
	})
	self.watchCancel()
	for {
		if err := self.ctx.Err(); err != nil {
			return
		}
		messageType, message, err := self.conn.ReadMessage()
		if err != nil && self.handler.OnError(self, err) {
			return
		}
		switch messageType {
		case websocket.TextMessage:
			self.handler.OnTextMessage(self, message)
		case websocket.BinaryMessage:
			self.handler.OnBinaryMessage(self, message)
		default:
			self.handler.OnMessage(self, messageType, message, err)
		}
	}
}

// watchCancel listens for context cancellation and triggers Shutdown.
func (self *WebsocketSession) watchCancel() {
	go func() {
		<-self.ctx.Done()
		slog.Info("websocket session context canceled")
		self.Shutdown()
	}()
}

// sendCloseMessage writes a normal closure control frame to the client.
func (self *WebsocketSession) sendCloseMessage() {
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutdown")
	_ = self.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
}

// Write sends a message of the given type in a thread-safe manner.
func (self *WebsocketSession) Write(messageType int, message []byte) error {
	if err := self.ctx.Err(); err != nil {
		return err
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.conn.WriteMessage(messageType, message)
}

// WriteTextMessage sends a text message in a thread-safe manner.
func (self *WebsocketSession) WriteTextMessage(message []byte) error {
	return self.Write(websocket.TextMessage, message)
}

// WriteBinaryMessage sends a binary message in a thread-safe manner.
func (self *WebsocketSession) WriteBinaryMessage(data []byte) error {
	return self.Write(websocket.BinaryMessage, data)
}

// WriteTextMessageWithJSON marshals the provided value to JSON or
// extracts raw bytes and sends as a text message.
func (self *WebsocketSession) WriteTextMessageWithJSON(message any) error {
	var body []byte
	switch v := message.(type) {
	case []byte:
		body = v
	case string:
		body = utils.UnsafeBytes(v)
	case strings.Builder:
		body = utils.UnsafeBytes(v.String())
	case *strings.Builder:
		body = utils.UnsafeBytes(v.String())
	case bytes.Buffer:
		body = v.Bytes()
	case *bytes.Buffer:
		body = v.Bytes()
	case io.Reader:
		data, err := io.ReadAll(v)
		if err != nil {
			return err
		}
		body = data
	default:
		data, err := sonic.Marshal(message)
		if err != nil {
			return err
		}
		body = data
	}
	return self.Write(websocket.TextMessage, body)
}

// Websocket returns the raw *websocket.Conn for advanced usage.
func (self *WebsocketSession) Websocket() *websocket.Conn {
	return self.conn
}

// Context returns the session's context, which is closed on Shutdown.
func (self *WebsocketSession) Context() context.Context {
	return self.ctx
}
