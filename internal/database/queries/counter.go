package queries

import (
	"context"
	"smm/internal/database"
	"smm/internal/database/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Counter interface {
	CreateOne(counter models.Counter) (*models.Counter, error)
	GenerateSeq(collection string) (int, error)
}
type counterQuery struct {
	ctx        context.Context
	collection *mongo.Collection
}

func NewCounter(ctx context.Context) Counter {
	return &counterQuery{
		ctx:        ctx,
		collection: database.DB.Collection(models.CounterCollectionName),
	}
}

func (q *counterQuery) CreateOne(counter models.Counter) (*models.Counter, error) {
	result, err := q.collection.InsertOne(q.ctx, counter)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "CreateOne").Str("funcInline", "q.collection.InsertOne").Msg("counterQuery")
		return nil, err
	}
	counter.Id = result.InsertedID.(primitive.ObjectID)
	return &counter, nil
}

func (q *counterQuery) GenerateSeq(collection string) (int, error) {
	after := options.After
	var res models.Counter
	if err := q.collection.FindOneAndUpdate(q.ctx,
		bson.M{"collection": collection}, bson.M{
			"$inc": bson.M{"seq": 1},
		},
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}).Decode(&res); err != nil {
		return 0, err
	}
	return res.Seq, nil
}
