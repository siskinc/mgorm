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
	db := mgorm.DefaultDatabase("testdb")
	col := db.C("name")
	col.Insert(bson.M{"_id": bson.NewObjectId(), "name": "daryl_test_db"})
	col2 := mgorm.Colletion("testdb", "name")
	col2.Insert(bson.M{"_id": bson.NewObjectId(), "name": "daryl_test_collection"})
}
