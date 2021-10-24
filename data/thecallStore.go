package data

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"sms-sorter/model/thecall"
	"time"
)

type thecallStore struct {
	collection *mgo.Collection
	context    *Context
}

func (ts *thecallStore) setContext() {
	ts.context = NewTokyoContext()
	ts.collection = ts.context.SpamDBCollection(CTheCall)
}

func NewTheCallStore() *thecallStore {
	return &thecallStore{}
}

func (ts thecallStore) Create(t *thecall.TheCall) error {
	ts.setContext()
	defer ts.context.Close()

	t.ID = bson.NewObjectId()
	t.CreatedAt = time.Now().Unix()

	return ts.collection.Insert(t)
}

func (ts thecallStore) GetBy(query bson.M) ([]thecall.TheCall, error) {
	ts.setContext()
	defer ts.context.Close()

	results := make([]thecall.TheCall, 0)
	err := ts.collection.Find(query).All(&results)

	return results, err
}

func (ts thecallStore) GetByDesc(query bson.M) ([]thecall.TheCall, error) {
	ts.setContext()
	defer ts.context.Close()

	results := make([]thecall.TheCall, 0)
	err := ts.collection.Find(query).Sort("-created_at").All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ts thecallStore) GetBySortLimit(query bson.M, sort string, limit int) ([]thecall.TheCall, error) {
	ts.setContext()
	defer ts.context.Close()

	results := make([]thecall.TheCall, 0)
	err := ts.collection.Find(query).Sort(sort).Limit(limit).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ts thecallStore) GetBySortLimitSkip(query bson.M, sort string, limit, skip int) ([]thecall.TheCall, error) {
	ts.setContext()
	defer ts.context.Close()

	results := make([]thecall.TheCall, 0)
	err := ts.collection.Find(query).Sort(sort).Skip(skip).Limit(limit).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (ts thecallStore) UpdateSet(what bson.M, set bson.M) error {
	ts.setContext()
	defer ts.context.Close()

	err := ts.collection.Update(what, bson.M{"$set": set})
	return err
}

func (ts thecallStore) Delete(what bson.M, all bool) error {
	ts.setContext()
	defer ts.context.Close()

	var err error
	if all {
		_, err = ts.collection.RemoveAll(what)
	} else {
		err = ts.collection.Remove(what)
	}
	if err != nil && err == mgo.ErrNotFound {
		return nil
	}
	return err
}

func (ts thecallStore) CountBy(by bson.M) (int, error) {
	ts.setContext()
	defer ts.context.Close()

	return ts.collection.Find(by).Count()
}
