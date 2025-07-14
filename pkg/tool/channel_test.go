package tool

import (
	"context"
	"testing"
	"time"
)

func TestSafeSend(t *testing.T) {
	ch := make(chan int, 1)
	ctx := context.Background()
	SafeSend(ctx, ch, 42)
	select {
	case v := <-ch:
		if v != 42 {
			t.Errorf("expected 42, got %v", v)
		}
	default:
		t.Error("expected value in channel, got none")
	}

	// 测试 context 已取消
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	SafeSend(ctx, ch, 100)
	select {
	case v := <-ch:
		t.Errorf("expected no value, got %v", v)
	default:
		// pass
	}
}

func TestSafeSendWithCallback(t *testing.T) {
	ch := make(chan int, 1)
	ctx := context.Background()
	called := false
	SafeSendWithCallback(ctx, ch, 99, func(err error) {
		called = true
	})
	select {
	case v := <-ch:
		if v != 99 {
			t.Errorf("expected 99, got %v", v)
		}
	default:
		t.Error("expected value in channel, got none")
	}
	if called {
		t.Error("callback should not be called when context is valid")
	}

	// 测试 context 已取消
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	cancel()
	cbCalled := false
	SafeSendWithCallback(ctx, ch, 100, func(err error) {
		cbCalled = true
		if err == nil {
			t.Error("expected error in callback, got nil")
		}
	})
	if !cbCalled {
		t.Error("callback should be called when context is canceled")
	}
}
