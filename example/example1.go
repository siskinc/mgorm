package main

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/siskinc/mgorm"
)

func main() {
	var err error
	// set the default mongodb infomation
	err = mgorm.DefaultMgoInfo(
		"127.0.0.1:27017",
		"testdb",
		"",
		"",
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
	err = client.Save(bson.M{"_id": bson.NewObjectId(), "name": "daryl"})
	if err != nil {
		log.Fatalln("Save is err", err)
	}
	err = client.Save(bson.M{"_id": bson.NewObjectId(), "name": "siskinc"})
	if err != nil {
		log.Fatalln("Save is err", err)
	}

	// Get one document
	// the first way
	result := make(map[string]interface{}, 0)
	query := bson.M{"name": "daryl"}
	err = client.Find(result, query, true)
	if err != nil {
		log.Fatalln("Get one document is err", err)
	}
	fmt.Println("Get one document, the first way, result", result)
	// the second way
	result = make(map[string]interface{}, 0)
	err = client.FindOne(result, query)
	if err != nil {
		log.Fatalln("Get one document is err", err)
	}
	fmt.Println("Get one document, the second way, result", result)

	// Get all document
	// the first way
	allResult := make([]map[string]interface{}, 0)
	allQuery := bson.M{}
	err = client.Find(&allResult, allQuery, false)
	if err != nil {
		log.Fatalln("Get all document is err", err)
	}
	fmt.Println("Get all document, the first way, result", allResult)

	// the second way
	allResult = make([]map[string]interface{}, 0)
	err = client.FindAll(&allResult, allQuery)
	if err != nil {
		log.Fatalln("Get all document is err", err)
	}
	fmt.Println("Get all document, the second way, result", allResult)

	// the third way
	allResult = make([]map[string]interface{}, 0)
	iter, err := client.FindAll4Iter(allQuery)
	if err != nil {
		log.Fatalln("Get all document is err", err)
	}
	fmt.Print("Get all document, the third way, result ")
	for iter.Next(result) {
		fmt.Print(result, " ")
	}
}
