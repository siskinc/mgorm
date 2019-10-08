package main

import (
	"github.com/globalsign/mgo"
	"github.com/siskinc/mgorm"
)

func main() {
	err := mgorm.DefaultMgoInfo(
		"127.0.0.1:27017",
		"testdb",
		"",
		"",
		2,
	)
	if err != nil {
		panic(err)
	}
	client := mgorm.DefaultMongoDBClient("people")
	DropIndex(client)
	//CreateIndex(client)

}

func CreateIndex(client *mgorm.MongoDBClient) {
	index := mgo.Index{
		Key:    []string{"-1:name"},
		Unique: true,
	}
	err := client.EnsureIndex(index)
	if nil != err {
		panic(err)
	}
	index2 := mgo.Index{
		Key:              []string{"username", "name"},
		Unique:           false,
		DropDups:         false,
		Background:       false,
		Sparse:           false,
		PartialFilter:    nil,
		ExpireAfter:      0,
		Name:             "",
		Min:              0,
		Max:              0,
		Minf:             0,
		Maxf:             0,
		BucketSize:       0,
		Bits:             0,
		DefaultLanguage:  "",
		LanguageOverride: "",
		Weights:          nil,
		Collation:        nil,
	}
	err = client.EnsureIndex(index2)
	if nil != err {
		panic(err)
	}
}

func DropIndex(client *mgorm.MongoDBClient) {
	err := client.DropIndex("-1:name")
	if nil != err {
		panic(err)
	}
	err = client.DropIndex("username", "name")
	if nil != err {
		panic(err)
	}
}
