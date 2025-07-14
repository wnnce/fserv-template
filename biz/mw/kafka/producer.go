package kafka

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/utils/v2"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wnnce/fserv-template/internal/constat"
)

// ProducerWithSync sends a message to the specified topic synchronously.
// The value is automatically serialized, and traceId is injected into headers if present in context.
func (self *Service) ProducerWithSync(ctx context.Context, topic string, key []byte, value any, headers ...kgo.RecordHeader) error {
	body, err := self.valueProtocol(value)
	if err != nil {
		return err
	}
	record := &kgo.Record{
		Topic:   topic,
		Key:     key,
		Value:   body,
		Headers: self.injectTraceID(ctx, headers),
	}
	slog.DebugContext(ctx, "sync producing kafka message", slog.Group("data",
		slog.String("name", self.name),
		slog.String("topic", topic),
		slog.Any("key", key),
		slog.Any("value", value),
	))
	return self.client.ProduceSync(ctx, record).FirstErr()
}

// ProducerWithAsync sends a message to the specified topic asynchronously.
// The value is automatically serialized, and traceId is injected into headers if present in context.
// The callback is invoked upon completion.
func (self *Service) ProducerWithAsync(ctx context.Context, topic string, key []byte, value any,
	callback func(record *kgo.Record, err error), headers ...kgo.RecordHeader) {
	body, err := self.valueProtocol(value)
	if err != nil {
		if callback != nil {
			callback(nil, err)
		}
		return
	}
	record := &kgo.Record{
		Topic:   topic,
		Key:     key,
		Value:   body,
		Headers: self.injectTraceID(ctx, headers),
	}
	slog.DebugContext(ctx, "async producing kafka message", slog.Group("data",
		slog.String("name", self.name),
		slog.String("topic", topic),
		slog.Any("key", key),
		slog.Any("value", value),
	))
	self.client.Produce(ctx, record, callback)
}

// injectTraceID injects the traceId from context into Kafka record headers if not already present.
func (self *Service) injectTraceID(ctx context.Context, headers []kgo.RecordHeader) []kgo.RecordHeader {
	traceID, ok := ctx.Value(constat.ContextTraceKey).(string)
	if !ok {
		return headers
	}
	for _, h := range headers {
		if h.Key == constat.ContextTraceKey {
			return headers
		}
	}
	return append(headers, kgo.RecordHeader{Key: constat.ContextTraceKey, Value: utils.UnsafeBytes(traceID)})
}

// valueProtocol serializes the value to a byte slice, supporting various types including []byte, string, io.Reader, and struct.
func (_ *Service) valueProtocol(value any) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case strings.Builder:
		return []byte(v.String()), nil
	case *strings.Builder:
		return []byte(v.String()), nil
	case bytes.Buffer:
		return v.Bytes(), nil
	case *bytes.Buffer:
		return v.Bytes(), nil
	case io.Reader:
		return io.ReadAll(v)
	default:
		return sonic.Marshal(value)
	}
}
