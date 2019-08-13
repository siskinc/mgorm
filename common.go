package mgorm

import (
	"github.com/globalsign/mgo/bson"
	"reflect"
)

// NewUpdater func
func NewUpdater(updater interface{}) interface{} {
	return map[string]interface{}{
		"$set": updater,
	}
}

func setSomeField(model interface{}, fieldName string, fieldValue reflect.Value) {
	v := reflect.ValueOf(model)
	v.Elem().FieldByName(fieldName).Set(fieldValue)
}

func getSomeField(model interface{}, fieldName string) reflect.Value {
	v := reflect.ValueOf(model)
	return v.Elem().FieldByName(fieldName)
}

func setObjectID(model interface{}, id bson.ObjectId) {
	vid := reflect.ValueOf(id)
	setSomeField(model, "ID", vid)
}

func getObjectID(model interface{}) bson.ObjectId {
	ret := getSomeField(model, "ID")
	return ret.Interface().(bson.ObjectId)
}
