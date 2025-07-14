package kafka

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wnnce/fserv-template/config"
)

// Service manages the Kafka client, consumers, and context for message processing.
type Service struct {
	runner     atomic.Bool
	name       string
	client     *kgo.Client
	autoCommit bool
	mutex      *sync.Mutex
	consumers  map[string]*Consumer
	ctx        context.Context
	cancel     context.CancelFunc
	once       sync.Once
}

// NewService creates a new Service instance with the given name, Kafka client, and auto-commit setting.
func NewService(name string, client *kgo.Client, autoCommit bool) *Service {
	return &Service{
		name:       name,
		client:     client,
		autoCommit: autoCommit,
		mutex:      &sync.Mutex{},
		consumers:  make(map[string]*Consumer),
	}
}

// Client returns the underlying kgo.Client instance.
func (self *Service) Client() *kgo.Client {
	return self.client
}

// Shutdown gracefully closes the Kafka client and cancels all consumers.
func (self *Service) Shutdown() {
	self.once.Do(func() {
		clear(self.consumers)
		self.client.Close()
		self.cancel()
	})
}

var (
	defaultService *Service
)

// InitKafkaService initializes the global Kafka Service using configuration and returns a cleanup function.
func InitKafkaService(ctx context.Context) (func(), error) {
	brokers := viper.GetStringSlice("kafka.brokers")
	autoCommit := config.ViperGet[bool]("kafka.commit", true)
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ClientID(config.ViperGet[string]("kafka.client-id")),
		kgo.ConsumerGroup(config.ViperGet[string]("kafka.consumer-group")),
		func(commit bool) kgo.Opt {
			if !commit {
				return kgo.DisableAutoCommit()
			}
			return kgo.GreedyAutoCommit()
		}(autoCommit),
	)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx); err != nil {
		return nil, err
	}
	childCtx, cancel := context.WithCancel(ctx)
	defaultService = &Service{
		client:     client,
		autoCommit: autoCommit,
		mutex:      &sync.Mutex{},
		consumers:  make(map[string]*Consumer),
		ctx:        childCtx,
		cancel:     cancel,
	}
	return func() {
		defaultService.Shutdown()
	}, nil
}

// Instance returns the default global Kafka Service instance.
func Instance() *Service {
	return defaultService
}
