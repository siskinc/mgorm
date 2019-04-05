package mgorm

import (
	"strings"
	"time"

	"github.com/globalsign/mgo"
)

// Session struct
type Session struct {
	S  *mgo.Session
	Db *mgo.Database
	C  *mgo.Collection
}

// MongoDBClient struct
type MongoDBClient struct {
	Hosts                   string // example:127.0.0.1, 11.100.1.1
	Database                string
	Collection              string
	Username                string
	Password                string
	CollectionTimeoutSecond int
	Session
}

// Connect func 连接MongoDB
func (m *MongoDBClient) Connect() error {
	if m.S != nil {
		return nil
	}
	dialInfo := &mgo.DialInfo{
		Addrs:    strings.Split(m.Hosts, ","),
		Timeout:  time.Duration(m.CollectionTimeoutSecond) * time.Second,
		Database: m.Database,
		Username: m.Username,
		Password: m.Password,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}
	m.S = session
	// Optional. Switch the session to a monotonic behavior.
	m.S.SetMode(mgo.Monotonic, true)
	return nil
}

// DB func 设置MongoDBClient的DB
func (m *MongoDBClient) DB() {
	db := m.Session.S.DB(m.Database)
	m.Session.Db = db
}

// C func 设置MongoDBClient的Collection
func (m *MongoDBClient) C() {
	if m.Session.Db == nil {
		m.DB()
	}
	c := m.Session.Db.C(m.Collection)
	m.Session.C = c
}

// Close func
func (m *MongoDBClient) Close() {
	if m.S != nil {
		m.S.LogoutAll()
		m.S.Close()
	}
}

// GetColletion func
func (m *MongoDBClient) GetColletion() *mgo.Collection {
	if m.Session.C == nil {
		m.C()
	}
	return m.Session.C
}

// Save func
func (m *MongoDBClient) Save(model interface{}) (err error) {
	err = m.Connect()
	if err != nil {
		return
	}
	defer m.Close()
	collection := m.GetColletion()
	err = collection.Insert(model)
	return
}

// Find func
func (m *MongoDBClient) Find(result interface{}, query interface{}, isOne bool) (err error) {
	err = m.Connect()
	if err != nil {
		return
	}
	defer m.Close()
	collection := m.GetColletion()
	if isOne {
		err = collection.Find(query).One(result)
	} else {
		err = collection.Find(query).All(result)
	}
	return
}

// FindOne func
func (m *MongoDBClient) FindOne(result interface{}, query interface{}) error {
	return m.Find(result, query, true)
}

// FindAll func
func (m *MongoDBClient) FindAll(result interface{}, query interface{}) error {
	return m.Find(result, query, false)
}

// FindAll4Iter func
func (m *MongoDBClient) FindAll4Iter(query interface{}) (iter *mgo.Iter, err error) {
	err = m.Connect()
	if err != nil {
		return
	}
	defer m.Close()
	collection := m.GetColletion()
	iter = collection.Find(query).Iter()
	err = iter.Err()
	if err != nil {
		return
	}
	return
}

func (m *MongoDBClient) FindPage(result interface{}, query interface{}, iPageSize, iPageIndex int, SortedStrs ...string) error {
	err := m.Connect()
	if err != nil {
		return err
	}
	defer m.Close()
	collection := m.GetColletion()
	skip := iPageSize * (iPageIndex - 1)
	if len(SortedStrs) == 0 {
		return collection.Find(query).Skip(skip).Limit(iPageSize).All(result)
	} else {
		return collection.Find(query).Sort(SortedStrs...).Skip(skip).Limit(iPageSize).All(result)
	}
}

// Count func
func (m *MongoDBClient) Count(query interface{}) (count int, err error) {
	count = -1
	err = m.Connect()
	if err != nil {
		return
	}
	defer m.Close()
	collection := m.GetColletion()
	count, err = collection.Find(query).Count()
	return
}

// Exist func
func (m *MongoDBClient) Exist(query interface{}) (exist bool, err error) {
	exist = false
	count := -1
	count, err = m.Count(query)
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

// Update func
func (m *MongoDBClient) Update(query, updater interface{}, isOne bool) error {
	err := m.Connect()
	if err != nil {
		return err
	}
	defer m.Close()
	collection := m.GetColletion()
	if isOne {
		err = collection.Update(query, updater)
	} else {
		_, err = collection.UpdateAll(query, updater)
	}
	return err
}

// Delete func
func (m *MongoDBClient) Delete(query interface{}, isOne bool) error {
	err := m.Connect()
	if err != nil {
		return err
	}
	defer m.Close()
	collection := m.GetColletion()
	if isOne {
		err = collection.Remove(query)
	} else {
		_, err = collection.RemoveAll(query)
	}
	return err
}

// DeleteById func
func (m *MongoDBClient) DeleteById(query interface{}, id string) error {
	err := m.Connect()
	if err != nil {
		return err
	}
	defer m.Close()
	collection := m.GetColletion()
	err = collection.RemoveId(id)
	return err
}

//DeleteOne func
func (m *MongoDBClient) DeleteOne(query interface{}) error {
	return m.Delete(query, false)
}

//DeleteAll func
func (m *MongoDBClient) DeleteAll(query interface{}) error {
	return m.Delete(query, true)
}

// UpdateOne func
func (m *MongoDBClient) UpdateOne(query, updater interface{}) error {
	return m.Update(query, updater, true)
}

// UpdateAll func
func (m *MongoDBClient) UpdateAll(query, updater interface{}) error {
	return m.Update(query, updater, false)
}
