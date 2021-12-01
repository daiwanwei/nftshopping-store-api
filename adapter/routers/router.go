package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"nftshopping-store-api/adapter"
	"nftshopping-store-api/adapter/middlewares"
)

var routerInstance *gin.Engine

func GetRouter() (instance *gin.Engine, err error) {
	if routerInstance == nil {
		instance, err = newRouter()
		if err != nil {
			return nil, err
		}
		routerInstance = instance
	}
	return routerInstance, nil
}

func newRouter() (router *gin.Engine, err error) {
	middleware, err := middlewares.GetMiddleware()
	if err != nil {
		return
	}
	engine := gin.Default()
	engine.Use(middleware.Cors.Cors())

	//k8s探針
	engine.GET("/probe", func(context *gin.Context) {
		context.JSON(http.StatusOK, adapter.NonDataResp{Code: 200, Msg: "OK"})
	})

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = InitUserRouter(engine)
	if err != nil {
		return
	}
	err = InitCreationRouter(engine)
	if err != nil {
		return
	}
	err = InitItemRouter(engine)
	if err != nil {
		return
	}
	err = InitCollectionRouter(engine)
	if err != nil {
		return
	}
	err = InitTradeRouter(engine)
	if err != nil {
		return
	}
	err = InitBrandRouter(engine)
	if err != nil {
		return
	}
	err = InitStockRouter(engine)
	if err != nil {
		return
	}
	return engine, nil
}
