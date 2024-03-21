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

type Order interface {
	CreateOne(order models.Order) (*models.Order, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Order, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdateById(id primitive.ObjectID, doc OrderUpdateByIdDoc) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (order *models.Order, err error)
	DeleteById(id primitive.ObjectID) error
}
type orderQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewOrder(ctx context.Context) Order {
	return &orderQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.OrderCollectionName),
	}
}

func (q *orderQuery) CreateOne(order models.Order) (*models.Order, error) {
	now := time.Now()
	order.UpdatedAt = now
	order.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, order)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("orderQuery")
		return nil, err
	}
	order.Id = result.InsertedID.(primitive.ObjectID)
	return &order, nil
}

func (q *orderQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Order, error) {
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
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("orderQuery")
		return nil, err
	}
	var categories []models.Order
	if err = cursor.All(q.ctx, &categories); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("orderQuery")
		return nil, err
	}
	return categories, nil
}

func (q *orderQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("orderQuery")
		return 0, err
	}
	return total, nil
}

func (q *orderQuery) UpdateById(id primitive.ObjectID, doc OrderUpdateByIdDoc) error {
	doc.UpdatedAt = time.Now()
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": doc,
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("orderQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeOrderNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *orderQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.Order, error) {
	var order models.Order
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&order); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("orderQuery")
		return nil, err
	}
	return &order, nil
}

func (q *orderQuery) DeleteById(id primitive.ObjectID) error {
	result, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("orderQuery")
		return err
	}
	if result.DeletedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}
