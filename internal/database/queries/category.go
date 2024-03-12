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

type Category interface {
	CreateOne(category models.Category) (*models.Category, error)
	GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Category, error)
	GetTotalByFilter(filter bson.M) (total int64, err error)
	UpdateById(id primitive.ObjectID, doc CategoryUpdateByIdDoc) error
	GetById(id primitive.ObjectID, opts ...QueryOption) (category *models.Category, err error)
}
type categoryQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewCategory(ctx context.Context) Category {
	return &categoryQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.CategoryCollectionName),
	}
}

func (q *categoryQuery) CreateOne(category models.Category) (*models.Category, error) {
	now := time.Now()
	category.UpdatedAt = now
	category.CreatedAt = now
	result, err := q.collection.InsertOne(q.ctx, category)
	if err != nil {
		logger.Error().Err(err).Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("categoryQuery")
		return nil, err
	}
	category.Id = result.InsertedID.(primitive.ObjectID)
	return &category, nil
}

func (q *categoryQuery) GetByFilter(filter bson.M, opts ...QueryOption) ([]models.Category, error) {
	var opt QueryOption
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
		logger.Error().Err(err).Str("func", "GetByFilter").Str("funcInline", "q.collection.Find").Msg("categoryQuery")
		return nil, err
	}
	var categories []models.Category
	if err = cursor.All(q.ctx, &categories); err != nil {
		logger.Error().Err(err).Str("func", "GetByFilter").Str("funcInline", "cursor.All").Msg("categoryQuery")
		return nil, err
	}
	return categories, nil
}

func (q *categoryQuery) GetTotalByFilter(filter bson.M) (total int64, err error) {
	if total, err = q.collection.CountDocuments(q.ctx, filter); err != nil {
		logger.Error().Err(err).Str("func", "GetTotalByFilter").Str("funcInline", "q.collection.CountDocuments").Msg("categoryQuery")
		return 0, err
	}
	return total, nil
}

func (q *categoryQuery) UpdateById(id primitive.ObjectID, doc CategoryUpdateByIdDoc) error {
	doc.UpdatedAt = time.Now()
	result, err := q.collection.UpdateByID(q.ctx, id, bson.M{
		"$set": doc,
	})
	if err != nil {
		logger.Error().Err(err).Str("func", "UpdateById").Str("funcInline", "q.collection.UpdateByID").Msg("categoryQuery")
		return err
	}
	if result.MatchedCount == 0 {
		return response.NewError(fiber.StatusNotFound, response.ErrorResponse{Code: constants.ErrCodeCategoryNotFound, Err: constants.ErrMsgResourceNotFound})
	}
	return nil
}

func (q *categoryQuery) GetById(id primitive.ObjectID, opts ...QueryOption) (*models.Category, error) {
	var category models.Category
	var opt QueryOption
	if len(opts) > 0 {
		opt = opts[0]
	}
	findOpt := options.FindOneOptions{
		Projection: opt.QueryOnlyField(),
	}
	if err := q.collection.FindOne(q.ctx, bson.M{"_id": id}, &findOpt).Decode(&category); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, response.NewError(fiber.StatusNotFound, response.ErrorResponse{Code: constants.ErrCodeUserNotFound, Err: constants.ErrMsgResourceNotFound})
		}
		logger.Error().Err(err).Str("func", "GetById").Str("funcInline", "q.collection.FindOne").Msg("categoryQuery")
		return nil, err
	}
	return &category, nil
}
