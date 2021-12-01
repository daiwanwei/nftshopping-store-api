package services

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"items"
	"nftshopping-store-api/event/messages"
	"nftshopping-store-api/event/publishers"
	"nftshopping-store-api/persistence/repositories"
	"nftshopping-store-api/pkg/nftshopping"
	"nftshopping-store-api/pkg/utils"
)

type ItemService interface {
	OrderItem(ctx context.Context, dto OrderItemDto) (err error)
	DeliverItem(ctx context.Context, dto DeliverItemDto) (itemDto *ItemDto, err error)
	FindItem(ctx context.Context, contract, token string) (itemDto *ItemDto, err error)
	GetAmountOfItemByBrand(ctx context.Context, brandId string) (amount int64, err error)
	FindAllItemByFilter(ctx context.Context, filter ItemFilterDto) (itemDto []ItemDto, err error)
	FindAllItemByFilterAndPage(ctx context.Context, filter ItemFilterDto, pageable utils.Pageable) (itemDto []ItemDto, err error)
}

type itemService struct {
	creation      CreationService
	brand         BrandService
	user          UserService
	item          repositories.ItemDao
	factory       items.FactoryService
	itemPublisher publishers.ItemPublisher
}

func NewItemService(
	creation CreationService, brand BrandService, user UserService,
) (service ItemService, err error) {
	repository, err := repositories.GetRepository()
	if err != nil {
		return
	}
	item, err := nftshopping.GetItem()
	if err != nil {
		return
	}
	publisher, err := publishers.GetPublisher()
	if err != nil {
		return
	}
	return &itemService{
		brand:         brand,
		creation:      creation,
		user:          user,
		item:          repository.Item,
		factory:       item.Factory,
		itemPublisher: publisher.Item,
	}, nil
}

func (service *itemService) FindItem(
	ctx context.Context, contract, token string,
) (itemDto *ItemDto, err error) {
	itemId := &repositories.ItemID{
		Contract: contract,
		Token:    token,
	}
	item, err := service.item.Find(ctx, itemId)
	if err != nil || item == nil {
		return
	}
	itemDto = &ItemDto{}
	if err = copier.Copy(itemDto, item); err != nil {
		return nil, err
	}
	return
}

func (service *itemService) GetAmountOfItemByBrand(ctx context.Context, brandId string) (amount int64, err error) {
	if isExisted, err := service.brand.Exist(ctx, brandId); err != nil {
		return 0, err
	} else {
		if !isExisted {
			return 0, NewItemServiceError(BrandNotFound)
		}
	}
	amount, err = service.item.CountByBrandOwner(ctx, brandId)
	if err != nil {
		return
	}
	return
}

func (service *itemService) FindAllItemByFilter(
	ctx context.Context, dto ItemFilterDto,
) (itemsDto []ItemDto, err error) {
	selector := repositories.ItemSelector{
		Owner: dto.Owner,
	}
	items, err := service.item.FindAllByFilter(ctx, repositories.SelectorOfItem(selector))
	if err != nil {
		return
	}
	if err = copier.Copy(&itemsDto, &items); err != nil {
		return nil, err
	}
	return
}

func (service *itemService) FindAllItemByFilterAndPage(
	ctx context.Context, dto ItemFilterDto, pageable utils.Pageable,
) (itemsDto []ItemDto, err error) {
	selector := repositories.ItemSelector{
		Owner: dto.Owner,
	}
	page, err := service.item.FindAllByFilterAndPage(ctx, repositories.SelectorOfItem(selector), pageable)
	if err != nil {
		return
	}
	items, ok := page.Content.([]repositories.Item)
	if !ok {
		return nil, utils.ErrCovertContent
	}
	if err = copier.Copy(&itemsDto, &items); err != nil {
		return nil, err
	}
	return
}

func (service *itemService) OrderItem(ctx context.Context, dto OrderItemDto) (err error) {
	creation, err := service.creation.FindCreationByID(ctx, dto.CreationId)
	if err != nil {
		return err
	}
	if creation == nil {
		return NewItemServiceError(CreationNotFound)
	}
	//err = service.factory.OrderItem(items.OrderItemRequest{
	//	ProductId: creation.CreationID,
	//	Contract:  creation.ContractAddress,
	//	Amount:    dto.Amount,
	//})
	//if err != nil {
	//	return
	//}
	err = service.itemPublisher.PublishToOrderItem(messages.OrderItemMessage{
		Contract:   creation.ContractAddress,
		CreationId: creation.CreationID,
		Amount:     dto.Amount,
	})
	if err != nil {
		return
	}
	return
}

func (service *itemService) DeliverItem(
	ctx context.Context, dto DeliverItemDto,
) (itemDto *ItemDto, err error) {
	id, err := primitive.ObjectIDFromHex(dto.CreationId)
	if err != nil {
		return nil, NewItemServiceError(CreationNotFound)
	}
	creation, err := service.creation.FindCreationByID(ctx, dto.CreationId)
	if err != nil {
		return
	}
	if creation == nil {
		return nil, NewItemServiceError(CreationNotFound)
	}
	item := &repositories.Item{
		ID: repositories.ItemID{
			Contract: dto.Contract,
			Token:    dto.Token,
		}, CreationID: id,
		BrandOwner: creation.BrandID,
	}
	err = service.item.Create(ctx, item)
	if err != nil {
		return
	}
	itemDto = &ItemDto{}
	if err = copier.Copy(itemDto, item); err != nil {
		return nil, err
	}
	return
}

type ItemDto struct {
	Contract   string `json:"contract"`
	Token      string `json:"token"`
	Creation   string `json:"creation"`
	Owner      string `json:"owner"`
	BrandOwner string `json:"brandOwner"`
}

func (dto *ItemDto) ID(id repositories.ItemID) {
	dto.Token = id.Token
	dto.Contract = id.Contract
}

func (dto *ItemDto) CreationID(id primitive.ObjectID) {
	dto.Creation = id.Hex()
}

type OrderItemDto struct {
	CreationId string `json:"creationId"`
	Amount     int    `json:"amount"`
}

type DeliverItemDto struct {
	CreationId string `json:"creationId"`
	Contract   string `json:"contract"`
	Token      string `json:"token"`
}

type ItemFilterDto struct {
	Owner      *string `json:"owner"`
	BrandOwner *string `json:"brandOwner"`
}

type ItemServiceError struct {
	ServiceError
}

func NewItemServiceError(e ServiceEvent) error {
	return &ItemServiceError{ServiceError{ServiceName: "ItemService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}
