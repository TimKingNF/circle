package sqldriver

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataM map[string]interface{}

type SQLCollection interface {
	Insert(docs ...interface{}) error

	//	select one
	SelectOne(selector, result interface{}) error

	//	select muti data
	SelectAll(selector, doc interface{}) []DataM

	//	range all data
	Iter(doc interface{}) []DataM

	//	update one
	UpdateOne(selector, newData interface{}) error

	//	update muti data
	UpdateAll(selector, newData interface{}) error

	//	delete one
	RemoveOne(selector interface{}) error

	//	delete muti data
	RemoveAll(selector interface{}) error
}

type mySQLCollection struct {
	collection *mgo.Collection
}

func (collection *mySQLCollection) Insert(docs ...interface{}) error {
	var lastErr error
	for k, _ := range docs {
		err := collection.collection.Insert(docs[k])
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (collection *mySQLCollection) SelectOne(selector, result interface{}) error {
	return collection.collection.Find(selector).One(result)
}

func (collection *mySQLCollection) SelectAll(selector, doc interface{}) []DataM {
	var result []DataM
	iter := collection.collection.Find(selector).Iter()
	for iter.Next(&doc) {
		result = append(result, DataM(doc.(bson.M)))
	}
	return result
}

func (collection *mySQLCollection) Iter(doc interface{}) []DataM {
	var result []DataM
	iter := collection.collection.Find(nil).Iter()
	for iter.Next(&doc) {
		result = append(result, DataM(doc.(bson.M)))
	}
	return result
}

func (collection *mySQLCollection) UpdateOne(selector, newData interface{}) error {
	return collection.collection.Update(selector, newData)
}

func (collection *mySQLCollection) UpdateAll(selector, newData interface{}) error {
	_, err := collection.collection.UpdateAll(selector, newData)
	return err
}

func (collection *mySQLCollection) RemoveOne(selector interface{}) error {
	return collection.collection.Remove(selector)
}

func (collection *mySQLCollection) RemoveAll(selector interface{}) error {
	_, err := collection.collection.RemoveAll(selector)
	return err
}
