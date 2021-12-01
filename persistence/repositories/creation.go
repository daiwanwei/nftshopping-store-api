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
	"time"
)

type CreationDao interface {
	Exist(ctx context.Context, id primitive.ObjectID) (isExisted bool, err error)
	Find(ctx context.Context, creationId primitive.ObjectID) (creation *Creation, err error)
	Create(ctx context.Context, creation *Creation) (err error)
	Save(ctx context.Context, creation *Creation) (err error)
	Delete(ctx context.Context, creationId primitive.ObjectID) (err error)
	FindAll(ctx context.Context) (creations []Creation, err error)
	FindAllByPage(ctx context.Context, pageable utils.Pageable) (creations *utils.Page, err error)
	FindAllByCreationName(ctx context.Context, creationName string) (creations []Creation, err error)
	FindAllByFilter(ctx context.Context, filter CreationFilter) (creations []Creation, err error)
	FindAllByFilterAndPage(ctx context.Context, filter CreationFilter, pageable utils.Pageable) (creations *utils.Page, err error)
}

type creationDao struct {
	collection *mongo.Collection
}

func NewCreationDao() (dao CreationDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &creationDao{db.Collection("creation")}, nil
}

func (dao *creationDao) Exist(ctx context.Context, id primitive.ObjectID) (isExisted bool, err error) {
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

func (dao *creationDao) Create(ctx context.Context, creation *Creation) (err error) {
	_, err = dao.collection.InsertOne(ctx, creation)
	return
}

func (dao *creationDao) Find(ctx context.Context, creationId primitive.ObjectID) (creation *Creation, err error) {
	creation = &Creation{}
	err = dao.collection.FindOne(ctx, bson.D{{"_id", creationId}}).Decode(creation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}
	return
}

func (dao *creationDao) Save(ctx context.Context, creation *Creation) (err error) {
	oldCreation := Creation{}
	filter := bson.D{{"_id", creation.ID}}
	update := bson.D{{"$set", bson.D{
		{"creation_name", creation.CreationName},
		{"creator", creation.Creator},
		{"small_image_url", creation.SmallImageURL},
		{"properties", creation.Properties},
		{"price", creation.Price},
		{"brand_id", creation.BrandID},
		{"sale_way", creation.SaleWay},
		{"sale_status", creation.SaleStatus},
		{"sale_start_at", creation.SaleStartAt},
		{"sale_end_at", creation.SaleEndAt},
		{"description", creation.Description},
	}}}
	option := options.FindOneAndUpdate().SetUpsert(true)
	err = dao.collection.FindOneAndUpdate(ctx, filter, update, option).Decode(&oldCreation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return
	}
	return
}

func (dao *creationDao) Delete(ctx context.Context, creationId primitive.ObjectID) (err error) {
	creation := Creation{}
	err = dao.collection.FindOneAndDelete(ctx, bson.D{{"_id", creationId}}).Decode(&creation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return CreationNotFound
		}
		return
	}
	return
}

