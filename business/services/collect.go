package services

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nftshopping-store-api/persistence/repositories"
	"nftshopping-store-api/pkg/utils"
)

type CollectService interface {
	FindCollect(ctx context.Context, owner, creationId string) (collection *CollectDto, err error)
	FindAllCollect(ctx context.Context) (collections []CollectDto, err error)
	FindAllCollectByFilter(
		ctx context.Context, filter CollectFilterDto,
	) (collectionsDto []CollectDto, err error)
	FindAllCollectByFilterAndPage(
		ctx context.Context, filter CollectFilterDto, pageable utils.Pageable,
	) (collectionsDto []CollectDto, err error)
}

type CollectDto struct {
	CreationID string `json:"creationId"`
	Owner      string `json:"owner"`
	Amount     int    `json:"amount"`
}

func (dto *CollectDto) ID(id repositories.CollectID) {
	dto.CreationID = id.CreationID.Hex()
	dto.Owner = id.Owner
}

type CollectFilterDto struct {
	Owner      *string `json:"owner"`
	CreationID *string `json:"creationId"`
}

type CollectionService interface {
	CollectService
}

type collectionService struct {
	user       UserService
	creation   CreationService
	collection repositories.CollectionDao
}

func NewCollectionService(user UserService, creation CreationService) (service CollectionService, err error) {
	dao, err := repositories.GetRepository()
	if err != nil {
		return nil, err
	}
	return &collectionService{
		user:       user,
		creation:   creation,
		collection: dao.Collection,
	}, nil
}

func (service *collectionService) FindCollect(
	ctx context.Context, ownerName, creationId string,
) (collectionDto *CollectDto, err error) {
	id, err := primitive.ObjectIDFromHex(creationId)
	if err != nil {
		return nil, nil
	}
	collectionId := &repositories.CollectID{
		Owner:      ownerName,
		CreationID: id,
	}
	collection, err := service.collection.Find(ctx, collectionId)
	if err != nil || collection == nil {
		return
	}
	collectionDto = &CollectDto{}
	if err = copier.Copy(collectionDto, collection); err != nil {
		return nil, err
	}
	return
}

func (service *collectionService) FindAllCollect(ctx context.Context) (collectionsDto []CollectDto, err error) {
	collections, err := service.collection.FindAll(ctx)
	if err != nil {
		return
	}
	if err = copier.Copy(&collectionsDto, &collections); err != nil {
		return nil, err
	}
	return
}

func (service *collectionService) FindAllCollectByFilter(
	ctx context.Context, filter CollectFilterDto,
) (collectionsDto []CollectDto, err error) {
	selector := repositories.CollectSelector{
		Owner: filter.Owner,
	}
	if filter.CreationID != nil {
		if creationId, err := primitive.ObjectIDFromHex(*filter.CreationID); err != nil {
			return nil, err
		} else {
			selector.CreationID = &creationId
		}
	}
	collections, err := service.collection.FindAllByFilter(
		ctx, repositories.SelectorOfCollect(selector),
	)
	if err != nil {
		return
	}
	if err = copier.Copy(&collectionsDto, &collections); err != nil {
		return nil, err
	}
	return
}

func (service *collectionService) FindAllCollectByFilterAndPage(
	ctx context.Context, filter CollectFilterDto, pageable utils.Pageable,
) (collectionsDto []CollectDto, err error) {
	selector := repositories.CollectSelector{
		Owner: filter.Owner,
	}
	if filter.CreationID != nil {
		if creationId, err := primitive.ObjectIDFromHex(*filter.CreationID); err != nil {
			return nil, err
		} else {
			selector.CreationID = &creationId
		}
	}
	page, err := service.collection.FindAllByFilterAndPage(
		ctx, repositories.SelectorOfCollect(selector), pageable,
	)
	if err != nil {
		return
	}
	if collections, ok := page.Content.([]repositories.Collect); !ok {
		return nil, utils.ErrCovertContent
	} else {
		err = copier.Copy(&collectionsDto, &collections)
		if err != nil {
			return nil, err
		}
	}
	return
}

type CollectionServiceError struct {
	ServiceError
}

func NewCollectServiceError(e ServiceEvent) error {
	return &CollectionServiceError{ServiceError{ServiceName: "CollectionService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}

type StockService interface {
	CollectService
}

type stockService struct {
	stock    repositories.StockDao
	creation CreationService
	brand    BrandService
}

func NewStockService(brand BrandService, creation CreationService) (service StockService, err error) {
	dao, err := repositories.GetRepository()
	if err != nil {
		return nil, err
	}
	return &stockService{
		stock:    dao.Stock,
		creation: creation,
		brand:    brand,
	}, nil
}

func (service *stockService) FindCollect(
	ctx context.Context, ownerName, creationId string,
) (collect *CollectDto, err error) {
	id, err := primitive.ObjectIDFromHex(creationId)
	if err != nil {
		return nil, nil
	}
	stockId := &repositories.CollectID{
		Owner:      ownerName,
		CreationID: id,
	}
	stock, err := service.stock.Find(ctx, stockId)
	if err != nil || stock == nil {
		return
	}
	collect = &CollectDto{}
	if err = copier.Copy(collect, stock); err != nil {
		return nil, err
	}
	return
}

func (service *stockService) FindAllCollect(ctx context.Context) (collectDto []CollectDto, err error) {
	collect, err := service.stock.FindAll(ctx)
	if err != nil {
		return
	}
	if err = copier.Copy(&collectDto, &collect); err != nil {
		return nil, err
	}
	return
}

func (service *stockService) FindAllCollectByFilter(
	ctx context.Context, filter CollectFilterDto,
) (collectDto []CollectDto, err error) {
	return
}

func (service *stockService) FindAllCollectByFilterAndPage(
	ctx context.Context, filter CollectFilterDto, pageable utils.Pageable,
) (collectDto []CollectDto, err error) {
	selector := repositories.CollectSelector{
		Owner: filter.Owner,
	}
	if filter.CreationID != nil {
		if creationId, err := primitive.ObjectIDFromHex(*filter.CreationID); err != nil {
			return nil, err
		} else {
			selector.CreationID = &creationId
		}
	}
	collect, err := service.stock.FindAllByFilter(
		ctx, repositories.SelectorOfCollect(selector),
	)
	if err != nil {
		return
	}
	if err = copier.Copy(&collectDto, &collect); err != nil {
		return nil, err
	}
	return
}

type StockServiceError struct {
	ServiceError
}

func NewStockServiceError(e ServiceEvent) error {
	return &StockServiceError{ServiceError{ServiceName: "StockService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}
