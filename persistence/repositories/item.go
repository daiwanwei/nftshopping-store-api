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

type ItemDao interface {
	Exist(ctx context.Context, id *ItemID) (isExisted bool, err error)
	Find(ctx context.Context, id *ItemID) (item *Item, err error)
	Create(ctx context.Context, item *Item) (err error)
	Save(ctx context.Context, item *Item) (err error)
	Delete(ctx context.Context, id ItemID) (err error)
	CountByBrandOwner(ctx context.Context, brandOwner string) (amount int64, err error)
	CountByOwner(ctx context.Context, owner string) (amount int64, err error)
	FindAll(ctx context.Context) (items []Item, err error)
	FindAllByPage(ctx context.Context, pageable utils.Pageable) (items *utils.Page, err error)
	FindAllByFilter(ctx context.Context, filter ItemFilter) (items []Item, err error)
	FindAllByFilterAndPage(ctx context.Context, filter ItemFilter, pageable utils.Pageable) (items *utils.Page, err error)
}

type itemDao struct {
	collection *mongo.Collection
}

func NewItemDao() (dao ItemDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &itemDao{db.Collection("creation_item")}, nil
}

func (dao *itemDao) Exist(ctx context.Context, id *ItemID) (isExisted bool, err error) {
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

func (dao *itemDao) Create(ctx context.Context, item *Item) (err error) {
	_, err = dao.collection.InsertOne(ctx, item)
	return
}

func (dao *itemDao) Find(ctx context.Context, id *ItemID) (item *Item, err error) {
	item = &Item{}
	err = dao.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}
	return
}

func (dao *itemDao) Save(ctx context.Context, item *Item) (err error) {
	oldItem := Item{}
	filter := bson.D{{"_id", oldItem.ID}}
	update := bson.D{{"$set", bson.D{
		{"owner", item.Owner},
		{"brand_owner", item.BrandOwner},
	}}}
	option := options.FindOneAndUpdate().SetUpsert(true)
	err = dao.collection.FindOneAndUpdate(ctx, filter, update, option).Decode(&oldItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return
	}
	return
}

func (dao *itemDao) Delete(ctx context.Context, id ItemID) (err error) {
	item := Item{}
	err = dao.collection.FindOneAndDelete(ctx, bson.D{{"_id", id}}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ItemNotFound
		}
		return
	}
	return
}

func (dao *itemDao) CountByBrandOwner(ctx context.Context, brandOwner string) (amount int64, err error) {
	amount, err = dao.collection.CountDocuments(ctx, bson.D{{"brand_owner", brandOwner}})
	if err != nil {
		return
	}
	return
}

func (dao *itemDao) CountByOwner(ctx context.Context, owner string) (amount int64, err error) {
	amount, err = dao.collection.CountDocuments(ctx, bson.D{{"owner", owner}})
	if err != nil {
		return
	}
	return
}

