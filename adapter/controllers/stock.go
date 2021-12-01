package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/business/services"
)

type StockController interface {
	FindAllStock(ctx *gin.Context)
}

type stockController struct {
	stock services.StockService
}

func NewStockController() (controller StockController, err error) {
	service, err := services.GetService()
	if err != nil {
		return
	}
	return &stockController{
		stock: service.Stock,
	}, nil
}

// FindAllStock godoc
// @Summary 取得所有庫存
// @Tags stock
// @produce application/json
// @Param page query string false "search by page"
// @Param size query string false "search by size"
// @Param creationId query string false "search by creationId"
// @Param brandId query string false "search by brandId"
// @Success 200 {object}  adapter.DataResp{data=[]services.CollectDto} "成功後返回的值"
// @Router /api/stock/findAllStock [get]
// @Security JWT
func (controller *stockController) FindAllStock(ctx *gin.Context) {
	pageable, err := getPageFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	filter, err := getStockFilterFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}
	var stocks []services.CollectDto

	if pageable.Page < 0 {
		stocks, err = controller.stock.FindAllCollectByFilter(context.TODO(), filter)
		if err != nil {
			respondWithData(ctx, nil, err)
			return
		}
	} else {
		stocks, err = controller.stock.FindAllCollectByFilterAndPage(context.TODO(), filter, *pageable)
		if err != nil {
			respondWithData(ctx, nil, err)
			return
		}
	}
	respondWithData(ctx, stocks, err)
}

func getStockFilterFromQuery(ctx *gin.Context) (filter services.CollectFilterDto, err error) {
	if brandId := ctx.Query("brandId"); len(brandId) > 0 {
		filter.Owner = &brandId
	}

	if creationId := ctx.Query("creationId"); len(creationId) > 0 {
		filter.CreationID = &creationId
	}
	return
}
