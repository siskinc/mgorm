package mgorm

import "github.com/globalsign/mgo/bson"

type CustomCollectionName interface {
	CollectionName() string
}

type CustomID interface {
	GetID() bson.ObjectId
	SetID(bson.ObjectId)
}
