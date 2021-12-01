package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"nftshopping-store-api/pkg/databases"
	"nftshopping-store-api/pkg/utils"
	"time"
)

type BrandDao interface {
	Exist(ctx context.Context, id string) (isExisted bool, err error)
	Find(ctx context.Context, id string) (brand *Brand, err error)
	Save(ctx context.Context, brand *Brand) (err error)
	Delete(ctx context.Context, id string) (err error)
	Create(ctx context.Context, brand *Brand) (err error)
	FindAll(ctx context.Context) (brands []Brand, err error)
	FindAllByPage(ctx context.Context, pageable utils.Pageable) (brands *utils.Page, err error)
	FindAllByFilter(
		ctx context.Context, filter BrandFilter,
	) (brands []Brand, err error)
	FindAllByFilterAndPage(
		ctx context.Context, filter BrandFilter, pageable utils.Pageable,
	) (brands *utils.Page, err error)
}

type brandDao struct {
	collection *mongo.Collection
}

func NewBrandDao() (dao BrandDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &brandDao{db.Collection("brand")}, nil
}

func (dao *brandDao) Exist(ctx context.Context, id string) (isExisted bool, err error) {
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

func (dao *brandDao) Create(ctx context.Context, brand *Brand) (err error) {
	_, err = dao.collection.InsertOne(ctx, brand)
	return
}

func (dao *brandDao) Find(ctx context.Context, id string) (brand *Brand, err error) {
	brand = &Brand{}
	err = dao.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(brand)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}
	return
}

func (dao *brandDao) Save(ctx context.Context, brand *Brand) (err error) {
	oldBrand := Brand{}
	filter := bson.D{{"_id", brand.ID}}
	update := bson.D{{"$set", bson.D{
		{"name", brand.Name},
		{"image_url", brand.ImageURL},
		{"create_at", brand.CreateAt},
		{"description", brand.Description},
	}}}
	option := options.FindOneAndUpdate().SetUpsert(true)
	err = dao.collection.FindOneAndUpdate(ctx, filter, update, option).Decode(&oldBrand)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return
	}
	return
}

func (dao *brandDao) Delete(ctx context.Context, id string) (err error) {
	brand := Brand{}
	err = dao.collection.FindOneAndDelete(ctx, bson.D{{"_id", id}}).Decode(&brand)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return BrandNotFound
		}
		return
	}
	return
}

func (dao *brandDao) FindAll(ctx context.Context) (brands []Brand, err error) {
	brands, err = dao.findList(ctx, bson.D{{}})
	if err != nil {
		return
	}
	return
}

func (dao *brandDao) FindAllByPage(
	ctx context.Context, pageable utils.Pageable,
) (brands *utils.Page, err error) {
	brands, err = dao.findPage(ctx, bson.D{{}}, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *brandDao) FindAllByFilter(
	ctx context.Context, filter BrandFilter,
) (brands []Brand, err error) {
	brands, err = dao.findList(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (dao *brandDao) FindAllByFilterAndPage(
	ctx context.Context, filter BrandFilter, pageable utils.Pageable,
) (brands *utils.Page, err error) {
	brands, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *brandDao) findList(ctx context.Context, filter interface{}) (brands []Brand, err error) {
	cur, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var brand Brand
		err := cur.Decode(&brand)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *brandDao) findPage(
	ctx context.Context, filter interface{}, pageable utils.Pageable,
) (brands *utils.Page, err error) {
	total, err := dao.collection.CountDocuments(ctx, filter)
	brands = &utils.Page{Size: pageable.Size, Page: pageable.Page, Total: total}
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
	var content []Brand
	for cur.Next(ctx) {
		var brand Brand
		err := cur.Decode(&brand)
		if err != nil {
			return nil, err
		}
		content = append(content, brand)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	brands.Content = content
	brands.TotalPage = utils.GetTotalPage(int64(brands.Size), brands.Total)

	return
}

type Brand struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	ImageURL    string    `bson:"image_url" json:"imageUrl"`
	CreateAt    time.Time `bson:"create_at" json:"createAt"`
}

type BrandFilter bson.D

func SelectorOfBrand(selector BrandSelector) (filter BrandFilter) {
	filter = BrandFilter{}
	if selector.Name != nil {
		filter = append(filter, bson.E{
			Key: "name", Value: selector.Name,
		})
	}

	if selector.CreateBefore != nil || selector.CreateAfter != nil {
		intervalFilter := bson.A{}
		if selector.CreateAfter != nil {
			intervalFilter = append(intervalFilter, bson.D{{"$lte", bson.A{selector.CreateAfter, "$create_at"}}})
		}
		if selector.CreateBefore != nil {
			intervalFilter = append(intervalFilter, bson.D{{"$gte", bson.A{selector.CreateBefore, "$create_at"}}})
		}
		filter = append(filter, bson.E{Key: "$expr", Value: bson.D{{"$and", intervalFilter}}})
	}
	return
}

type BrandSelector struct {
	Name         *string    `json:"name"`
	CreateBefore *time.Time `json:"createBefore"`
	CreateAfter  *time.Time `json:"createAfter"`
}

var BrandNotFound = errors.New("brand not found")
