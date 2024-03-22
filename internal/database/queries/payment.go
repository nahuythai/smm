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

type Payment interface {
	CreateOne(payment models.Payment) (*models.Payment, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Payment, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdateById(id primitive.ObjectID, doc PaymentUpdateByIdDoc) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (payment *models.Payment, err error)
	DeleteById(id primitive.ObjectID) error
	GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.Payment, error)
}
type paymentQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewPayment(ctx context.Context) Payment {
	return &paymentQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.PaymentCollectionName),
	}
}

func (q *paymentQuery) CreateOne(payment models.Payment) (*models.Payment, error) {
	now := time.Now()
	payment.UpdatedAt = now
	payment.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, payment)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("paymentQuery")
		return nil, err
	}
	payment.Id = result.InsertedID.(primitive.ObjectID)
	return &payment, nil
}

func (q *paymentQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Payment, error) {
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
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("paymentQuery")
		return nil, err
	}
	var payments []models.Payment
	if err = cursor.All(q.ctx, &payments); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("paymentQuery")
		return nil, err
	}
	return payments, nil
}

func (q *paymentQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("paymentQuery")
		return 0, err
	}
	return total, nil
}

func (q *paymentQuery) UpdateById(id primitive.ObjectID, doc PaymentUpdateByIdDoc) error {
	doc.UpdatedAt = time.Now()
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": doc,
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("paymentQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodePaymentNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *paymentQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.Payment, error) {
	var payment models.Payment
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&payment); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("paymentQuery")
		return nil, err
	}
	return &payment, nil
}

func (q *paymentQuery) DeleteById(id primitive.ObjectID) error {
	result, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("paymentQuery")
		return err
	}
	if result.DeletedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *paymentQuery) GetByIds(ids []primitive.ObjectID, opts ...QueryOption) ([]models.Payment, error) {
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpts := options.FindOptions{
		Projection: opt.QueryOnlyField(),
	}
	cursor, err := q.collection.Find(q.ctx, bson.M{"_id": bson.M{"$in": ids}}, &findOpts)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "q.collection.Find").Msg("paymentQuery")
		return nil, err
	}
	var payments []models.Payment
	if err = cursor.All(q.ctx, &payments); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByIds").Str("funcInline", "cursor.All").Msg("paymentQuery")
		return nil, err
	}
	return payments, nil
}
