package databases

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var dbInstance *db

func GetDB() (instance *db, err error) {
	if dbInstance == nil {
		mongo, err := GetMongoDB()
		if err != nil {
			return nil, err
		}
		dbInstance = &db{
			Mongo: mongo,
		}
	}
	return dbInstance, nil
}

type db struct {
	Mongo *mongo.Database
}
