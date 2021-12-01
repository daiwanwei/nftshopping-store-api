package services

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nftshopping-store-api/persistence/repositories"
	"nftshopping-store-api/pkg/utils"
	"time"
)

type TradeService interface {
	FindTransaction(ctx context.Context, tnxId string) (transactionDto *TransactionDto, err error)
	FindAllTransactionByFilter(
		ctx context.Context, dto TransactionFilterDto,
	) (transactionsDto []TransactionDto, err error)
	FindAllTransactionByFilterAndPage(
		ctx context.Context, dto TransactionFilterDto, pageable utils.Pageable,
	) (transactionsDto []TransactionDto, err error)
	TradeInCreation(ctx context.Context, dto TradeInCreationDto) (txn *TransactionDto, err error)
}

type tradeService struct {
	user        UserService
	creation    CreationService
	transaction repositories.TransactionDao
}

func NewTradeService(user UserService, creation CreationService) (service TradeService, err error) {
	dao, err := repositories.GetRepository()
	if err != nil {
		return nil, err
	}
	return &tradeService{
		transaction: dao.Transaction,
		user:        user,
		creation:    creation,
	}, nil
}

func (service *tradeService) FindTransaction(ctx context.Context, tnxId string) (transactionDto *TransactionDto, err error) {
	id, err := primitive.ObjectIDFromHex(tnxId)
	if err != nil {
		return nil, nil
	}
	transaction, err := service.transaction.Find(ctx, id)
	if err != nil {
		return
	}
	transactionDto = &TransactionDto{}
	if err = copier.Copy(transactionDto, transaction); err != nil {
		return nil, err
	}
	return
}

func (service *tradeService) FindAllTransactionByFilter(
	ctx context.Context, dto TransactionFilterDto,
) (transactionsDto []TransactionDto, err error) {
	selector := repositories.TransactionSelector{}
	err = copier.Copy(&selector, &dto)
	if err != nil {
		return
	}

	creations, err := service.transaction.FindAllByFilter(ctx, repositories.SelectorOfTransaction(selector))
	if err != nil {
		return
	}
	if err = copier.Copy(&transactionsDto, &creations); err != nil {
		return nil, err
	}
	return
}

func (service *tradeService) FindAllTransactionByFilterAndPage(
	ctx context.Context, dto TransactionFilterDto, pageable utils.Pageable,
) (transactionsDto []TransactionDto, err error) {
	selector := repositories.TransactionSelector{}
	err = copier.Copy(&selector, &dto)
	if err != nil {
		return
	}
	page, err := service.transaction.FindAllByFilterAndPage(ctx, repositories.SelectorOfTransaction(selector), pageable)
	if err != nil {
		return
	}
	creations, ok := page.Content.([]repositories.Transaction)
	if !ok {
		return nil, utils.ErrCovertContent
	}
	if err = copier.Copy(&transactionsDto, &creations); err != nil {
		return nil, err
	}
	return
}

func (service *tradeService) TradeInCreation(
	ctx context.Context, dto TradeInCreationDto,
) (txn *TransactionDto, err error) {
	if isExisted, err := service.user.ExistByID(ctx, dto.Buyer); err != nil {
		return nil, err
	} else {
		if !isExisted {
			return nil, NewTradeServiceError(UserNotFound)
		}
	}

	if isExisted, err := service.user.ExistByID(ctx, dto.Seller); err != nil {
		return nil, err
	} else {
		if !isExisted {
			return nil, NewTradeServiceError(UserNotFound)
		}
	}
	creation, err := service.creation.FindCreationByID(ctx, dto.CreationID)
	if err != nil {
		return
	}
	if creation == nil {
		return nil, NewTradeServiceError(CreationNotFound)
	}

	transaction := &repositories.Transaction{}
	if err = copier.Copy(transaction, &dto); err != nil {
		return
	}
	transaction.ID = primitive.NewObjectID()
	transaction.CreationID = dto.CreationID
	transaction.BrandID = creation.BrandID
	transaction.Price = creation.Price * dto.Amount
	transaction.TradeAt = time.Now()
	err = service.transaction.Create(ctx, transaction)
	txn = &TransactionDto{}
	err = copier.Copy(txn, &transaction)
	if err != nil {
		return
	}
	return
}

type TransactionDto struct {
	TransactionID string    `json:"transactionId"`
	Creation      string    `json:"creation"`
	BrandID       string    `json:"brandId"`
	Buyer         string    `json:"buyer"`
	Seller        string    `json:"seller"`
	Price         int       `json:"price"`
	Amount        int       `json:"amount"`
	TradeAt       time.Time `json:"tradeAt"`
}

func (dto *TransactionDto) ID(id primitive.ObjectID) {
	dto.TransactionID = id.Hex()
}

type TradeInCreationDto struct {
	CreationID string `json:"creationId"`
	Buyer      string `json:"buyer"`
	Seller     string `json:"seller"`
	Amount     int    `json:"amount"`
}

type TransactionFilterDto struct {
	CreationID   *string    `json:"creationID"`
	BrandID      *string    `json:"brandId"`
	Buyer        *string    `json:"buyer"`
	Seller       *string    `json:"seller"`
	CreationName *string    `json:"creationName"`
	Properties   []string   `json:"properties"`
	TradedBefore *time.Time `json:"tradedBefore"`
	TradedAfter  *time.Time `json:"tradedAfter"`
	MaxPrice     *int       `json:"maxPrice"`
	MinPrice     *int       `json:"minPrice"`
}

type TradeServiceError struct {
	ServiceError
}

func NewTradeServiceError(e ServiceEvent) error {
	return &TradeServiceError{ServiceError{ServiceName: "TradeService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}
