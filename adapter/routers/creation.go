package routers

import (
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/adapter/controllers"
)

func InitCreationRouter(engine *gin.Engine) (err error) {
	controller, err := controllers.GetController()
	if err != nil {
		return
	}
	app := engine.Group("api")

	creation := app.Group("creation")
	creation.GET("/findCreation", controller.Creation.FindCreation)
	creation.GET("/findAllCreation", controller.Creation.FindAllCreation)
	creation.POST("/postCreation", controller.Creation.PostCreation)
	creation.DELETE("/deleteCreation", controller.Creation.DeleteCreation)
	creation.POST("/updateCreation", controller.Creation.UpdateCreation)
	return
}
