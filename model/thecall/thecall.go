package thecall

import (
	"github.com/globalsign/mgo/bson"
)

type TheCall struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	PhoneNumber string        `bson:"phone_number" json:"phone_number"`
	Subject     string        `bson:"subject" json:"subject"`
	IsWhiteList bool          `bson:"is_white_list" json:"is_white_list"`
	//AddedMember string        `bson:"added_member" json:"added_member"`
	//View        int           `bson:"view" json:"view"`
	//AddedAt     int64         `bson:"added_at"`

	CreatedAt int64 `bson:"created_at" json:"created_at"`
	UpdatedAt int64 `bson:"updated_at" json:"updated_at"`

	//AskCount int `bson:"ask_count" json:"ask_count"`
}

//type Comment struct {
//	Nickname  string `bson:"nickname" json:"nickname"`
//	Content   string `bson:"content" json:"content"`
//	At        int64  `bson:"at" json:"at"`
//	Recommend int    `bson:"recommend" json:"recommend"`
//}

func New() *TheCall {
	return &TheCall{ID: bson.NewObjectId()}
}

func (c *TheCall) Upsert() error {
	ps, err := GetByPhoneNumber(c.PhoneNumber)
	if err != nil {
		return nil
	}

	// Create a new.
	if ps == nil {
		return store.Create(c)
	}

	// Check Update
	if ps.Subject != c.Subject && ps.IsWhiteList == c.IsWhiteList {
		return store.UpdateSet(bson.M{KeyID: ps.ID}, bson.M{KeySubject: c.Subject})
	}

	// Nothing to update.
	return nil
}
