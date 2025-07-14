package tool

import "context"

func SafeSend[T any](ctx context.Context, ch chan<- T, data T) {
	if ctx.Err() != nil {
		return
	}
	select {
	case <-ctx.Done():
	case ch <- data:
	}
}

func SafeSendWithCallback[T any](ctx context.Context, ch chan<- T, data T, callback func(err error)) {
	if err := ctx.Err(); err != nil {
		if callback != nil {
			callback(err)
		}
		return
	}
	select {
	case <-ctx.Done():
		if callback != nil {
			callback(ctx.Err())
		}
	case ch <- data:
	}
}
