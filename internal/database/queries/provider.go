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

type Provider interface {
	CreateOne(provider models.Provider) (*models.Provider, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Provider, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdateById(id primitive.ObjectID, doc ProviderUpdateByIdDoc) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (provider *models.Provider, err error)
	DeleteById(id primitive.ObjectID) error
	GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.Provider, error)
}
type providerQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewProvider(ctx context.Context) Provider {
	return &providerQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.ProviderCollectionName),
	}
}

func (q *providerQuery) CreateOne(provider models.Provider) (*models.Provider, error) {
	now := time.Now()
	provider.UpdatedAt = now
	provider.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, provider)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("providerQuery")
		return nil, err
	}
	provider.Id = result.InsertedID.(primitive.ObjectID)
	return &provider, nil
}

func (q *providerQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Provider, error) {
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpts := options.FindOptions{
		Skip:       opt.QuerySkip(),
		Limit:      opt.QueryLimit(),
		Projection: opt.QueryOnlyField(),
		Sort:       opt.QuerySort(),
	}
	cursor, err := q.collection.Find(q.ctx, filter, &findOpts)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("providerQuery")
		return nil, err
	}
	var providers []models.Provider
	if err = cursor.All(q.ctx, &providers); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("providerQuery")
		return nil, err
	}
	return providers, nil
}

func (q *providerQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("providerQuery")
		return 0, err
	}
	return total, nil
}

func (q *providerQuery) UpdateById(id primitive.ObjectID, doc ProviderUpdateByIdDoc) error {
	doc.UpdatedAt = time.Now()
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": doc,
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("providerQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeProviderNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *providerQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.Provider, error) {
	var provider models.Provider
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&provider); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("providerQuery")
		return nil, err
	}
	return &provider, nil
}

func (q *providerQuery) DeleteById(id primitive.ObjectID) error {
	result, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("providerQuery")
		return err
	}
	if result.DeletedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *providerQuery) GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.Provider, error) {
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpts := options.FindOptions{
		Projection: opt.QueryOnlyField(),
	}
	cursor, err := q.collection.Find(q.ctx, bson.M{"_id": bson.M{"$in": ids}}, &findOpts)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "q.collection.Find").Msg("providerQuery")
		return nil, err
	}
	var providers []models.Provider
	if err = cursor.All(q.ctx, &providers); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "cursor.All").Msg("providerQuery")
		return nil, err
	}
	return providers, nil
}
