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
	// TextMessageHandler is called for text messages.
	TextMessageHandler(message []byte)

	// BinaryMessageHandler is called for binary messages.
	BinaryMessageHandler(message []byte)

	// MessageHandler handles other message types or errors.
	MessageHandler(messageType int, message []byte, err error)

	// PongHandler handles pong (heartbeat) responses.
	PongHandler(data string) error

	// OnError is invoked on read errors. Return true to exit.
	OnError(err error) bool

	// OnClose is called once when the session is closed.
	OnClose()
}

type EchoWebsocketHandler struct {
	ctx SessionContext
}

func NewEchoWebsocketHandler(ctx SessionContext) WebsocketHandler {
	return &EchoWebsocketHandler{
		ctx: ctx,
	}
}

func (self *EchoWebsocketHandler) TextMessageHandler(message []byte) {
	data := utils.UnsafeString(message)
	slog.Info("echo websocket reader text message", slog.String("message", data))
	if data == "END" {
		self.ctx.Shutdown()
		return
	}
	_ = self.ctx.WriteTextMessageWithJSON(handler.OkWithData[string](data))
}

func (self *EchoWebsocketHandler) BinaryMessageHandler(message []byte) {
	//TODO implement me
	panic("implement me")
}

func (self *EchoWebsocketHandler) MessageHandler(messageType int, message []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (self *EchoWebsocketHandler) PongHandler(data string) error {
	//TODO implement me
	panic("implement me")
}

func (self *EchoWebsocketHandler) OnError(err error) bool {
	//TODO implement me
	return true
}

func (self *EchoWebsocketHandler) OnClose() {
	//TODO implement me
	slog.Info("EchoWebsocketHandler onClose")
}
