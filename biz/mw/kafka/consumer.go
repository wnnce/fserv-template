package kafka

import (
	"context"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Consumer represents a Kafka topic consumer with batch settings and a message handler.
type Consumer struct {
	Topic         string
	Batch         bool
	BatchMaxCount int
	Handler       MessageHandler
}

func NewConsumer(topic string, batch bool, batchMaxCount int, handler MessageHandler) Consumer {
	return Consumer{
		Topic:         topic,
		Batch:         batch,
		BatchMaxCount: batchMaxCount,
		Handler:       handler,
	}
}

// MessageHandler defines the function signature for processing consumed Kafka records.
type MessageHandler func(ctx context.Context, client *kgo.Client, autoCommit bool, records ...*kgo.Record)

// RegisterConsumers registers one or more consumers to the Service and adds their topics to the client.
func (self *Service) RegisterConsumers(consumers ...Consumer) {
	if self.ctx.Err() != nil {
		return
	}
	if consumers == nil || len(consumers) == 0 {
		return
	}
	self.mutex.Lock()
	for _, consumer := range consumers {
		if consumer.Topic == "" || strings.TrimSpace(consumer.Topic) == "" {
			continue
		}
		if consumer.Batch && consumer.BatchMaxCount <= 1 {
			consumer.Batch = false
		}
		if _, ok := self.consumers[consumer.Topic]; !ok {
			self.consumers[consumer.Topic] = &consumer
			self.client.AddConsumeTopics(consumer.Topic)
		}
	}
	self.mutex.Unlock()
}

// RemoveConsumers removes consumers for the specified topics and pauses fetching for those topics.
func (self *Service) RemoveConsumers(topics ...string) {
	if self.ctx.Err() != nil {
		return
	}
	if topics == nil || len(topics) == 0 {
		return
	}
	self.mutex.Lock()
	for _, topic := range topics {
		if _, ok := self.consumers[topic]; ok {
			delete(self.consumers, topic)
			self.client.PauseFetchTopics(topic)
		}
	}
	self.mutex.Unlock()
}
