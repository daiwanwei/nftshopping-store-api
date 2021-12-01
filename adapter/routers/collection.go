package routers

import (
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/adapter/controllers"
)

func InitCollectionRouter(engine *gin.Engine) (err error) {
	controller, err := controllers.GetController()
	if err != nil {
		return
	}
	app := engine.Group("api")

	collect := app.Group("collection")
	collect.GET("/findCollection", controller.Collection.FindCollection)
	collect.GET("/findAllCollection", controller.Collection.FindAllCollection)
	return
}
