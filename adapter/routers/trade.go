package routers

import (
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/adapter/controllers"
)

func InitTradeRouter(engine *gin.Engine) (err error) {
	controller, err := controllers.GetController()
	if err != nil {
		return
	}
	app := engine.Group("api")

	user := app.Group("trade")
	user.GET("/findTransaction", controller.Trade.FindTransaction)
	user.GET("/findAllTransaction", controller.Trade.FindAllTransaction)
	user.POST("/tradeInCreation", controller.Trade.TradeInCreation)
	return
}
