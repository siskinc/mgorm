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

type collections map[string]map[string]*mgo.Collection
type dbs map[string]*mgo.Database

var (
	defaultHosts                   string
	defaultDatabase                string
	defaultUsername                string
	defaultPassword                string
	defaultCollectionTimeoutSecond int
	defaultMgoSession              *mgo.Session
	defaultMgoDatabase             *mgo.Database
	allColletions                  collections
	allDbs                         dbs
)

func init() {
	allColletions = make(map[string]map[string]*mgo.Collection, 0)
	allDbs = make(map[string]*mgo.Database, 0)
}

func Colletion(database, collection string) *mgo.Collection {
	databaseMap, ok := allColletions[database]
	if !ok {
		databaseMap = make(map[string]*mgo.Collection, 0)
		allColletions[database] = databaseMap
		db := allDbs[database]
		if db == nil {
			db = DefaultDatabase(database)
		}
		allColletions[database][collection] = db.C(collection)
	}
	return databaseMap[collection]
}

func DefaultDatabase(database string) *mgo.Database {
	db := allDbs[database]
	if db == nil {
		db = defaultMgoSession.DB(database)
	}
	return db
}

func DefaultMgoInfo(hosts, database, username, password string, collectionTimeoutSecond int) error {
	defaultHosts = hosts
	defaultDatabase = database
	defaultUsername = username
	defaultPassword = password
	defaultCollectionTimeoutSecond = collectionTimeoutSecond
	dialInfo := &mgo.DialInfo{
		Addrs:    strings.Split(defaultHosts, ","),
		Timeout:  time.Duration(defaultCollectionTimeoutSecond) * time.Second,
		Database: defaultDatabase,
		Username: defaultUsername,
		Password: defaultPassword,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}
	defaultMgoSession = session
	defaultMgoSession.SetMode(mgo.Monotonic, true)
	defaultMgoDatabase = session.DB(defaultDatabase)
	return nil
}

func DefaultMongoDBClient(collection string) *MongoDBClient {
	return &MongoDBClient{
		Hosts:                   defaultHosts,
		Database:                defaultDatabase,
		Collection:              collection,
		Username:                defaultUsername,
		Password:                defaultPassword,
		CollectionTimeoutSecond: defaultCollectionTimeoutSecond,
		Session: Session{
			S:  defaultMgoSession,
			Db: defaultMgoDatabase,
			C:  defaultMgoDatabase.C(collection),
		},
	}
}

// Connect return an error if connect mongodb is exception
// If the MongoDBClient Object Session's S is not nil point,
// this function will stop run and return an error of nil.
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

// DB will set the MongoDBClient Object Session's DB filed
func (m *MongoDBClient) DB() {
	if m.Session.Db != nil {
		return
	}
	db := m.Session.S.DB(m.Database)
	m.Session.Db = db
}

// C will set the MongoDBClient Object Session's C filed
func (m *MongoDBClient) C() {
	if m.Session.Db == nil {
		m.DB()
	}
	if m.Session.C != nil {
		return
	}
	c := m.Session.Db.C(m.Collection)
	m.Session.C = c
}

// Close
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
