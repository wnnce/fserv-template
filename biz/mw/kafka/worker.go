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
	ch     chan *kgo.Record
}

// Run starts the main loop for polling and dispatching Kafka messages to topic workers.
// Each topic is handled by a separate goroutine for concurrent processing.
func (self *Service) Run() {
	if self.runner.Load() {
		return
	}
	go func() {
		self.runner.Store(true)
		defer self.runner.Store(false)
		workerCtxMap := make(map[string]*workerCtx)
		for {
			if self.ctx.Err() != nil {
				slog.Info("kafka service context is canceled, service exit")
				return
			}
			func() {
				defer func() {
					if err := recover(); err != nil {
						slog.Error("kafka service run is panic", slog.Any("error", err), slog.String("name", self.name))
					}
				}()
				self.mutex.Lock()
				if len(self.consumers) == 0 {
					slog.Error("kafka service exit, consumer topics is empty", slog.String("name", self.name))
					self.mutex.Unlock()
					return
				}
				consumersClone := make(map[string]*Consumer, len(self.consumers))
				for topic, handler := range self.consumers {
					consumersClone[topic] = handler
				}
				self.mutex.Unlock()
				for {
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
					wg := &sync.WaitGroup{}
					iter := fetches.RecordIter()
					for !iter.Done() {
						record := iter.Next()
						work, ok := workerCtxMap[record.Topic]
						if !ok {
							ctx, cancel := context.WithCancel(self.ctx)
							work = &workerCtx{
								ctx:    ctx,
								cancel: cancel,
								topic:  record.Topic,
								ch:     make(chan *kgo.Record, 1024),
							}
							workerCtxMap[record.Topic] = work
							wg.Add(1)
							go self.worker(work, wg, consumersClone[record.Topic])
						}
						tool.SafeSend(work.ctx, work.ch, record)
					}
					for _, work := range workerCtxMap {
						close(work.ch)
					}
					wg.Wait()
					clear(workerCtxMap)
				}
			}()
			time.Sleep(3 * time.Second)
			slog.Info("kafka service restart", slog.String("name", self.name))
		}
	}()
}

// worker processes messages for a specific topic, supporting batch or single-message handling.
// It invokes the consumer's Handler for each batch or message.
func (self *Service) worker(ctx *workerCtx, wg *sync.WaitGroup, consumer *Consumer) {
	defer func() {
		wg.Done()
		ctx.cancel()
	}()
	batchRecord := make([]*kgo.Record, 0)
	var lastTime int64
Loop:
	for {
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
			batchRecord = append(batchRecord, record)
			now := time.Now().UnixMilli()
			if len(batchRecord) >= consumer.BatchMaxCount || now-lastTime > 3000 {
				consumer.Handler(ctx.ctx, self.client, self.autoCommit, batchRecord...)
				lastTime = now
				batchRecord = batchRecord[:0]
			}
		}
	}
	if len(batchRecord) > 0 {
		consumer.Handler(ctx.ctx, self.client, self.autoCommit, batchRecord...)
	}
}
