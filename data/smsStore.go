package data

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"sms-sorter/model/sms"
	"time"
)

type smsStore struct {
	collection *mgo.Collection
	context    *Context
}

func (hs *smsStore) setContext() {
	hs.context = NewTokyoContext()
	hs.collection = hs.context.SmsDBCollection(CSms)
}

func NewSmsStore() *smsStore {
	return &smsStore{}
}

func (hs smsStore) Create(h *sms.Sms) error {
	hs.setContext()
	defer hs.context.Close()

	h.ID = bson.NewObjectId()
	h.CreatedAt = time.Now().Unix()

	return hs.collection.Insert(h)
}

func (hs smsStore) GetBy(query bson.M) ([]sms.Sms, error) {
	hs.setContext()
	defer hs.context.Close()

	results := make([]sms.Sms, 0)
	err := hs.collection.Find(query).All(&results)

	return results, err
}

func (hs smsStore) GetByDesc(query bson.M) ([]sms.Sms, error) {
	hs.setContext()
	defer hs.context.Close()

	results := make([]sms.Sms, 0)
	err := hs.collection.Find(query).Sort("-created_at").All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (hs smsStore) GetBySortLimit(query bson.M, sort string, limit int) ([]sms.Sms, error) {
	hs.setContext()
	defer hs.context.Close()

	results := make([]sms.Sms, 0)
	err := hs.collection.Find(query).Sort(sort).Limit(limit).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (hs smsStore) GetBySortLimitSkip(query bson.M, sort string, limit, skip int) ([]sms.Sms, error) {
	hs.setContext()
	defer hs.context.Close()

	results := make([]sms.Sms, 0)
	err := hs.collection.Find(query).Sort(sort).Skip(skip).Limit(limit).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (hs smsStore) UpdateSet(what bson.M, set bson.M) error {
	hs.setContext()
	defer hs.context.Close()

	err := hs.collection.Update(what, bson.M{"$set": set})
	return err
}

func (hs smsStore) Delete(what bson.M, all bool) error {
	hs.setContext()
	defer hs.context.Close()

	var err error
	if all {
		_, err = hs.collection.RemoveAll(what)
	} else {
		err = hs.collection.Remove(what)
	}
	if err != nil && err == mgo.ErrNotFound {
		return nil
	}
	return err
}

func (hs smsStore) CountBy(by bson.M) (int, error) {
	hs.setContext()
	defer hs.context.Close()

	return hs.collection.Find(by).Count()
}
