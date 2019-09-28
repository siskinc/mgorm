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
	if customIDModel, ok := model.(CustomID); ok {
		customIDModel.SetID(id)
		return
	}
	vid := reflect.ValueOf(id)
	setSomeField(model, "ID", vid)
}

func getObjectID(model interface{}) bson.ObjectId {
	if customIDModel, ok := model.(CustomID); ok {
		return customIDModel.GetID()
	}
	ret := getSomeField(model, "ID")
	return ret.Interface().(bson.ObjectId)
}

func MergeBsonM(obj1, obj2 bson.M) bson.M {
	result := bson.M{}
	for k, v := range obj1 {
		result[k] = v
	}

	for k, v := range obj2 {
		result[k] = v
	}

	return result
}
