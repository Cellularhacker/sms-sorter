package finefssCategory

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FineFssCategory struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Value     string             `bson:"value" json:"value"`
	Text      string             `bson:"text" json:"text"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
}

func New() *FineFssCategory {
	return &FineFssCategory{}
}

func (f *FineFssCategory) Create() error {
	if f.Text == "" {
		return nil
	}
	return Create(f)
}
