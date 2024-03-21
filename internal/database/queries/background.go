package queries

import (
	"context"
	"smm/internal/database"
	"smm/internal/database/models"
	"smm/pkg/configure"
	"smm/pkg/constants"
	"smm/pkg/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	cfg = configure.GetConfig()
)

type Background interface {
	CreateOne(background models.Background) (*models.Background, error)
	DeleteById(id primitive.ObjectID) error
	CreateIndexes() error
}
type backgroundQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewBackground(ctx context.Context) Background {
	return &backgroundQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.BackgroundCollectionName),
	}
}

func (q *backgroundQuery) CreateOne(background models.Background) (*models.Background, error) {
	now := time.Now()
	background.UpdatedAt = now
	background.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, background)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, response.NewError(fiber.StatusConflict, response.Option{Code: constants.ErrBackgroundTaskExist})
		}
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("backgroundQuery")
		return nil, err
	}
	background.Id = result.InsertedID.(primitive.ObjectID)
	return &background, nil
}

func (q *backgroundQuery) DeleteById(id primitive.ObjectID) error {
	if _, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id}); err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("backgroundQuery")
		return err
	}
	return nil
}
func (q *backgroundQuery) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "type", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "type", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(int32(cfg.BackgroundTaskDuration.Seconds()))}}
	if _, err := database.DB.Collection(models.UserCollectionName).Indexes().CreateMany(context.Background(), indexes); err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateIndexes").Str("funcInline", "Indexes().CreateMany").Msg("backgroundQuery")
		return err
	}
	return nil
}
