package database

import (
	"context"
	"smm/pkg/configure"
	"smm/pkg/logging"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Database
	logger = logging.GetLogger()
	cfg    = configure.GetConfig()
)

func GetTimeoutWithContext(ctx context.Context) (c context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(ctx, cfg.MongoRequestTimeout)
}

func InitDatabase() {
	ctx, cancel := GetTimeoutWithContext(context.Background())
	defer cancel()
	// Create a MongoDB client
	logLevel := options.LogLevelInfo
	if cfg.Debug {
		logLevel = options.LogLevelDebug
	}
	loggerOptions := options.Logger().SetComponentLevel(options.LogComponentCommand, logLevel)
	clientOptions := options.Client().ApplyURI(cfg.MongoDBUrl).SetLoggerOptions(loggerOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal().Err(err).Str("func", "InitDatabase").Str("funcInline", "mongo.Connect").Msg("database")
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal().Err(err).Str("func", "InitDatabase").Str("funcInline", "client.Ping").Msg("database")
	}
	DB = client.Database(cfg.MongoDBName)
}
