package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"

	"github.com/siskinc/mgorm"
)

func main() {
	log.SetFlags(1)
	var err error
	// set the default mongodb infomation
	err = mgorm.DefaultMongoInfo(
		"mongodb://127.0.0.1:27017/",
		"testdb",
		30,
	)
	if err != nil {
		log.Fatalln("Connect mongodb is err:", err)
	}

	// Get MongoDBClient Object
	client := mgorm.DefaultMongoDBClient("name")
	if client == nil {
		log.Fatalln("Get MongoDBClient Object is nil")
	}
	_, err = client.Save(bson.M{"name": "daryl"})
	if err != nil {
		log.Fatalln("Save is err", err)
	}
	_, err = client.Save(bson.M{"name": "siskinc"})
	if err != nil {
		log.Fatalln("Save is err", err)
	}

	// Get one document
	// the first way
	query := bson.M{"name": "daryl"}
	cursor, err := client.Find(query)
	if err != nil {
		log.Fatalln("Get one document is err", err)
	}
	fmt.Println("Get one document, the first way, cursor", cursor)
	// the second way
	result := client.FindOne(query)
	fmt.Println("Get one document, the second way, result", result)
}
