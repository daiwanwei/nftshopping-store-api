package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nftshopping-store-api/pkg/databases"
)

type Auth struct {
	Name        string   `bson:"_id" json:"name"`
	Password    string   `bson:"password" json:"password"`
	Authorities []string `bson:"authorities" json:"authorities"`
}

func (a *Auth) GetName() string {
	return a.Name
}

func (a *Auth) GetAuthorities() []string {
	return a.Authorities
}

type AuthDao interface {
	FindByName(ctx context.Context, name string) (auth *Auth, err error)
	Create(ctx context.Context, auth *Auth) (err error)
	DeleteByName(ctx context.Context, name string) (err error)
	Save(ctx context.Context, auth *Auth) (err error)
}

type authDao struct {
	collection *mongo.Collection
}

func NewAuthDao() (dao AuthDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &authDao{db.Collection("auth")}, nil
}

func (dao *authDao) Create(ctx context.Context, auth *Auth) (err error) {
	_, err = dao.collection.InsertOne(ctx, auth)
	return
}

func (dao *authDao) DeleteByName(ctx context.Context, name string) (err error) {
	auth := &Auth{}
	err = dao.collection.FindOneAndDelete(ctx, bson.D{{"_id", name}}).Decode(auth)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return AuthNotFound
		}
		return
	}
	return
}

func (dao *authDao) FindByName(ctx context.Context, name string) (auth *Auth, err error) {
	auth = &Auth{}
	err = dao.collection.FindOne(ctx, bson.D{{"_id", name}}).Decode(auth)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}
	return
}

func (dao *authDao) Save(ctx context.Context, auth *Auth) (err error) {
	oldAuth := Auth{}
	filter := bson.D{{"_id", auth.Name}}
	update := bson.D{{"$set", bson.D{
		{"password", auth.Password},
		{"authorities", auth.Authorities},
	}}}
	option := options.FindOneAndUpdate().SetUpsert(true)
	err = dao.collection.FindOneAndUpdate(ctx, filter, update, option).Decode(&oldAuth)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return
	}
	return
}

var (
	AuthNotFound = errors.New("auth not found")
)
