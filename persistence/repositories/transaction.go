package repositories

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nftshopping-store-api/pkg/databases"
	"nftshopping-store-api/pkg/utils"
	"time"
)

type TransactionDao interface {
	Find(ctx context.Context, tnxId primitive.ObjectID) (transaction *Transaction, err error)
	Create(ctx context.Context, transaction *Transaction) (err error)
	Delete(ctx context.Context, creationId primitive.ObjectID) (err error)
	FindAll(ctx context.Context) (transactions []Transaction, err error)
	FindAllByPage(ctx context.Context, pageable utils.Pageable) (transactions *utils.Page, err error)
	FindAllByFilter(ctx context.Context, filter TransactionFilter) (transactions []Transaction, err error)
	FindAllByFilterAndPage(ctx context.Context, filter TransactionFilter, pageable utils.Pageable) (transactions *utils.Page, err error)
}

type transactionDao struct {
	collection *mongo.Collection
}

func NewTransactionDao() (dao TransactionDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &transactionDao{db.Collection("creation_transaction")}, nil
}

func (dao *transactionDao) Create(ctx context.Context, transaction *Transaction) (err error) {
	_, err = dao.collection.InsertOne(ctx, transaction)
	return
}

func (dao *transactionDao) Find(
	ctx context.Context, creationId primitive.ObjectID,
) (transaction *Transaction, err error) {
	transaction = &Transaction{}
	err = dao.collection.FindOne(ctx, bson.D{{"_id", creationId}}).Decode(transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}
	return
}

func (dao *transactionDao) Delete(ctx context.Context, creationId primitive.ObjectID) (err error) {
	transaction := Transaction{}
	err = dao.collection.FindOneAndDelete(ctx, bson.D{{"_id", creationId}}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return TransactionNotFound
		}
		return
	}
	return
}

func (dao *transactionDao) FindAll(ctx context.Context) (transactions []Transaction, err error) {
	cur, err := dao.collection.Find(ctx, bson.D{{}})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var transaction Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *transactionDao) FindAllByPage(
	ctx context.Context, pageable utils.Pageable,
) (transactions *utils.Page, err error) {
	filter := bson.D{{}}
	transactions, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *transactionDao) FindAllByFilter(
	ctx context.Context, filter TransactionFilter,
) (transactions []Transaction, err error) {
	transactions, err = dao.findList(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (dao *transactionDao) FindAllByFilterAndPage(
	ctx context.Context, filter TransactionFilter, pageable utils.Pageable,
) (transactions *utils.Page, err error) {
	transactions, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *transactionDao) findList(
	ctx context.Context, filter interface{},
) (transactions []Transaction, err error) {
	cur, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var transaction Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *transactionDao) findPage(
	ctx context.Context, filter interface{}, pageable utils.Pageable,
) (transactions *utils.Page, err error) {
	total, err := dao.collection.CountDocuments(ctx, filter)
	transactions = &utils.Page{Size: pageable.Size, Page: pageable.Page, Total: total}
	option := options.Find()
	option.SetSkip(int64(pageable.Size * pageable.Page))
	option.SetLimit(int64(pageable.Size))
	if len(pageable.Sort) > 0 {
		sort := bson.D{}
		for key, value := range pageable.Sort {
			sort = append(sort, bson.E{Key: key, Value: value})
		}
		option.SetSort(sort)
	}
	cur, err := dao.collection.Find(ctx, filter, option)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	var content []Transaction
	for cur.Next(ctx) {
		var transaction Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		content = append(content, transaction)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	transactions.Content = content
	return
}

type Transaction struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	CreationID string             `bson:"creation_id" json:"creationId"`
	BrandID    string             `bson:"brand_id" json:"brandId"`
	Buyer      string             `bson:"buyer" json:"buyer"`
	Seller     string             `bson:"seller" json:"seller"`
	Amount     int                `bson:"amount" json:"amount"`
	Price      int                `bson:"price" json:"price"`
	TradeAt    time.Time          `bson:"trade_at" json:"tradeAt"`
}

type TransactionFilter bson.D

func SelectorOfTransaction(selector TransactionSelector) (filter TransactionFilter) {
	filter = TransactionFilter{}

	if selector.Buyer != nil {
		filter = append(filter, bson.E{
			Key: "buyer", Value: selector.Buyer,
		})
	}

	if selector.Seller != nil {
		filter = append(filter, bson.E{
			Key: "seller", Value: selector.Seller,
		})
	}

	if selector.CreationID != nil {
		filter = append(filter, bson.E{
			Key: "creation_id", Value: selector.CreationID,
		})
	}

	if selector.BrandID != nil {
		filter = append(filter, bson.E{
			Key: "brand_id", Value: selector.BrandID,
		})
	}

	if selector.TradedBefore != nil || selector.TradedAfter != nil {
		tradeFilter := bson.D{}
		if selector.TradedAfter != nil {
			tradeFilter = append(tradeFilter, bson.E{Key: "$gte", Value: selector.TradedAfter})
		}
		if selector.TradedBefore != nil {
			tradeFilter = append(tradeFilter, bson.E{Key: "$lte", Value: selector.TradedBefore})
		}
		filter = append(filter, bson.E{Key: "trade_at", Value: tradeFilter})
	}

	if selector.MaxPrice != nil || selector.MinPrice != nil {
		priceFilter := bson.D{}
		if selector.MinPrice != nil {
			priceFilter = append(priceFilter, bson.E{Key: "$gte", Value: selector.MinPrice})
		}
		if selector.MaxPrice != nil {
			priceFilter = append(priceFilter, bson.E{Key: "$lte", Value: selector.MaxPrice})
		}
		filter = append(filter, bson.E{Key: "price", Value: priceFilter})
	}
	return
}

type TransactionSelector struct {
	CreationID   *string    `json:"creationID"`
	BrandID      *string    `json:"brandId"`
	Buyer        *string    `json:"buyer"`
	Seller       *string    `json:"seller"`
	TradedBefore *time.Time `json:"tradedBefore"`
	TradedAfter  *time.Time `json:"tradedAfter"`
	MaxPrice     *int       `json:"maxPrice"`
	MinPrice     *int       `json:"minPrice"`
}

var TransactionNotFound = errors.New("transaction not found")
