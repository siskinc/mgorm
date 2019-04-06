package main

import (
	"github.com/globalsign/mgo/bson"
	"github.com/siskinc/mgorm"
)

func main() {
	// set the default mongodb infomation
	mgorm.DefaultMgoInfo(
		"127.0.0.1:27017",
		"testdb",
		"",
		"",
		30,
	)
	client := mgorm.DefaultMongoDBClient("name")
	client.Save(bson.M{"_id": bson.NewObjectId(), "name": "daryl"})
}
