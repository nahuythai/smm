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

type Otp interface {
	CreateOne(otp models.Otp) (*models.Otp, error)
	GetByUserId(userId primitive.ObjectID, opts ...QueryOption) (otp *models.Otp, err error)
}
type otpQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewOtp(ctx context.Context) Otp {
	return &otpQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.OtpCollectionName),
	}
}

func (q *otpQuery) CreateOne(otp models.Otp) (*models.Otp, error) {
	now := time.Now()
	otp.UpdatedAt = now
	otp.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, otp)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, response.NewError(fiber.StatusConflict, response.Option{Code: constants.ErrCodeUserExist})
		}
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("otpQuery")
		return nil, err
	}
	otp.Id = result.InsertedID.(primitive.ObjectID)
	return &otp, nil
}

func (q *otpQuery) GetByUserId(userId primitive.ObjectID, opts ...QueryOption) (*models.Otp, error) {
	var otp models.Otp
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"user_id": userId}, &findOpt).Decode(&otp); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeOtpNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetByUserId").Str("funcInline", "q.collection.FindOne").Msg("otpQuery")
		return nil, err
	}
	return &otp, nil
}

func (q *otpQuery) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)}}
	if _, err := database.DB.Collection(models.UserCollectionName).Indexes().CreateMany(context.Background(), indexes); err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateIndexes").Str("funcInline", "Indexes().CreateMany").Msg("otpQuery")
		return err
	}
	return nil
}
