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

type CustomRate interface {
	CreateOne(customRate models.CustomRate) (*models.CustomRate, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.CustomRate, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdatePriceById(id primitive.ObjectID, price float64) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (customRate *models.CustomRate, err error)
	DeleteById(id primitive.ObjectID) error
	GetByUserIdAndServiceId(userId primitive.ObjectID, serviceId primitive.ObjectID, opts ...QueryOption) (*models.CustomRate, error)
}
type customRateQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewCustomRate(ctx context.Context) CustomRate {
	return &customRateQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.CustomRateCollectionName),
	}
}

func (q *customRateQuery) CreateOne(customRate models.CustomRate) (*models.CustomRate, error) {
	now := time.Now()
	customRate.UpdatedAt = now
	customRate.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, customRate)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("customRateQuery")
		return nil, err
	}
	customRate.Id = result.InsertedID.(primitive.ObjectID)
	return &customRate, nil
}

func (q *customRateQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.CustomRate, error) {
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
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("customRateQuery")
		return nil, err
	}
	var categories []models.CustomRate
	if err = cursor.All(q.ctx, &categories); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("customRateQuery")
		return nil, err
	}
	return categories, nil
}

func (q *customRateQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("customRateQuery")
		return 0, err
	}
	return total, nil
}

func (q *customRateQuery) UpdatePriceById(id primitive.ObjectID, price float64) error {
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": bson.M{"price": price, "updated_at": time.Now()},
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("customRateQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeCustomRateNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *customRateQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.CustomRate, error) {
	var customRate models.CustomRate
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&customRate); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeCustomRateNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("customRateQuery")
		return nil, err
	}
	return &customRate, nil
}

func (q *customRateQuery) GetByUserIdAndServiceId(userId primitive.ObjectID, serviceId primitive.ObjectID, opts ...QueryOption) (*models.CustomRate, error) {
	var customRate models.CustomRate
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"user_id": userId, "service_id": serviceId}, &findOpt).Decode(&customRate); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeCustomRateNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetByUserIdAndServiceId").Str("funcInline", "q.collection.FindOne").Msg("customRateQuery")
		return nil, err
	}
	return &customRate, nil
}

func (q *customRateQuery) DeleteById(id primitive.ObjectID) error {
	result, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("customRateQuery")
		return err
	}
	if result.DeletedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeCustomRateNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}
