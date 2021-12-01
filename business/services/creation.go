package services

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"items"
	"nftshopping-store-api/persistence/repositories"
	"nftshopping-store-api/pkg/caches"
	"nftshopping-store-api/pkg/nftshopping"
	"nftshopping-store-api/pkg/utils"
	"time"
)

type CreationService interface {
	Exist(ctx context.Context, id string) (isExisted bool, err error)
	FindCreationByID(ctx context.Context, id string) (creationDto *CreationDto, err error)
	FindAllCreationByFilter(ctx context.Context, dto CreationFilterDto) (creationsDto []CreationDto, err error)
	FindAllCreationByFilterAndPage(ctx context.Context, dto CreationFilterDto, pageable utils.Pageable) (creationsDto []CreationDto, err error)
	PostCreation(ctx context.Context, dto PostCreationDto) (creationDto *CreationDto, err error)
	DeleteCreation(ctx context.Context, id primitive.ObjectID) (err error)
	UpdateCreation(ctx context.Context, dto UpdateCreationDto) (err error)
}

type creationService struct {
	creation          repositories.CreationDao
	creationCache     *lru.Cache
	creationListCache *lru.Cache
	brand             BrandService
	contractManager   items.ContractManagerService
}

func NewCreationService(brand BrandService) (service CreationService, err error) {
	dao, err := repositories.GetRepository()
	if err != nil {
		return nil, err
	}
	cacheManager, err := caches.GetCacheManager()
	if err != nil {
		return
	}
	creationCache, err := cacheManager.GetCache("creation")
	if err != nil {
		return
	}
	creationListCache, err := cacheManager.GetCache("creationListCache")
	if err != nil {
		return
	}
	item, err := nftshopping.GetItem()
	if err != nil {
		return
	}
	return &creationService{
		creation:          dao.Creation,
		creationCache:     creationCache,
		creationListCache: creationListCache,
		brand:             brand,
		contractManager:   item.ContractManager,
	}, nil
}

func (service *creationService) Exist(ctx context.Context, id string) (isExisted bool, err error) {
	creationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, nil
	}
	isExisted, err = service.creation.Exist(ctx, creationId)
	return
}

func (service *creationService) FindCreationByID(ctx context.Context, id string) (creationDto *CreationDto, err error) {
	creationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}
	creation, err := service.creation.Find(ctx, creationId)
	if err != nil {
		return
	}
	creationDto = &CreationDto{}
	if err = copier.Copy(creationDto, creation); err != nil {
		return nil, err
	}
	return
}

func (service *creationService) FindAllCreationByFilter(ctx context.Context, dto CreationFilterDto) (creationsDto []CreationDto, err error) {
	//creationsDto,ok:=service.getCreationListFromCache(ctx, dto)
	//if ok{
	//	return
	//}
	var creationIds []primitive.ObjectID
	for _, id := range dto.CreationIDs {
		creationId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		creationIds = append(creationIds, creationId)
	}
	selector := repositories.CreationSelector{
		CreationIDs:     creationIds,
		CreationName:    dto.CreationName,
		Properties:      dto.Properties,
		SaleStartAfter:  dto.SaleEndAfter,
		SaleStartBefore: dto.SaleStartBefore,
		MaxPrice:        dto.MaxPrice,
		MinPrice:        dto.MinPrice,
		BrandID:         dto.BrandID,
	}
	creations, err := service.creation.FindAllByFilter(ctx, repositories.SelectorOfCreation(selector))
	if err != nil {
		return
	}
	if err = copier.Copy(&creationsDto, &creations); err != nil {
		return nil, err
	}
	//service.addCreationListToCache(ctx, dto,creationsDto)
	return
}

func (service *creationService) FindAllCreationByFilterAndPage(ctx context.Context, dto CreationFilterDto, pageable utils.Pageable) (creationsDto []CreationDto, err error) {
	//creationsDto,ok:=service.getCreationListFromCache(ctx, dto)
	//if ok{
	//	return
	//}
	var creationIds []primitive.ObjectID
	for _, id := range dto.CreationIDs {
		creationId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		creationIds = append(creationIds, creationId)
	}
	selector := repositories.CreationSelector{
		CreationIDs:     creationIds,
		CreationName:    dto.CreationName,
		Creator:         dto.Creator,
		Properties:      dto.Properties,
		SaleStartAfter:  dto.SaleEndAfter,
		SaleStartBefore: dto.SaleStartBefore,
		MaxPrice:        dto.MaxPrice,
		MinPrice:        dto.MinPrice,
		BrandID:         dto.BrandID,
	}
	page, err := service.creation.FindAllByFilterAndPage(ctx, repositories.SelectorOfCreation(selector), pageable)
	if err != nil {
		return
	}
	creations, ok := page.Content.([]repositories.Creation)
	if !ok {
		return nil, utils.ErrCovertContent
	}
	if err = copier.Copy(&creationsDto, &creations); err != nil {
		return nil, err
	}
	//service.addCreationListToCache(ctx, dto,creationsDto)
	return
}

