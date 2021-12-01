package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/business/services"
	"strconv"
	"strings"
	"time"
)

type TradeController interface {
	FindTransaction(ctx *gin.Context)
	FindAllTransaction(ctx *gin.Context)
	TradeInCreation(ctx *gin.Context)
}

type tradeController struct {
	trade services.TradeService
}

func NewTradeController() (controller TradeController, err error) {
	service, err := services.GetService()
	if err != nil {
		return
	}
	return &tradeController{
		trade: service.Trade,
	}, nil
}

// FindTransaction godoc
// @Summary 取得交易資訊
// @Tags trade
// @produce application/json
// @Param transactionId query string false "search by transactionId"
// @Success 200 {object}  adapter.DataResp{data=services.TransactionDto} "成功後返回的值"
// @Router /api/trade/findTransaction [get]
// @Security JWT
func (controller *tradeController) FindTransaction(ctx *gin.Context) {
	transactionId := ctx.Query("transactionId")
	transaction, err := controller.trade.FindTransaction(context.TODO(), transactionId)
	respondWithData(ctx, transaction, err)
}

// FindAllTransaction godoc
// @Summary 取得所有交易資訊
// @Tags trade
// @produce application/json
// @Param page query string false "search by page"
// @Param size query string false "search by size"
// @Param properties query string false "search by properties"
// @Param tradedBefore query string false "search by tradedBefore"
// @Param tradedAfter query string false "search by tradedAfter"
// @Param maxPrice query int false "search by maxPrice"
// @Param minPrice query int false "search by minPrice"
// @Param properties query string false "search by properties"
// @Param creationName query string false "search by creationName"
// @Param creationId query string false "search by creationId"
// @Param brandId query string false "search by brandId"
// @Param buyer query string false "search by buyer"
// @Param seller query string false "search by seller"
// @Param sort query string false "search by sort"
// @Param order query int false "search by order"
// @Success 200 {object}  adapter.DataResp{data=[]services.TransactionDto} "成功後返回的值"
// @Router /api/trade/findAllTransaction [get]
// @Security JWT
func (controller *tradeController) FindAllTransaction(ctx *gin.Context) {
	pageable, err := getPageFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	filter, err := getTransactionFilterFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	var transactions []services.TransactionDto

	if pageable.Page < 0 {
		transactions, err = controller.trade.FindAllTransactionByFilter(context.TODO(), filter)
	} else {
		transactions, err = controller.trade.FindAllTransactionByFilterAndPage(context.TODO(), filter, *pageable)
	}
	respondWithData(ctx, transactions, err)
}

// TradeInCreation godoc
// @Summary 交易藝術品
// @Tags trade
// @produce application/json
// @Param tradeInCreation body services.TradeInCreationDto true "藝術品交易資料"
// @Success 200 {object}  adapter.DataResp{data=services.TransactionDto} "成功後返回的值"
// @Router /api/trade/tradeInCreation [post]
// @Security JWT
func (controller *tradeController) TradeInCreation(ctx *gin.Context) {
	trade := services.TradeInCreationDto{}
	if err := ctx.ShouldBindJSON(&trade); err != nil {
		respondWithData(ctx, nil, err)
		return
	}
	txn, err := controller.trade.TradeInCreation(context.TODO(), trade)
	respondWithData(ctx, txn, err)
}

func getTransactionFilterFromQuery(ctx *gin.Context) (filter services.TransactionFilterDto, err error) {
	if properties := ctx.Query("properties"); len(properties) > 0 {
		filter.Properties = strings.Split(properties, ",")
	}

	if creationId := ctx.Query("creationId"); len(creationId) > 0 {
		filter.CreationID = &creationId
	}

	if brandId := ctx.Query("brandId"); len(brandId) > 0 {
		filter.BrandID = &brandId
	}

	if buyer := ctx.Query("buyer"); len(buyer) > 0 {
		filter.Buyer = &buyer
	}

	if seller := ctx.Query("seller"); len(seller) > 0 {
		filter.Seller = &seller
	}

	if creationName := ctx.Query("creationName"); len(creationName) > 0 {
		filter.CreationName = &creationName
	}

	if after := ctx.Query("tradedAfter"); len(after) > 0 {
		tradedAfter, err := time.Parse(time.RFC3339, after)
		if err != nil {
			return filter, err
		}
		filter.TradedAfter = &tradedAfter
	}

	if before := ctx.Query("tradedBefore"); len(before) > 0 {
		tradedBefore, err := time.Parse(time.RFC3339, before)
		if err != nil {
			return filter, err
		}
		filter.TradedBefore = &tradedBefore
	}

	if min := ctx.Query("minPrice"); len(min) > 0 {
		minPrice, err := strconv.Atoi(min)
		if err != nil {
			return filter, err
		}
		filter.MinPrice = &minPrice
	}

	if max := ctx.Query("maxPrice"); len(max) > 0 {
		maxPrice, err := strconv.Atoi(max)
		if err != nil {
			return filter, err
		}
		filter.MaxPrice = &maxPrice
	}
	return
}
