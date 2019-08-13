package main

import (
	"github.com/globalsign/mgo/bson"
	"github.com/siskinc/mgorm"
)

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Username string
	Password string
}

func main() {
	// set the default mongodb infomation
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
	user := &User{}
	user.Username = "123"
	user.Password = "123"
	col := mgorm.Colletion("testdb", "user")
	user.ID = bson.NewObjectId()
	err = mgorm.Save(col, user, user.ID)
	if err != nil {
		panic(err)
	}

	user.Password = "1234"
	err = mgorm.Save(col, user, user.ID)
	if err != nil {
		panic(err)
	}

}
