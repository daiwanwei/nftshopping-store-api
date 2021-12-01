package routers

import (
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/adapter/controllers"
)

func InitUserRouter(engine *gin.Engine) (err error) {
	controller, err := controllers.GetController()
	if err != nil {
		return
	}
	app := engine.Group("api")

	user := app.Group("user")
	user.GET("/exist", controller.User.Exist)
	user.POST("/register", controller.User.Register)
	user.GET("/findUser", controller.User.FindUser)
	user.GET("/deleteUser", controller.User.DeleteUser)
	return
}
