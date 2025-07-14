package main

import (
	"context"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/wnnce/fserv-template/biz/route"
	"github.com/wnnce/fserv-template/biz/route/ws"
)

var upgrade websocket.FastHTTPUpgrader

type User struct {
	Name string `json:"name,omitempty" validate:"required,min=3"`
	Age  int    `json:"age,omitempty" validate:"required,min=1"`
}

// custom router register
func customRouter(app *fiber.App) {
	app.Get("/ping", ping)
	app.Get("/ws/echo", func(ctx fiber.Ctx) error {
		return upgrade.Upgrade(ctx.RequestCtx(), func(conn *websocket.Conn) {
			session := ws.NewWebsocketSession(context.Background(), conn)
			echoHandler := ws.NewEchoWebsocketHandler(session)
			session.SetHandler(echoHandler)
			// Sync
			session.ReadLoop()
		})
	})
	app.Get("/validator", func(ctx fiber.Ctx) error {
		user := &User{}
		if err := ctx.Bind().Body(user); err != nil {
			return err
		}
		return ctx.JSON(route.OkWithData(user))
	})
}

func ping(ctx fiber.Ctx) error {
	return ctx.JSON(route.Ok())
}
