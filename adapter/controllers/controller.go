package controllers

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"nftshopping-store-api/adapter"
	"nftshopping-store-api/pkg/log"
	"nftshopping-store-api/pkg/utils"
	"strconv"
)

var controllerInstance *controller

func GetController() (instance *controller, err error) {
	if controllerInstance == nil {
		instance, err = newController()
		if err != nil {
			return nil, err
		}
		controllerInstance = instance
	}
	return controllerInstance, nil
}

type controller struct {
	User       UserController
	Creation   CreationController
	Item       ItemController
	Collection CollectionController
	Trade      TradeController
	Brand      BrandController
	Stock      StockController
}

func newController() (instance *controller, err error) {
	user, err := NewUserController()
	if err != nil {
		return
	}
	creation, err := NewCreationController()
	if err != nil {
		return
	}
	item, err := NewItemController()
	if err != nil {
		return
	}
	collection, err := NewCollectionController()
	if err != nil {
		return
	}
	trade, err := NewTradeController()
	if err != nil {
		return
	}
	brand, err := NewBrandController()
	if err != nil {
		return
	}
	stock, err := NewStockController()
	if err != nil {
		return
	}
	return &controller{
		User:       user,
		Creation:   creation,
		Item:       item,
		Collection: collection,
		Trade:      trade,
		Brand:      brand,
		Stock:      stock,
	}, nil
}

func getPageFromQuery(ctx *gin.Context) (pageable *utils.Pageable, err error) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	if err != nil {
		return
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		return
	}
	var sort map[string]int
	if sortType := ctx.Query("sort"); len(sortType) != 0 {
		order, err := strconv.Atoi(ctx.DefaultQuery("order", "1"))
		if err != nil {
			return nil, err
		}
		sort = map[string]int{
			sortType: order,
		}
	}
	return &utils.Pageable{
		Size: size, Page: page, Sort: sort,
	}, nil
}

func respondWithData(ctx *gin.Context, data interface{}, err error) {
	if err != nil {
		sentry.CaptureException(err)
		logger, logErr := log.GetLog()
		if logErr != nil {
			ctx.JSON(http.StatusOK, adapter.DataResp{Code: 500, Msg: logErr.Error(), Data: nil})
			return
		}
		logger.Error(err)
		if e, ok := err.(utils.CustomError); ok {
			ctx.JSON(http.StatusOK, adapter.DataResp{Code: e.GetCode(), Msg: e.GetMsg(), Data: nil})
			return
		}
		ctx.JSON(http.StatusOK, adapter.DataResp{Code: 500, Msg: err.Error(), Data: nil})
		return
	} else {
		ctx.JSON(http.StatusOK, adapter.DataResp{Code: 200, Msg: "OK", Data: data})
	}
}

func respond(ctx *gin.Context, err error) {
	if err != nil {
		sentry.CaptureException(err)
		logger, logErr := log.GetLog()
		if logErr != nil {
			ctx.JSON(http.StatusOK, adapter.DataResp{Code: 500, Msg: logErr.Error(), Data: nil})
			return
		}
		logger.Error(err)
		if e, ok := err.(utils.CustomError); ok {
			ctx.JSON(http.StatusOK, adapter.NonDataResp{Code: e.GetCode(), Msg: e.GetMsg()})
			return
		}
		ctx.JSON(http.StatusOK, adapter.NonDataResp{Code: 500, Msg: err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, adapter.NonDataResp{Code: 200, Msg: "OK"})
	}
}
