package main

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/wnnce/fserv-template/biz/handler"
)

var (
	head  *handlerNode
	mutex = &sync.Mutex{}
)

// handlerNode represents a node in the singly linked list used to store the
// chain of registered error handlers.
type handlerNode struct {
	handler ErrorHandler
	next    *handlerNode
}

// ErrorHandler defines a function that handles an error within a Fiber context.
// It returns a potentially modified error and a boolean indicating whether the
// error was handled (true) or should be passed to the next handler (false).
type ErrorHandler func(ctx *fiber.Ctx, err error) (error, bool)

// RegisterErrorHandler registers one or more ErrorHandler functions into the global
// error handling chain. Handlers are added in a last-in-first-called order.
func RegisterErrorHandler(handlers ...ErrorHandler) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, handle := range handlers {
		node := &handlerNode{
			handler: handle,
			next:    head,
		}
		head = node
	}
}

// ChainErrorHandler traverses the registered error handler chain and attempts
// to handle the given error. If none of the handlers return handled = true,
// it falls back to the defaultErrorHandler.
func ChainErrorHandler(ctx *fiber.Ctx, err error) error {
	mutex.Lock()
	node := head
	mutex.Unlock()
	for node != nil {
		handlerErr, ok := node.handler(ctx, err)
		if ok {
			return handlerErr
		}
		node = node.next
	}
	slog.InfoContext(ctx.UserContext(), "No custom error handler handled the error, falling back to default")
	return defaultErrorHandler(ctx, err)
}

// defaultErrorHandler is the fallback handler that logs the error and sends
// a generic 500 Internal Server Error response to the client.
func defaultErrorHandler(ctx *fiber.Ctx, err error) error {
	slog.ErrorContext(ctx.UserContext(), err.Error())
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		return ctx.JSON(handler.Fail(fiberError.Code, fiberError.Message))
	}
	return ctx.JSON(handler.FailWithServer(err.Error()))
}
