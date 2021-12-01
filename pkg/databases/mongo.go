package databases

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"nftshopping-store-api/pkg/config"
)

var (
	mongoInstance   *mongo.Client
	mongoDBInstance *mongo.Database
)

func GetMongo() (instance *mongo.Client, err error) {
	if mongoInstance == nil {
		instance, err = newMongo()
		if err != nil {
			panic(err)
		}
		mongoInstance = instance
	}
	return mongoInstance, nil
}

func newMongo() (*mongo.Client, error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	mongoConfig := c.Database.Mongo
	clientOpts := options.Client().ApplyURI(mongoConfig.Uri)
	instance, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return nil, err
	}
	return instance, err
}

func GetMongoDB() (instance *mongo.Database, err error) {
	if mongoDBInstance == nil {
		instance, err = newMongoDB()
		if err != nil {
			panic(err)
		}
		mongoDBInstance = instance
	}
	return mongoDBInstance, nil
}

func newMongoDB() (instance *mongo.Database, err error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	mongoConfig := c.Database.Mongo
	m, err := GetMongo()
	if err != nil {
		return nil, err
	}
	instance = m.Database(mongoConfig.Source)
	err = updateCollection(instance)
	if err != nil {
		return
	}
	return
}

var collections = []string{"auth", "user", "validation"}

func updateCollection(mongodb *mongo.Database) (err error) {
	for _, collection := range collections {
		err := mongodb.CreateCollection(context.Background(), collection)
		if err != nil {
			if commonErr, ok := err.(mongo.CommandError); !ok {
				return err
			} else {
				if commonErr.Code != 48 {
					fmt.Println(commonErr.Code)
					return err
				}
			}
		}
	}
	return
}

type MongoTransaction struct {
	Session mongo.Session
	Options *options.TransactionOptions
}

func NewMongoTransaction() (mongoTransaction *MongoTransaction, err error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	m, err := GetMongo()
	if err != nil {
		return nil, err
	}
	session, err := m.StartSession()
	if err != nil {
		return
	}
	return &MongoTransaction{
		Session: session,
		Options: txnOpts,
	}, nil
}

func (txn *MongoTransaction) Start() error {
	return txn.Session.StartTransaction(txn.Options)
}

func (txn *MongoTransaction) Commit(ctx context.Context) (err error) {
	err = txn.Session.CommitTransaction(ctx)
	return
}

func (txn *MongoTransaction) Abort(ctx context.Context) error {
	return txn.Session.AbortTransaction(ctx)
}

func (txn *MongoTransaction) End(ctx context.Context) {
	txn.Session.EndSession(ctx)
}
