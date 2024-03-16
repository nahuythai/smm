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

type User interface {
	CreateOne(user models.User) (*models.User, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.User, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdateApiKeyById(id primitive.ObjectID, apiKey string) error
	UpdateById(id primitive.ObjectID, doc UserUpdateByIdDoc) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (user *models.User, err error)
	GetByUsername(username string, opts ...QueryOption) (user *models.User, err error)
	UpdateBalanceById(id primitive.ObjectID, balance float64) error
	UpdatePasswordById(id primitive.ObjectID, password string) error
	CreateIndexes() error
	UpdateEmailVerificationById(id primitive.ObjectID, verified bool) error
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
		if mongo.IsDuplicateKeyError(err) {
			return nil, response.NewError(fiber.StatusConflict, response.Option{Code: constants.ErrCodeUserExist})
		}
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("userQuery")
		return nil, err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (q *userQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.User, error) {
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
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("userQuery")
		return nil, err
	}
	var users []models.User
	if err = cursor.All(q.ctx, &users); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("userQuery")
		return nil, err
	}
	return users, nil
}

func (q *userQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Caller().Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("userQuery")
		return 0, err
	}
	return total, nil
}

func (q *userQuery) UpdateApiKeyById(id primitive.ObjectID, apiKey string) error {
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": bson.M{"api_key": apiKey, "updated_at": time.Now()},
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateApiKeyById").Str("funcInline", "q.collection.UpdateByID").Msg("userQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *userQuery) UpdateById(id primitive.ObjectID, doc UserUpdateByIdDoc) error {
	doc.UpdatedAt = time.Now()
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": doc,
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("userQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *userQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.User, error) {
	var user models.User
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("userQuery")
		return nil, err
	}
	return &user, nil
}

func (q *userQuery) UpdateBalanceById(id primitive.ObjectID, balance float64) error {
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": bson.M{"balance": balance, "updated_at": time.Now()},
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateBalanceById").Str("funcInline", "q.collection.UpdateByID").Msg("userQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *userQuery) UpdatePasswordById(id primitive.ObjectID, password string) error {
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": bson.M{"password": password, "updated_at": time.Now()},
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdatePasswordById").Str("funcInline", "q.collection.UpdateByID").Msg("userQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *userQuery) GetByUsername(username string, opts ...QueryOption) (*models.User, error) {
	var user models.User
	opt := NewOption()
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"username": username}, &findOpt).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Caller().Str("func", "GetByUsername").Str("funcInline", "q.collection.FindOne").Msg("userQuery")
		return nil, err
	}
	return &user, nil
}

func (q *userQuery) CreateIndexes() error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)}}
	if _, err := database.DB.Collection(models.UserCollectionName).Indexes().CreateMany(context.Background(), indexes); err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateIndexes").Str("funcInline", "Indexes().CreateMany").Msg("userQuery")
		return err
	}
	return nil
}

func (q *userQuery) UpdateEmailVerificationById(id primitive.ObjectID, verified bool) error {
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": bson.M{"email_verification": verified, "updated_at": time.Now()},
	})
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UpdateEmailVerificationById").Str("funcInline", "q.collection.UpdateByID").Msg("userQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.Option{Code: constants.ErrCodeUserNotFound, Data: constants.ErrMsgResourceNotFound})
	}
	return nil
}
