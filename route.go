package main

import (
	"context"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/wnnce/fserv-template/biz/handler"
	"github.com/wnnce/fserv-template/biz/route/ws"
	"github.com/wnnce/fserv-template/pkg/tool"
)

var upgrade = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

type User struct {
	Name string `json:"name,omitempty" validate:"required,min=3"`
	Age  int    `json:"age,omitempty" validate:"required,min=1"`
}

// custom router register
func customRouter(app *fiber.App) {
	app.Get("/health", health)
	app.Get("/ws/echo", func(ctx *fiber.Ctx) error {
		return upgrade.Upgrade(ctx.Context(), func(conn *websocket.Conn) {
			session := ws.NewWebsocketSession(context.Background(), conn, ws.NewEchoWebsocketHandler())
			// Sync
			session.ReadLoop()
		})
	})
	app.Get("/validator", func(ctx *fiber.Ctx) error {
		user := &User{}
		if err := ctx.BodyParser(user); err != nil {
			return err
		}
		if err := tool.Validator().Validate(user); err != nil {
			return err
		}
		return ctx.JSON(handler.OkWithData(user))
	})
}

func health(ctx *fiber.Ctx) error {
	return ctx.JSON(handler.Ok())
}
