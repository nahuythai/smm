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

type Transaction interface {
	CreateOne(transaction models.Transaction) (*models.Transaction, error)
	GetById(id primitive.ObjectID, opts ...QueryOption) (transaction *models.Transaction, err error)
	DeleteById(id primitive.ObjectID) error
}
type transactionQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewTransaction(ctx context.Context) Transaction {
	return &transactionQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.TransactionCollectionName),
	}
}

func (q *transactionQuery) CreateOne(transaction models.Transaction) (*models.Transaction, error) {
	now := time.Now()
	transaction.UpdatedAt = now
	transaction.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, transaction)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("transactionQuery")
		return nil, err
	}
	transaction.Id = result.InsertedID.(primitive.ObjectID)
	return &transaction, nil
}

func (q *transactionQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.Transaction, error) {
	var transaction models.Transaction
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&transaction); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeTransactionNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("transactionQuery")
		return nil, err
	}
	return &transaction, nil
}

func (q *transactionQuery) DeleteById(id primitive.ObjectID) error {
	if _, err := q.collection.DeleteOne(q.ctx, bson.M{"_id": id}); err != nil {
		logger.Error().Err(err).Caller().Str("func", "DeleteById").Str("funcInline", "q.collection.DeleteOne").Msg("transactionQuery")
		return err
	}
	return nil
}

func (q *transactionQuery) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "expired_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)}}
	if _, err := database.DB.Collection(models.UserCollectionName).Indexes().CreateMany(context.Background(), indexes); err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateIndexes").Str("funcInline", "Indexes().CreateMany").Msg("transactionQuery")
		return err
	}
	return nil
}
