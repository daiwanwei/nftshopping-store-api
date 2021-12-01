package routers

import (
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/adapter/controllers"
)

func InitBrandRouter(engine *gin.Engine) (err error) {
	controller, err := controllers.GetController()
	if err != nil {
		return
	}
	app := engine.Group("api")

	brand := app.Group("brand")
	brand.GET("/findBrand", controller.Brand.FindBrand)
	brand.GET("/findAllBrand", controller.Brand.FindAllBrand)
	brand.POST("/postBrand", controller.Brand.PostBrand)
	brand.POST("/updateBrand", controller.Brand.UpdateBrand)
	brand.DELETE("/deleteBrand", controller.Brand.DeleteBrand)
	return
}
