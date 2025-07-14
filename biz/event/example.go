// example.go
// Author:      cola
// Description: TODO: Describe this file
// Created:     2025/7/13 21:51

package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	ExampleTopic = "template-example"
)

type ExampleEvent struct {
	Offset    int   `json:"offset"`
	Timestamp int64 `json:"timestamp"`
}

func NewExampleEvent(offset int) ExampleEvent {
	return ExampleEvent{
		Offset:    offset,
		Timestamp: time.Now().UnixMilli(),
	}
}

func HandlerExampleEvent(_ context.Context, _ *kgo.Client, _ bool, records ...*kgo.Record) {
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		key, value, topic := record.Key, record.Value, record.Topic
		slog.Info(
			"handler example event message",
			slog.String("topic", topic),
			slog.String("key", string(key)),
			slog.String("value", string(value)),
		)
	}
}
