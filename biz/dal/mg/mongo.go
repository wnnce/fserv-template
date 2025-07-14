package mg

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/wnnce/fserv-template/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	mongoClient   *mongo.Client                // mongoClient holds the MongoDB client instance.
	mongoDB       *mongo.Database              // mongoDB references the connected MongoDB database.
	collectionMap map[string]*mongo.Collection // collectionMap caches collection instances by name.
	mutex         sync.Mutex                   // mutex ensures safe concurrent access to collectionMap.
)

// InitMongoDB initializes a MongoDB client using configuration loaded from Viper.
// It validates the connection with a ping and caches the connected database.
//
// Returns a cleanup function to disconnect the client and clear internal caches,
// or an error if the connection or ping fails.
func InitMongoDB(ctx context.Context) (func(), error) {
	mongoURL := config.ViperGet[string]("mongo.url", "mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = client.Ping(timeout, readpref.Primary()); err != nil {
		slog.Error("mongoDB ping failed", slog.String("error", err.Error()))
		_ = client.Disconnect(ctx)
		return nil, err
	}
	slog.Info("mongoDB connection success", slog.Group("data", slog.String("url", mongoURL)))
	mongoClient = client
	mongoDB = client.Database(config.ViperGet[string]("mongo.database"))
	collectionMap = make(map[string]*mongo.Collection)
	return func() {
		_ = mongoClient.Disconnect(ctx)
		clear(collectionMap)
	}, nil
}

// Database returns the connected MongoDB database instance.
// Panics if the database has not been initialized.
func Database() *mongo.Database {
	return mongoDB
}

// Collection returns a cached MongoDB collection by name.
// If the collection is not cached yet, it creates and stores it.
//
// This function is safe for concurrent use.
func Collection(collection string) *mongo.Collection {
	mutex.Lock()
	defer mutex.Unlock()
	value, ok := collectionMap[collection]
	if !ok {
		value = mongoDB.Collection(collection)
		collectionMap[collection] = value
	}
	return value
}
