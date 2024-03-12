package queries

import (
	"context"
	"smm/internal/database"
	"smm/internal/database/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User interface {
	CreateOne(user models.User) (*models.User, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.User, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
}
type userQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewUser(ctx context.Context) User {
	return &userQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.UserCollectionName),
	}
}

func (q *userQuery) CreateOne(user models.User) (*models.User, error) {
	now := time.Now()
	user.UpdatedAt = now
	user.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, user)
	if err != nil {
		logger.Error().Err(err).Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("userQuery")
		return nil, err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (q *userQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.User, error) {
	var opt QueryOption
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpts := options.FindOptions{
		Skip:       opt.QuerySkip(),
		Limit:      opt.QueryLimit(),
		Projection: opt.QueryOnlyField(),
	}
	cursor, err := q.collection.Find(q.ctx, filter, &findOpts)
	if err != nil {
		logger.Error().Err(err).Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("userQuery")
		return nil, err
	}
	var users []models.User
	if err = cursor.All(q.ctx, &users); err != nil {
		logger.Error().Err(err).Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("userQuery")
		return nil, err
	}
	return users, nil
}

func (q *userQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("userQuery")
		return 0, err
	}
	return total, nil
}