func (service *creationService) PostCreation(ctx context.Context, dto PostCreationDto) (creationDto *CreationDto, err error) {
	if isExisted, err := service.brand.Exist(ctx, dto.BrandID); err != nil {
		return nil, err
	} else {
		if !isExisted {
			return nil, NewBrandServiceError(BrandNotFound)
		}
	}
	creation := &repositories.Creation{}
	if err = copier.Copy(creation, &dto); err != nil {
		return
	}
	creation.ID = primitive.NewObjectID()
	creation.CreateAt = time.Now()
	creation.SaleStatus = "WAIT_FOR_SALE"
	creation.BrandID = dto.BrandID
	err = service.creation.Create(ctx, creation)
	if err != nil {
		return
	}
	creationDto = &CreationDto{}
	err = copier.Copy(creationDto, creation)
	if err != nil {
		return
	}
	err = service.contractManager.DeployCreation(items.DeployCreationRequest{
		Contract: dto.ContractAddress,
	})
	if err != nil {
		if respErr, ok := err.(*items.ItemError); ok {
			switch respErr.Code {
			case int(items.ContractExisted):
				return nil, NewCreationServiceError(ContractDuplicate)
			default:
				return nil, err
			}
		}
		return
	}
	return
}

func (service *creationService) DeleteCreation(ctx context.Context, id primitive.ObjectID) (err error) {
	err = service.creation.Delete(ctx, id)
	if err != nil {
		if err != repositories.CreationNotFound {
			return nil
		}
	}
	return
}

func (service *creationService) getCreationListFromCache(ctx context.Context, dto CreationFilterDto) (creationsDto []CreationDto, ok bool) {
	key, err := generateKeyOfCache(dto)
	if err != nil {
		return nil, false
	}
	val, ok := service.creationListCache.Get(key)
	if !ok {
		return nil, false
	}
	creationsDto, ok = val.([]CreationDto)
	if !ok {
		return nil, false
	}
	return
}

func (service *creationService) addCreationListToCache(ctx context.Context, dto CreationFilterDto, creationsDto []CreationDto) (ok bool) {
	key, err := generateKeyOfCache(dto)
	if err != nil {
		return false
	}
	ok = service.creationListCache.Add(key, creationsDto)
	if !ok {
		return
	}
	return
}

func (service *creationService) UpdateCreation(ctx context.Context, dto UpdateCreationDto) (err error) {
	id, err := primitive.ObjectIDFromHex(dto.CreationID)
	if err != nil {
		return NewCreationServiceError(CreationNotFound)
	}
	creation, err := service.creation.Find(ctx, id)
	if err != nil {
		return
	}
	if creation == nil {
		return NewCreationServiceError(CreationNotFound)
	}
	err = copier.Copy(creation, dto)
	if err != nil {
		return
	}
	err = service.creation.Save(ctx, creation)
	if err != nil {
		return
	}
	return
}

type CreationDto struct {
	CreationID      string    `json:"creationId"`
	CreationName    string    `json:"creationName"`
	Amount          int       `json:"amount"`
	SmallImageURL   string    `json:"smallImageUrl"`
	Creator         string    `json:"creator"`
	Properties      []string  `json:"properties"`
	Price           int       `json:"price"`
	BrandID         string    `json:"brandId"`
	SaleWay         string    `json:"saleWay"`
	SaleStatus      string    `json:"saleStatus"`
	SaleStartAt     time.Time `json:"saleStartAt"`
	SaleEndAt       time.Time `json:"saleEndAt"`
	Description     string    `json:"description"`
	ContractAddress string    `json:"contractAddress"`
}

func (dto *CreationDto) ID(id primitive.ObjectID) {
	dto.CreationID = id.Hex()
}

type CreationFilterDto struct {
	CreationIDs     []string   `json:"creationIDs"`
	CreationName    *string    `json:"creationName"`
	Properties      []string   `json:"properties"`
	Creator         *string    `json:"creator"`
	BrandID         *string    `json:"brandId"`
	SaleStartBefore *time.Time `json:"saleStartBefore"`
	SaleEndAfter    *time.Time `json:"saleEndAfter"`
	MaxPrice        *int       `json:"maxPrice"`
	MinPrice        *int       `json:"minPrice"`
}

type PostCreationDto struct {
	CreationName    string    `json:"creationName"`
	Creator         string    `json:"creator"`
	SmallImageURL   string    `json:"smallImageUrl"`
	Amount          int       `json:"amount"`
	Price           int       `bson:"price" json:"price"`
	Properties      []string  `json:"properties"`
	BrandID         string    `json:"brandId"`
	SaleWay         string    `bson:"sale_way" json:"saleWay"`
	SaleStatus      string    `bson:"sale_status" json:"saleStatus"`
	SaleStartAt     time.Time `bson:"sale_start_at" json:"saleStartAt"`
	SaleEndAt       time.Time `bson:"sale_end_at" json:"saleEndAt"`
	Description     string    `json:"description"`
	ContractAddress string    `json:"contractAddress"`
}

type UpdateCreationDto struct {
	CreationID   string    `json:"creationId"`
	SaleStartAt  time.Time `json:"saleStartAt"`
	SaleEndAt    time.Time `json:"saleEndAt"`
	CreationName string    `json:"creationName"`
	Description  string    `json:"description"`
}

type CreationServiceError struct {
	ServiceError
}

func NewCreationServiceError(e ServiceEvent) error {
	return &CreationServiceError{ServiceError{ServiceName: "CreationService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}
