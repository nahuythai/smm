package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Counter struct {
	Seq        int                `bson:"seq"`
	Collection string             `bson:"collection"`
	Id         primitive.ObjectID `bson:"_id,omitempty"`
}

const CounterCollectionName = "counters"
