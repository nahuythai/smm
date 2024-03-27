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

type Service interface {
	CreateOne(service models.Service) (*models.Service, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Service, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdateById(id primitive.ObjectID, doc ServiceUpdateByIdDoc) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (service *models.Service, err error)
	DeleteById(id primitive.ObjectID) error
	GetActiveById(id primitive.ObjectID, opts ...QueryOption) (*models.Service, error)
	GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.Service, error)
	GetActiveBySeq(seq int, opts ...QueryOption) (*models.Service, error)
}
type serviceQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewService(ctx context.Context) Service {
	return &serviceQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.ServiceCollectionName),
	}
}

func (q *serviceQuery) CreateOne(service models.Service) (*models.Service, error) {
	now := time.Now()
	service.UpdatedAt = now
	service.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, service)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("serviceQuery")
		return nil, err
	}
	service.Id = result.InsertedID.(primitive.ObjectID)
	return &service, nil
}

func (q *serviceQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Service, error) {
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
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("serviceQuery")
		return nil, err
	}
	var services []models.Service
	if err = cursor.All(q.ctx, &services); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("serviceQuery")
		return nil, err
	}
	return services, nil
}

func (q *serviceQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("serviceQuery")
		return 0, err
	}
	return total, nil
}

func (q *serviceQuery) UpdateById(id primitive.ObjectID, doc ServiceUpdateByIdDoc) error {
	doc.UpdatedAt = time.Now()
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": doc,
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("serviceQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeServiceNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *serviceQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.Service, error) {
	var service models.Service
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&service); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeServiceNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("serviceQuery")
		return nil, err
	}
	return &service, nil
}

func (q *serviceQuery) DeleteById(id primitive.ObjectID) error {
	result, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("serviceQuery")
		return err
	}
	if result.DeletedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeServiceNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *serviceQuery) GetActiveById(id primitive.ObjectID, opts ...QueryOption) (*models.Service, error) {
	var service models.Service
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id, "status": constants.ServiceStatusOn}, &findOpt).Decode(&service); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeServiceNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetActiveById").Str("funcInline", "q.collection.FindOne").Msg("serviceQuery")
		return nil, err
	}
	return &service, nil
}

func (q *serviceQuery) GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.Service, error) {
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpts := options.FindOptions{
		Projection: opt.QueryOnlyField(),
	}
	cursor, err := q.collection.Find(q.ctx, bson.M{"_id": bson.M{"$in": ids}}, &findOpts)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "q.collection.Find").Msg("serviceQuery")
		return nil, err
	}
	var services []models.Service
	if err = cursor.All(q.ctx, &services); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "cursor.All").Msg("serviceQuery")
		return nil, err
	}
	return services, nil
}

func (q *serviceQuery) GetActiveBySeq(seq int, opts ...QueryOption) (*models.Service, error) {
	var service models.Service
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"seq": seq, "status": constants.ServiceStatusOn}, &findOpt).Decode(&service); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeServiceNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetActiveBySeq").Str("funcInline", "q.collection.FindOne").Msg("serviceQuery")
		return nil, err
	}
	return &service, nil
}
