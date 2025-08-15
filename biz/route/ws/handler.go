// handler.go
//
// Author:      cola
// Description: TODO: Describe this file
// Created:     2025/7/13 15:10

package ws

import (
	"log/slog"
	"time"

	"github.com/fasthttp/websocket"
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

	// OnPing handles ping (heartbeat) responses.
	OnPing(ctx SessionContext, data string) error

	// OnPong handles pong (heartbeat) responses.
	OnPong(ctx SessionContext, data string) error

	// OnClose is called once when the session is closed.
	OnClose(ctx SessionContext, err error)
}

type BuiltinWebsocketHandler struct {
}

func (_ *BuiltinWebsocketHandler) OnTextMessage(_ SessionContext, _ []byte) {}

func (_ *BuiltinWebsocketHandler) OnBinaryMessage(_ SessionContext, _ []byte) {}

func (_ *BuiltinWebsocketHandler) OnMessage(_ SessionContext, _ int, _ []byte, _ error) {
}

func (_ *BuiltinWebsocketHandler) OnPing(ctx SessionContext, data string) error {
	return ctx.Websocket().WriteControl(websocket.PongMessage, []byte(data), time.Now().Add(time.Second))
}

func (_ *BuiltinWebsocketHandler) OnPong(_ SessionContext, _ string) error {
	return nil
}

func (_ *BuiltinWebsocketHandler) OnClose(_ SessionContext, _ error) {}

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