func (dao *itemDao) FindAll(ctx context.Context) (items []Item, err error) {
	items, err = dao.findList(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	return
}

func (dao *itemDao) FindAllByPage(ctx context.Context, pageable utils.Pageable) (items *utils.Page, err error) {
	filter := bson.D{}
	items, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *itemDao) FindAllByFilter(ctx context.Context, filter ItemFilter) (items []Item, err error) {
	items, err = dao.findList(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (dao *itemDao) FindAllByFilterAndPage(
	ctx context.Context, filter ItemFilter, pageable utils.Pageable,
) (items *utils.Page, err error) {
	items, err = dao.findPage(ctx, filter, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *itemDao) findList(ctx context.Context, filter interface{}) (items []Item, err error) {
	cur, err := dao.collection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var item Item
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *itemDao) findPage(
	ctx context.Context, filter interface{}, pageable utils.Pageable,
) (items *utils.Page, err error) {
	total, err := dao.collection.CountDocuments(ctx, filter)
	items = &utils.Page{Size: pageable.Size, Page: pageable.Page, Total: total}
	option := options.Find()
	option.SetSkip(int64(pageable.Size * pageable.Page))
	option.SetLimit(int64(pageable.Size))
	cur, err := dao.collection.Find(ctx, filter, option)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	var content []Item
	for cur.Next(ctx) {
		var item Item
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		content = append(content, item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	items.Content = content

	return
}

type Item struct {
	ID         ItemID             `bson:"_id" json:"id"`
	CreationID primitive.ObjectID `bson:"creation_id" json:"creationId"`
	Owner      string             `bson:"owner" json:"owner"`
	BrandOwner string             `bson:"brand_owner" json:"brandOwner"`
}

type ItemID struct {
	Contract string `bson:"contract" json:"contract"`
	Token    string `bson:"token" json:"token"`
}

type ItemFilter bson.D

func SelectorOfItem(selector ItemSelector) (filter ItemFilter) {
	filter = ItemFilter{}
	if selector.Owner != nil {
		filter = append(filter, bson.E{
			Key: "owner", Value: selector.Owner,
		})
	}
	if selector.BrandOwner != nil {
		filter = append(filter, bson.E{
			Key: "brand_owner", Value: selector.BrandOwner,
		})
	}
	return
}

type ItemSelector struct {
	Owner      *string `json:"owner"`
	BrandOwner *string `json:"brandOwner"`
}

type ItemDetailDao interface {
	Find(ctx context.Context, id *ItemID) (item *ItemDetail, err error)
	FindAll(ctx context.Context) (items []ItemDetail, err error)
	FindAllByPage(ctx context.Context, pageable utils.Pageable) (items *utils.Page, err error)
	FindAllByFilter(ctx context.Context, filter ItemDetailFilter) (items []ItemDetail, err error)
	FindAllByFilterAndPage(ctx context.Context, filter ItemDetailFilter, pageable utils.Pageable) (items *utils.Page, err error)
}

type itemDetailDao struct {
	collection *mongo.Collection
}

func NewItemDetailDao() (dao ItemDetailDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &itemDetailDao{db.Collection("creation_item")}, nil
}

func (dao *itemDetailDao) Find(ctx context.Context, id *ItemID) (item *ItemDetail, err error) {
	item = &ItemDetail{}

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{{"_id", id}}},
		},
	}
	pipeline = append(pipeline, stageOfItemDetail...)
	item, err = dao.findOne(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *itemDetailDao) FindAll(ctx context.Context) (items []ItemDetail, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfItemDetail...)
	items, err = dao.findList(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *itemDetailDao) FindAllByPage(
	ctx context.Context, pageable utils.Pageable,
) (items *utils.Page, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfItemDetail...)
	items, err = dao.findPage(ctx, pipeline, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *itemDetailDao) FindAllByFilter(ctx context.Context, filter ItemDetailFilter) (items []ItemDetail, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfItemDetail...)
	pipeline = append(pipeline, bson.D(filter))
	items, err = dao.findList(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *itemDetailDao) FindAllByFilterAndPage(
	ctx context.Context, filter ItemDetailFilter, pageable utils.Pageable,
) (items *utils.Page, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfItemDetail...)
	pipeline = append(pipeline, bson.D(filter))
	items, err = dao.findPage(ctx, pipeline, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *itemDetailDao) findList(
	ctx context.Context, pipeline []bson.D,
) (items []ItemDetail, err error) {
	cur, err := dao.collection.Aggregate(ctx, pipeline)
	if cur == nil || err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var item ItemDetail
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *itemDetailDao) findPage(
	ctx context.Context, pipeline []bson.D, pageable utils.Pageable,
) (items *utils.Page, err error) {
	items = &utils.Page{Size: pageable.Size, Page: pageable.Page}
	count, err := dao.count(context.Background(), pipeline)
	if err != nil {
		return
	}
	items.Total = count

	pipelineOfPage := mongo.Pipeline{}
	pageStage := []bson.D{
		{{"$limit", pageable.Size}},
		{{"$skip", pageable.Size * pageable.Page}},
	}
	if len(pageable.Sort) > 0 {
		sort := bson.D{}
		for key, value := range pageable.Sort {
			sort = append(sort, bson.E{Key: key, Value: value})
		}
		pageStage = append(pageStage, bson.D{{"$sort", sort}})
	}

	pipelineOfPage = append(pipelineOfPage, pipeline...)
	pipelineOfPage = append(pipelineOfPage, pageStage...)
	cur, err := dao.collection.Aggregate(ctx, pipelineOfPage)
	if cur == nil || err != nil {
		return
	}
	defer cur.Close(ctx)
	var content []ItemDetail
	for cur.Next(ctx) {
		var transaction ItemDetail
		err := cur.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		content = append(content, transaction)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	items.Content = content
	return
}

func (dao *itemDetailDao) findOne(ctx context.Context, pipeline []bson.D) (item *ItemDetail, err error) {
	if err != nil {
		return
	}

	pipelineOfPage := mongo.Pipeline{}
	oneStage := []bson.D{
		{{"$limit", 1}},
	}
	pipelineOfPage = append(pipelineOfPage, pipeline...)
	pipelineOfPage = append(pipelineOfPage, oneStage...)
	cur, err := dao.collection.Aggregate(ctx, pipelineOfPage)
	if cur == nil || err != nil {
		return
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		err = cur.Decode(&item)
		if err != nil {
			return nil, err
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *itemDetailDao) count(ctx context.Context, pipeline []bson.D) (count int64, err error) {
	pipelineOfCount := mongo.Pipeline{}
	countStage := bson.D{
		{"$count", "count"},
	}
	pipelineOfCount = append(pipelineOfCount, pipeline...)
	pipelineOfCount = append(pipelineOfCount, countStage)
	curOfCount, err := dao.collection.Aggregate(ctx, pipelineOfCount)
	if curOfCount == nil || err != nil {
		return
	}
	defer curOfCount.Close(ctx)
	var info []bson.M
	if err = curOfCount.All(ctx, &info); err != nil {
		return
	}
	if len(info) < 1 {
		return 0, nil
	}
	return int64(info[0]["count"].(int32)), nil
}

type ItemDetail struct {
	ID         ItemID   `bson:"_id" json:"id"`
	Owner      string   `bson:"owner" json:"owner"`
	BrandOwner string   `bson:"brand_owner" json:"brandOwner"`
	Creation   Creation `bson:"creation" json:"creation"`
}

var stageOfItemDetail = []bson.D{
	{{"$lookup", bson.D{
		{"from", "creation"},
		{"localField", "creation_id"},
		{"foreignField", "_id"},
		{"as", "creation"},
	}}},
	{{"$unwind", "$creation"}},
	{{"$project", bson.D{
		{"creation_id", 0},
	}}},
}

type ItemDetailFilter bson.D

func SelectorOfItemDetail(selector ItemDetailSelector) (filter ItemDetailFilter) {
	matchStage := bson.D{}
	if selector.Owner != nil {
		matchStage = append(matchStage, bson.E{
			Key: "owner", Value: selector.Owner,
		})
	}
	if selector.BrandOwner != nil {
		matchStage = append(matchStage, bson.E{
			Key: "brand_owner", Value: selector.BrandOwner,
		})
	}
	return ItemDetailFilter(bson.D{{"$match", matchStage}})
}

type ItemDetailSelector struct {
	Owner      *string `json:"owner"`
	BrandOwner *string `json:"brandOwner"`
}

var ItemNotFound = errors.New("item not found")
