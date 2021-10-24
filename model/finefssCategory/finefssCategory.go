package finefssCategory

import (
	"github.com/globalsign/mgo/bson"
)

type FineFssCategory struct {
	ID        bson.ObjectId `bson:"_id" json:"_id"`
	Value     string        `bson:"value" json:"value"`
	Text      string        `bson:"text" json:"text"`
	CreatedAt int64         `bson:"created_at" json:"created_at"`
}

func New() *FineFssCategory {
	return &FineFssCategory{ID: bson.NewObjectId()}
}

func (f *FineFssCategory) Create() error {
	if f.Text == "" {
		return nil
	}
	return store.Create(f)
}
