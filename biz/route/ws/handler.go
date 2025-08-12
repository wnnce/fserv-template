// handler.go
//
// Author:      cola
// Description: TODO: Describe this file
// Created:     2025/7/13 15:10

package ws

import (
	"log/slog"

	"github.com/gofiber/utils/v2"
	"github.com/wnnce/fserv-template/biz/handler"
)

// WebsocketHandler defines callbacks for session events and messages.
type WebsocketHandler interface {
	// OnTextMessage is called for text messages.
	OnTextMessage(ctx SessionContext, message []byte)

	// OnBinaryMessage is called for binary messages.
	OnBinaryMessage(ctx SessionContext, message []byte)

	// OnMessage handles other message types or errors.
	OnMessage(ctx SessionContext, messageType int, message []byte, err error)

	// OnPong handles pong (heartbeat) responses.
	OnPong(ctx SessionContext, data string) error

	// OnError is invoked on read errors. Return true to exit.
	OnError(ctx SessionContext, err error) bool

	// OnClose is called once when the session is closed.
	OnClose(ctx SessionContext)
}

type BuiltinWebsocketHandler struct {
}

func (_ *BuiltinWebsocketHandler) OnTextMessage(_ SessionContext, _ []byte) {}

func (_ *BuiltinWebsocketHandler) OnBinaryMessage(_ SessionContext, _ []byte) {}

func (_ *BuiltinWebsocketHandler) OnMessage(_ SessionContext, _ int, _ []byte, _ error) {
}

func (_ *BuiltinWebsocketHandler) OnPong(_ SessionContext, _ string) error {
	return nil
}

func (_ *BuiltinWebsocketHandler) OnError(_ SessionContext, _ error) bool {
	return false
}

func (_ *BuiltinWebsocketHandler) OnClose(_ SessionContext) {}

type EchoWebsocketHandler struct {
	BuiltinWebsocketHandler
}

func NewEchoWebsocketHandler() WebsocketHandler {
	return &EchoWebsocketHandler{}
}

func (self *EchoWebsocketHandler) OnTextMessage(ctx SessionContext, message []byte) {
	data := utils.UnsafeString(message)
	slog.Info("echo websocket reader text message", slog.String("message", data))
	if data == "END" {
		ctx.Shutdown()
		return
	}
	_ = ctx.WriteTextMessageWithJSON(handler.OkWithData[string](data))
}

func (self *EchoWebsocketHandler) OnError(ctx SessionContext, err error) bool {
	//TODO implement me
	return true
}

func (self *EchoWebsocketHandler) OnClose(ctx SessionContext) {
	slog.Info("EchoWebsocketHandler onClose")
}
