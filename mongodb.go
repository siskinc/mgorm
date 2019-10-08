package mgorm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Session struct
type Session struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
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

type collections map[string]map[string]*mongo.Collection

//type gridFSList map[string]map[string]*mongo.GridFS
type dbs map[string]*mongo.Database

var (
	defaultDatabase                string
	defaultCollectionTimeoutSecond int
	defaultMongoClient             *mongo.Client
	defaultMongoDatabase           *mongo.Database
	defaultContext                 = context.Background()
	allCollections                 collections
	//allGridFS                      gridFSList
	allDbs dbs
)

func init() {
	allCollections = make(collections, 0)
	//allGridFS = make(gridFSList, 0)
	allDbs = make(dbs, 0)
}

func Collection(database, collection string) *mongo.Collection {
	databaseMap, ok := allCollections[database]
	if !ok {
		databaseMap = make(map[string]*mongo.Collection, 0)
		allCollections[database] = databaseMap
		db, ok := allDbs[database]
		if !ok {
			db = DefaultDatabase(database)
			allDbs[database] = db
		}
		allCollections[database][collection] = db.Collection(collection)
	}
	c, ok := databaseMap[collection]
	if !ok {
		db := allDbs[database]
		c = db.Collection(collection)
		allCollections[database][collection] = c
	}
	return c
}

func DefaultDatabase(database string) *mongo.Database {
	db := allDbs[database]
	if db == nil {
		db = defaultMongoClient.Database(database)
	}
	return db
}

func DefaultMongoInfo(uri, database string, collectionTimeoutSecond int) (err error) {
	defaultDatabase = database
	defaultCollectionTimeoutSecond = collectionTimeoutSecond
	defaultContext, _ = context.WithTimeout(defaultContext, time.Duration(defaultCollectionTimeoutSecond)*time.Second)
	clientOptions := options.Client().ApplyURI(uri)
	defaultMongoClient, err = mongo.NewClient(clientOptions)
	if nil != err {
		return
	}
	err = defaultMongoClient.Connect(defaultContext)
	if nil != err {
		return
	}
	defaultMongoDatabase = defaultMongoClient.Database(defaultDatabase)
	return
}

func DefaultMongoDBClient(collection string) *MongoDBClient {
	return &MongoDBClient{
		Session: Session{
			Client:     defaultMongoClient,
			Database:   defaultMongoDatabase,
			Collection: defaultMongoDatabase.Collection(collection),
		},
	}
}

// DB will set the MongoDBClient Object Session's DB filed
func (m *MongoDBClient) DB() {
	if m.Session.Database != nil {
		return
	}
	db := m.Session.Client.Database(m.Database)
	m.Session.Database = db
}

// Collection will set the MongoDBClient Object Session's Collection filed
func (m *MongoDBClient) C() {
	if m.Session.Database == nil {
		m.DB()
	}
	if m.Session.Collection != nil {
		return
	}
	c := m.Session.Database.Collection(m.Collection)
	m.Session.Collection = c
}

// Close
func (m *MongoDBClient) Close() {
	if m.Client != nil {
		err := m.Client.Disconnect(defaultContext)
		if nil != err {
			panic(err)
		}
	}
}

// GetCollection func
func (m *MongoDBClient) GetCollection() *mongo.Collection {
	if m.Session.Collection == nil {
		m.C()
	}
	return m.Session.Collection
}

// Save func
func (m *MongoDBClient) Save(model interface{}) (result *mongo.InsertOneResult, err error) {
	collection := m.GetCollection()
	result, err = collection.InsertOne(defaultContext, model)
	return
}

// BatchSave func
func (m *MongoDBClient) BatchSave(modelList []interface{}) (result *mongo.InsertManyResult, err error) {
	collection := m.GetCollection()
	result, err = collection.InsertMany(defaultContext, modelList)
	return
}

// Find func
func (m *MongoDBClient) Find(filter interface{}, opts ...*options.FindOptions) (cursor *mongo.Cursor, err error) {
	collection := m.GetCollection()
	cursor, err = collection.Find(defaultContext, filter, opts...)
	return
}

// FindOne func
func (m *MongoDBClient) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	collection := m.GetCollection()
	return collection.FindOne(defaultContext, filter, opts...)
}

// FindAll func
func (m *MongoDBClient) FindAll(filter interface{}, opts ...*options.FindOptions) (cursor *mongo.Cursor, err error) {
	return m.Find(filter, opts...)
}

func (m *MongoDBClient) FindPage(filter interface{}, iPageSize, iPageIndex int64, SortedStrs ...string) (cursor *mongo.Cursor, err error) {
	opt := &options.FindOptions{}
	skip := iPageSize * (iPageIndex - 1)
	opt = opt.SetLimit(iPageSize).SetSkip(skip)
	for _, sortStr := range SortedStrs {
		opt = opt.SetSort(sortStr)
	}
	return m.Find(filter, opt)
}

// Count func
func (m *MongoDBClient) Count(filter interface{}) (count int64, err error) {
	collection := m.GetCollection()
	count, err = collection.CountDocuments(defaultContext, filter)
	return
}

// Exist func
func (m *MongoDBClient) Exist(query interface{}) (exist bool, err error) {
	exist = false
	var count int64 = -1
	count, err = m.Count(query)
	if err != nil {
		return
	}
	if count > 0 {
		exist = true
	}
	return
}

func (m *MongoDBClient) UpdateMany(filter, update interface{}) (*mongo.UpdateResult, error) {
	collection := m.GetCollection()
	return collection.UpdateMany(defaultContext, filter, update)
}

func (m *MongoDBClient) UpdateOne(filter, update interface{}) (*mongo.UpdateResult, error) {
	collection := m.GetCollection()
	return collection.UpdateOne(defaultContext, filter, update)
}

// DeleteModel func
func (m *MongoDBClient) DeleteModel(query interface{}, isOne bool) (*mongo.DeleteResult, error) {
	collection := m.GetCollection()

	if isOne {
		return collection.DeleteOne(defaultContext, query)
	}
	return collection.DeleteMany(defaultContext, query)
}

func (m *MongoDBClient) CreateIndex(index mongo.IndexModel) (string, error) {
	collection := m.GetCollection()
	return collection.Indexes().CreateOne(defaultContext, index)
}

func (m *MongoDBClient) DropIndex(name string) (bson.Raw, error) {
	collection := m.GetCollection()
	return collection.Indexes().DropOne(defaultContext, name)
}
