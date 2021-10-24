package sms

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

/*
{
	"from_number": "01065146909",
	"contact_name": "Cellularhacker@DEXEOS",
	"text": "Urlurl",
	"occurred_at": "January 19, 2020 at 09:54PM"
}
*/
type Sms struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	FromNumber  string        `bson:"from_number" json:"from_number"`
	ContactName string        `bson:"contact_name" json:"contact_name"`
	Text        string        `bson:"text" json:"text"`
	OccurredAt  string        `bson:"occurred_at" json:"occurred_at"`
	ToNumber    string        `bson:"to_number" json:"to_number"`
	ReceivedAt  int64         `bson:"received_at" json:"received_at"`

	TextType    int    `bson:"text_type" json:"text_type"`
	MessageHash string `bson:"message_hash" json:"message_hash"`
	CreatedAt   int64  `bson:"created_at" json:"created_at"`
}

func New() *Sms {
	return &Sms{ID: bson.NewObjectId()}
}

func (s *Sms) Create() error {
	if IsExistHash(s.MessageHash) {
		return nil
	}

	s.CreatedAt = time.Now().Unix()
	return store.Create(s)
}
