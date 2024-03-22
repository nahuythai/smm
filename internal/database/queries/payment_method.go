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

type PaymentMethod interface {
	CreateOne(paymentMethod models.PaymentMethod) (*models.PaymentMethod, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.PaymentMethod, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdateById(id primitive.ObjectID, doc PaymentMethodUpdateByIdDoc) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (paymentMethod *models.PaymentMethod, err error)
	DeleteById(id primitive.ObjectID) error
	GetActiveById(id primitive.ObjectID, opts ...QueryOption) (*models.PaymentMethod, error)
	GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.PaymentMethod, error)
}
type paymentMethodQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewPaymentMethod(ctx context.Context) PaymentMethod {
	return &paymentMethodQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.PaymentMethodCollectionName),
	}
}

func (q *paymentMethodQuery) CreateOne(paymentMethod models.PaymentMethod) (*models.PaymentMethod, error) {
	now := time.Now()
	paymentMethod.UpdatedAt = now
	paymentMethod.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, paymentMethod)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("paymentMethodQuery")
		return nil, err
	}
	paymentMethod.Id = result.InsertedID.(primitive.ObjectID)
	return &paymentMethod, nil
}

func (q *paymentMethodQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.PaymentMethod, error) {
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
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("paymentMethodQuery")
		return nil, err
	}
	var paymentMethods []models.PaymentMethod
	if err = cursor.All(q.ctx, &paymentMethods); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("paymentMethodQuery")
		return nil, err
	}
	return paymentMethods, nil
}

func (q *paymentMethodQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("paymentMethodQuery")
		return 0, err
	}
	return total, nil
}

func (q *paymentMethodQuery) UpdateById(id primitive.ObjectID, doc PaymentMethodUpdateByIdDoc) error {
	doc.UpdatedAt = time.Now()
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": doc,
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("paymentMethodQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodePaymentMethodNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *paymentMethodQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&paymentMethod); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("paymentMethodQuery")
		return nil, err
	}
	return &paymentMethod, nil
}

func (q *paymentMethodQuery) DeleteById(id primitive.ObjectID) error {
	result, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("paymentMethodQuery")
		return err
	}
	if result.DeletedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *paymentMethodQuery) GetActiveById(id primitive.ObjectID, opts ...QueryOption) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id, "status": constants.PaymentMethodStatusOn}, &findOpt).Decode(&paymentMethod); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetActiveById").Str("funcInline", "q.collection.FindOne").Msg("paymentMethodQuery")
		return nil, err
	}
	return &paymentMethod, nil
}

func (q *paymentMethodQuery) GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.PaymentMethod, error) {
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpts := options.FindOptions{
		Projection: opt.QueryOnlyField(),
	}
	cursor, err := q.collection.Find(q.ctx, bson.M{"_id": bson.M{"$in": ids}}, &findOpts)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "q.collection.Find").Msg("paymentMethodQuery")
		return nil, err
	}
	var paymentMethods []models.PaymentMethod
	if err = cursor.All(q.ctx, &paymentMethods); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "cursor.All").Msg("paymentMethodQuery")
		return nil, err
	}
	return paymentMethods, nil
}
