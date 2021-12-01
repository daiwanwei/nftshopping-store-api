package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nftshopping-store-api/pkg/databases"
	"nftshopping-store-api/pkg/utils"
)

type UserDao interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (user *User, err error)
	FindByAccount(ctx context.Context, account string) (user *User, err error)
	Save(ctx context.Context, user *User) (err error)
	Create(ctx context.Context, user *User) (err error)
	Delete(ctx context.Context, id primitive.ObjectID) (err error)
	ExistByID(ctx context.Context, id primitive.ObjectID) (isExisted bool, err error)
	ExistByAccount(ctx context.Context, account string) (isExisted bool, err error)
	FindAllByPage(ctx context.Context, pageable utils.Pageable) (users *utils.Page, err error)
}

type userDao struct {
	collection *mongo.Collection
}

func NewUserDao() (dao UserDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	col := db.Collection("user")
	opt := options.Index().SetUnique(true)
	col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"account": 1, // index in ascending order
			}, Options: opt,
		},
	})
	return &userDao{col}, nil
}

func (dao *userDao) FindByID(ctx context.Context, id primitive.ObjectID) (user *User, err error) {
	user = &User{}
	err = dao.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}
	return
}

func (dao *userDao) FindByAccount(ctx context.Context, account string) (user *User, err error) {
	user = &User{}
	err = dao.collection.FindOne(ctx, bson.D{{"account", account}}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}
	return
}

func (dao *userDao) Save(ctx context.Context, user *User) (err error) {
	oldUser := User{}
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"account", user.Account},
	}}}
	option := options.FindOneAndUpdate().SetUpsert(true)
	err = dao.collection.FindOneAndUpdate(ctx, filter, update, option).Decode(&oldUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return
	}
	return
}

func (dao *userDao) Delete(ctx context.Context, id primitive.ObjectID) (err error) {
	user := User{}
	err = dao.collection.FindOneAndDelete(ctx, bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return UserNotFound
		}
		return
	}
	return
}

func (dao *userDao) ExistByID(ctx context.Context, id primitive.ObjectID) (isExisted bool, err error) {
	count, err := dao.collection.CountDocuments(ctx, bson.D{{"_id", id}})
	if err != nil {
		return
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (dao *userDao) ExistByAccount(ctx context.Context, account string) (isExisted bool, err error) {
	count, err := dao.collection.CountDocuments(ctx, bson.D{{"account", account}})
	if err != nil {
		return
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (dao *userDao) FindAllByPage(ctx context.Context, pageable utils.Pageable) (users *utils.Page, err error) {
	filter := bson.D{}
	users, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *userDao) findPage(
	ctx context.Context, filter interface{}, pageable utils.Pageable,
) (page *utils.Page, err error) {
	total, err := dao.collection.CountDocuments(ctx, filter)
	page = &utils.Page{Size: pageable.Size, Page: pageable.Page, Total: total}
	option := options.Find()
	option.SetSkip(int64(pageable.Size * pageable.Page))
	option.SetLimit(int64(pageable.Size))
	cur, err := dao.collection.Find(ctx, filter, option)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	var users []User
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	page.Content = users
	page.TotalPage = utils.GetTotalPage(int64(page.Size), page.Total)
	return
}

func (dao *userDao) findList(ctx context.Context, filter interface{}) (users []User, err error) {
	cur, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *userDao) Create(ctx context.Context, user *User) (err error) {
	_, err = dao.collection.InsertOne(ctx, user)
	return
}

type User struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Account string             `bson:"account" json:"account"`
}

var UserNotFound = errors.New("user not found")
