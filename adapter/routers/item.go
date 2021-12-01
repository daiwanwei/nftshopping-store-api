package routers

import (
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/adapter/controllers"
)

func InitItemRouter(engine *gin.Engine) (err error) {
	controller, err := controllers.GetController()
	if err != nil {
		return
	}
	app := engine.Group("api")

	item := app.Group("item")
	item.GET("/:contract/:token", controller.Item.FindItem)
	item.POST("/orderItem", controller.Item.OrderItem)
	item.POST("/deliverItem", controller.Item.DeliverItem)
	item.GET("/findAllItem", controller.Item.FindAllItem)
	item.GET("/getAmountOfItem", controller.Item.GetAmountOfItem)
	return
}
