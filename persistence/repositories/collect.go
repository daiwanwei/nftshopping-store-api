package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nftshopping-store-api/pkg/databases"
	"nftshopping-store-api/pkg/utils"
)

type CollectDao interface {
	Find(ctx context.Context, id *CollectID) (collection *Collect, err error)
	FindAll(ctx context.Context) (collections []Collect, err error)
	FindAllByPage(ctx context.Context, pageable utils.Pageable) (collects *utils.Page, err error)
	FindAllByFilter(ctx context.Context, filter CollectFilter) (collects []Collect, err error)
	FindAllByFilterAndPage(ctx context.Context, filter CollectFilter, pageable utils.Pageable) (collects *utils.Page, err error)
}

type Collect struct {
	ID     CollectID `bson:"_id" json:"id"`
	Amount int       `bson:"amount" json:"amount"`
}

type CollectID struct {
	Owner      string             `bson:"owner" json:"owner"`
	CreationID primitive.ObjectID `bson:"creation_id" json:"creationId"`
}

type CollectFilter bson.D

func SelectorOfCollect(selector CollectSelector) (filter CollectFilter) {
	matchStage := bson.D{}
	if selector.Owner != nil {
		matchStage = append(matchStage, bson.E{
			Key: "_id.owner", Value: selector.Owner,
		})
	}
	if selector.CreationID != nil {
		matchStage = append(matchStage, bson.E{
			Key: "_id.creation_id", Value: selector.CreationID,
		})
	}

	return CollectFilter(bson.D{{"$match", matchStage}})
}

type CollectSelector struct {
	Owner      *string             `json:"owner"`
	CreationID *primitive.ObjectID `json:"creationId"`
}

type CollectionDao interface {
	CollectDao
}

type collectionDao struct {
	collection *mongo.Collection
}

func NewCollectionDao() (dao CollectDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &collectionDao{db.Collection("creation_item")}, nil
}

func (dao *collectionDao) Find(ctx context.Context, id *CollectID) (collection *Collect, err error) {
	collection = &Collect{}
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"owner", id.Owner},
				{"creation_id", id.CreationID},
			}},
		},
	}
	pipeline = append(pipeline, stageOfCollection...)
	collection, err = dao.findOne(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *collectionDao) FindAll(ctx context.Context) (items []Collect, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfCollection...)
	items, err = dao.findList(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *collectionDao) FindAllByPage(
	ctx context.Context, pageable utils.Pageable,
) (items *utils.Page, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfCollection...)
	items, err = dao.findPage(ctx, pipeline, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *collectionDao) FindAllByFilter(ctx context.Context, filter CollectFilter) (collections []Collect, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfCollection...)
	pipeline = append(pipeline, bson.D(filter))
	collections, err = dao.findList(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *collectionDao) FindAllByFilterAndPage(
	ctx context.Context, filter CollectFilter, pageable utils.Pageable,
) (collections *utils.Page, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfCollection...)
	pipeline = append(pipeline, bson.D(filter))
	collections, err = dao.findPage(ctx, pipeline, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *collectionDao) findOne(ctx context.Context, pipeline []bson.D) (collection *Collect, err error) {
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
		err = cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *collectionDao) findList(
	ctx context.Context, pipeline []bson.D,
) (collections []Collect, err error) {
	cur, err := dao.collection.Aggregate(ctx, pipeline)
	if cur == nil || err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var collection Collect
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *collectionDao) findPage(
	ctx context.Context, pipeline []bson.D, pageable utils.Pageable,
) (collections *utils.Page, err error) {
	collections = &utils.Page{Size: pageable.Size, Page: pageable.Page}
	count, err := dao.count(context.Background(), pipeline)
	if err != nil {
		return
	}
	collections.Total = count

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
	var content []Collect
	for cur.Next(ctx) {
		var collection Collect
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		content = append(content, collection)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	collections.Content = content
	return
}

func (dao *collectionDao) count(ctx context.Context, pipeline []bson.D) (count int64, err error) {
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

var stageOfCollection = []bson.D{
	{
		{"$group", bson.D{
			{"_id", bson.D{
				{"owner", "$owner"},
				{"creation_id", "$creation_id"},
			}},
			{"amount", bson.D{{"$sum", 1}}},
		}},
	},
}

var CollectionNotFound = errors.New("collection not found")

type StockDao interface {
	CollectDao
}

type stockDao struct {
	stock *mongo.Collection
}

func NewStockDao() (dao StockDao, err error) {
	db, err := databases.GetMongoDB()
	if err != nil {
		return nil, err
	}
	return &stockDao{db.Collection("creation_item")}, nil
}

func (dao *stockDao) Find(ctx context.Context, id *CollectID) (collection *Collect, err error) {
	collection = &Collect{}
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"brand_owner", id.Owner},
				{"creation_id", id.CreationID},
			}},
		},
	}
	pipeline = append(pipeline, stageOfStock...)
	collection, err = dao.findOne(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *stockDao) FindAll(ctx context.Context) (collections []Collect, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfStock...)
	collections, err = dao.findList(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *stockDao) FindAllByPage(
	ctx context.Context, pageable utils.Pageable,
) (collections *utils.Page, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfStock...)
	collections, err = dao.findPage(ctx, pipeline, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *stockDao) FindAllByFilter(ctx context.Context, filter CollectFilter) (collections []Collect, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfStock...)
	pipeline = append(pipeline, bson.D(filter))
	collections, err = dao.findList(ctx, pipeline)
	if err != nil {
		return
	}
	return
}

