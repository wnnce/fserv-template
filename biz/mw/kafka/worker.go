package kafka

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wnnce/fserv-template/pkg/tool"
)

// workerCtx holds the context, cancel function, topic, and channel for a topic worker.
type workerCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
	topic  string
	ticker *time.Ticker
	ch     chan *kgo.Record
	mutex  *sync.Mutex
}

// ReadLoop starts the main loop for polling and dispatching Kafka messages to topic workers.
// Each topic is handled by a separate goroutine for concurrent processing.
func (self *Service) ReadLoop() {
	if self.running.Load() {
		return
	}
	self.running.Store(true)
	defer self.running.Store(false)
	for {
		if self.ctx.Err() != nil {
			slog.Info("kafka service context is canceled, service exit")
			return
		}
		if len(self.consumers) == 0 {
			slog.Error("kafka service exit, consumer topics is empty", slog.String("name", self.name))
			return
		}
		fetches := self.client.PollFetches(self.ctx)
		if fetches.IsClientClosed() {
			slog.Error("kafka service run exit, client is closed")
			return
		}
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				slog.Error("kafka service poll fetches error", slog.Any("error", err))
			}
			continue
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			self.workerMutex.Lock()
			work, ok := self.workerMap[record.Topic]
			if !ok {
				ctx, cancel := context.WithCancel(self.ctx)
				work = &workerCtx{
					ctx:    ctx,
					cancel: cancel,
					topic:  record.Topic,
					ch:     make(chan *kgo.Record, 256),
					ticker: time.NewTicker(time.Second),
					mutex:  &sync.Mutex{},
				}
				self.workerMap[record.Topic] = work
				self.workerMutex.Unlock()
				self.mutex.Lock()
				consumer, ok := self.consumers[record.Topic]
				self.mutex.Unlock()
				if !ok {
					continue
				}
				go self.worker(work, consumer)
			}
			tool.SafeSend(work.ctx, work.ch, record)
		}
	}
}

// worker processes messages for a specific topic, supporting batch or single-message handling.
// It invokes the consumer's Handler for each batch or message.
func (self *Service) worker(ctx *workerCtx, consumer *Consumer) {
	defer func() {
		ctx.cancel()
		ctx.ticker.Stop()
	}()
	batchRecord := make([]*kgo.Record, 0)
Loop:
	for {
		if ctx.ctx.Err() != nil {
			return
		}
		select {
		case <-ctx.ctx.Done():
			slog.Info("kafka service worker exit is context canceled", slog.String("topic", ctx.topic))
			break Loop
		case record, ok := <-ctx.ch:
			if !ok {
				break Loop
			}
			if !consumer.Batch || consumer.BatchMaxCount <= 1 {
				consumer.Handler(ctx.ctx, self.client, self.autoCommit, record)
				continue
			}
			ctx.mutex.Lock()
			batchRecord = append(batchRecord, record)
			if len(batchRecord) >= consumer.BatchMaxCount {
				consumer.Handler(ctx.ctx, self.client, self.autoCommit, batchRecord...)
				batchRecord = batchRecord[:0]
			}
			ctx.mutex.Unlock()
		case <-ctx.ticker.C:
			if !consumer.Batch {
				continue
			}
			ctx.mutex.Lock()
			if len(batchRecord) > 0 {
				consumer.Handler(ctx.ctx, self.client, self.autoCommit, batchRecord...)
				batchRecord = batchRecord[:0]
			}
			ctx.mutex.Unlock()
		}
	}
	if len(batchRecord) > 0 {
		consumer.Handler(ctx.ctx, self.client, self.autoCommit, batchRecord...)
	}
}
