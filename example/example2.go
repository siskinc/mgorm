package main

import (
	"context"
	"github.com/siskinc/mgorm"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// set the default mongodb infomation
	mgorm.DefaultMongoInfo(
		"mongodb://127.0.0.1:27017/",
		"testdb",
		30,
	)
	db := mgorm.DefaultDatabase("testdb")
	col := db.Collection("name")
	col.InsertOne(context.Background(), bson.M{"name": "daryl_test_db"})
	col2 := mgorm.Collection("testdb", "name")
	col2.InsertOne(context.Background(), bson.M{"name": "daryl_test_collection"})
}