func (dao *creationDao) FindAll(ctx context.Context) (creations []Creation, err error) {
	creations, err = dao.findList(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	return
}

func (dao *creationDao) FindAllByPage(ctx context.Context, pageable utils.Pageable) (creations *utils.Page, err error) {
	filter := bson.D{}
	creations, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *creationDao) FindAllByCreationName(ctx context.Context, creationName string) (creations []Creation, err error) {
	creations, err = dao.findList(ctx, bson.D{{"creation_name", creationName}})
	if err != nil {
		return
	}
	return
}

func (dao *creationDao) FindAllByFilter(ctx context.Context, filter CreationFilter) (creations []Creation, err error) {
	creations, err = dao.findList(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (dao *creationDao) FindAllByFilterAndPage(
	ctx context.Context, filter CreationFilter, pageable utils.Pageable,
) (creations *utils.Page, err error) {
	creations, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *creationDao) findList(ctx context.Context, filter interface{}) (creations []Creation, err error) {
	cur, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var creation Creation
		err := cur.Decode(&creation)
		if err != nil {
			return nil, err
		}
		creations = append(creations, creation)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *creationDao) findPage(
	ctx context.Context, filter interface{}, pageable utils.Pageable,
) (creations *utils.Page, err error) {
	total, err := dao.collection.CountDocuments(ctx, filter)
	creations = &utils.Page{Size: pageable.Size, Page: pageable.Page, Total: total}
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
	var content []Creation
	for cur.Next(ctx) {
		var creation Creation
		err := cur.Decode(&creation)
		if err != nil {
			return nil, err
		}
		content = append(content, creation)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	creations.Content = content

	return
}

type Creation struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	CreationName    string             `bson:"creation_name" json:"creationName"`
	Creator         string             `bson:"creator" json:"creator"`
	SmallImageURL   string             `json:"smallImageUrl"`
	Properties      []string           `bson:"properties" json:"properties"`
	Amount          int                `bson:"amount" json:"amount"`
	Price           int                `bson:"price" json:"price"`
	CreateAt        time.Time          `bson:"create_at" json:"createAt"`
	BrandID         string             `bson:"brand_id" json:"brandId"`
	SaleWay         string             `bson:"sale_way" json:"saleWay"`
	SaleStatus      string             `bson:"sale_status" json:"saleStatus"`
	SaleStartAt     time.Time          `bson:"sale_start_at" json:"saleStartAt"`
	SaleEndAt       time.Time          `bson:"sale_end_at" json:"saleEndAt"`
	Description     string             `bson:"description" json:"description"`
	ContractAddress string             `bson:"contract_address" json:"contractAddress"`
}

type CreationFilter bson.D

func SelectorOfCreation(selector CreationSelector) (filter CreationFilter) {
	filter = CreationFilter{}
	if selector.CreationIDs != nil || len(selector.CreationIDs) > 0 {
		filter = append(filter, bson.E{
			Key: "_id", Value: bson.D{{Key: "$in", Value: selector.CreationIDs}},
		})
	}

	if selector.CreationName != nil {
		filter = append(filter, bson.E{
			Key: "creation_name", Value: selector.CreationName,
		})
	}

	if selector.Creator != nil {
		filter = append(filter, bson.E{
			Key: "creator", Value: selector.Creator,
		})
	}

	if selector.BrandID != nil {
		filter = append(filter, bson.E{
			Key: "brand_id", Value: selector.BrandID,
		})
	}

	if selector.Properties != nil || len(selector.Properties) > 0 {
		filter = append(filter, bson.E{
			Key: "properties", Value: bson.D{{Key: "$in", Value: selector.Properties}},
		})
	}

	if selector.SaleStartBefore != nil || selector.SaleStartAfter != nil {
		saleTimeFilter := bson.D{}
		if selector.SaleStartAfter != nil {
			saleTimeFilter = append(saleTimeFilter, bson.E{Key: "$gte", Value: selector.SaleStartAfter})
		}
		if selector.SaleStartBefore != nil {
			saleTimeFilter = append(saleTimeFilter, bson.E{Key: "$lte", Value: selector.SaleStartBefore})
		}
		filter = append(filter, bson.E{Key: "sale_start_at", Value: saleTimeFilter})
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

type CreationSelector struct {
	CreationIDs     []primitive.ObjectID `json:"creationIDs"`
	CreationName    *string              `json:"creationName"`
	Creator         *string              `json:"creator"`
	Properties      []string             `json:"properties"`
	SaleStartBefore *time.Time           `json:"saleStartBefore"`
	SaleStartAfter  *time.Time           `json:"saleStartAfter"`
	MaxPrice        *int                 `json:"maxPrice"`
	MinPrice        *int                 `json:"minPrice"`
	BrandID         *string              `json:"brandId"`
}

var CreationNotFound = errors.New("creation not found")
