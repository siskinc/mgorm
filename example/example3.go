package main

import (
	"fmt"
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
		"193.112.25.176:27017",
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
	err = mgorm.Save(col, user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	user.Password = "1234"
	err = mgorm.Save(col, user)
	if err != nil {
		panic(err)
	}

	err = mgorm.Delete(col, user)

}
