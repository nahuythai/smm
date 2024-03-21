package queries

import (
	"context"
	"smm/internal/database"
	"smm/internal/database/models"
	"smm/pkg/constants"
	"smm/pkg/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	CreateOne(session models.Session) (*models.Session, error)
	GetById(id primitive.ObjectID, opts ...QueryOption) (session *models.Session, err error)
	DeleteById(id primitive.ObjectID) error
}
type sessionQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewSession(ctx context.Context) Session {
	return &sessionQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.SessionCollectionName),
	}
}

func (q *sessionQuery) CreateOne(session models.Session) (*models.Session, error) {
	now := time.Now()
	session.UpdatedAt = now
	session.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, session)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("sessionQuery")
		return nil, err
	}
	session.Id = result.InsertedID.(primitive.ObjectID)
	return &session, nil
}

func (q *sessionQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.Session, error) {
	var session models.Session
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&session); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeSessionNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("sessionQuery")
		return nil, err
	}
	return &session, nil
}

func (q *sessionQuery) DeleteById(id primitive.ObjectID) error {
	if _, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id}); err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("sessionQuery")
		return err
	}
	return nil
}

func (q *sessionQuery) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "expired_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)}}
	if _, err := database.DB.Collection(models.UserCollectionName).Indexes().CreateMany(context.Background(), indexes); err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateIndexes").Str("funcInline", "Indexes().CreateMany").Msg("sessionQuery")
		return err
	}
	return nil
}