func (dao *stockDao) FindAllByFilterAndPage(
	ctx context.Context, filter CollectFilter, pageable utils.Pageable,
) (collections *utils.Page, err error) {
	pipeline := mongo.Pipeline{}
	pipeline = append(pipeline, stageOfStock...)
	pipeline = append(pipeline, bson.D(filter))
	collections, err = dao.findPage(ctx, pipeline, pageable)
	if err != nil {
		return
	}
	return
}

func (dao *stockDao) findOne(ctx context.Context, pipeline []bson.D) (collection *Collect, err error) {
	if err != nil {
		return
	}
	pipelineOfPage := mongo.Pipeline{}
	oneStage := []bson.D{
		{{"$limit", 1}},
	}
	pipelineOfPage = append(pipelineOfPage, pipeline...)
	pipelineOfPage = append(pipelineOfPage, oneStage...)
	cur, err := dao.stock.Aggregate(ctx, pipelineOfPage)
	if cur == nil || err != nil {
		return
	}
	defer cur.Close(ctx)
	if cur.Next(ctx) {
		err = cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *stockDao) findList(
	ctx context.Context, pipeline []bson.D,
) (collections []Collect, err error) {
	cur, err := dao.stock.Aggregate(ctx, pipeline)
	if cur == nil || err != nil {
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var collection Collect
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return
}

func (dao *stockDao) findPage(
	ctx context.Context, pipeline []bson.D, pageable utils.Pageable,
) (collections *utils.Page, err error) {
	collections = &utils.Page{Size: pageable.Size, Page: pageable.Page}
	count, err := dao.count(context.Background(), pipeline)
	if err != nil {
		return
	}
	collections.Total = count

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
	cur, err := dao.stock.Aggregate(ctx, pipelineOfPage)
	if cur == nil || err != nil {
		return
	}
	defer cur.Close(ctx)
	var content []Collect
	for cur.Next(ctx) {
		var collection Collect
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		content = append(content, collection)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	collections.Content = content
	collections.TotalPage = utils.GetTotalPage(int64(collections.Size), collections.Total)
	return
}

func (dao *stockDao) count(ctx context.Context, pipeline []bson.D) (count int64, err error) {
	pipelineOfCount := mongo.Pipeline{}
	countStage := bson.D{
		{"$count", "count"},
	}
	pipelineOfCount = append(pipelineOfCount, pipeline...)
	pipelineOfCount = append(pipelineOfCount, countStage)
	curOfCount, err := dao.stock.Aggregate(ctx, pipelineOfCount)
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

var stageOfStock = []bson.D{
	{
		{"$group", bson.D{
			{"_id", bson.D{
				{"owner", "$brand_owner"},
				{"creation_id", "$creation_id"},
			}},
			{"amount", bson.D{{"$sum", 1}}},
		}},
	},
}

var StockNotFound = errors.New("stock not found")
