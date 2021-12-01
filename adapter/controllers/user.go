package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/business/services"
	"nftshopping-store-api/pkg/transactions"
)

type UserController interface {
	Exist(ctx *gin.Context)
	Register(ctx *gin.Context)
	FindUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type userController struct {
	user services.UserService
}

func NewUserController() (controller UserController, err error) {
	service, err := services.GetService()
	if err != nil {
		return
	}
	return &userController{service.User}, nil
}

// Exist godoc
// @Summary  查看會員
// @Tags user
// @produce application/json
// @Param userId query string false "search by userId"
// @Param userName query string false "search by userName"
// @Success 200 {object}  adapter.DataResp{data=bool} "成功後返回的值"
// @Router /api/user/exist [get]
// @Security JWT
func (controller *userController) Exist(ctx *gin.Context) {
	userId := ctx.Query("userId")
	var isExisted bool
	var err error
	if len(userId) == 0 {
		isExisted, err = controller.user.ExistByAccount(context.TODO(), ctx.Query("userName"))
	} else {
		isExisted, err = controller.user.ExistByID(context.TODO(), ctx.Query("userId"))
	}
	respondWithData(ctx, isExisted, err)
}

// Register godoc
// @Summary 建立會員
// @Tags user
// @produce application/json
// @Param RegisterUser body services.RegisterUserDto true "註冊Dto"
// @Success 200 {object}  adapter.DataResp{data=services.UserDto} "成功後返回的值"
// @Router /api/user/register [post]
func (controller *userController) Register(ctx *gin.Context) {
	register := services.RegisterUserDto{}
	if err := ctx.ShouldBindJSON(&register); err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	txn, err := transactions.NewTransaction("Register")
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}
	defer txn.End(context.Background())

	callback := func(transactionCtx transactions.TransactionContext) (interface{}, error) {
		user, err := controller.user.Register(transactionCtx, register)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	user, err := txn.With(context.Background(), callback)
	respondWithData(ctx, user, err)
}

// DeleteUser godoc
// @Summary 刪除會員
// @Tags user
// @produce application/json
// @Param userId query string false "search by userId"
// @Success 200 {object}  adapter.NonDataResp "成功後返回的值"
// @Router /api/user/deleteUser [get]
// @Security JWT
func (controller *userController) DeleteUser(ctx *gin.Context) {
	err := controller.user.DeleteUser(context.TODO(), ctx.Query("userId"))
	respond(ctx, err)
}

// FindUser godoc
// @Summary  取的會員資料
// @Tags user
// @produce application/json
// @Param userId query string false "search by userId"
// @Param userName query string false "search by userName"
// @Success 200 {object}  adapter.DataResp{data=services.UserDto} "成功後返回的值"
// @Router /api/user/findUser [get]
// @Security JWT
func (controller *userController) FindUser(ctx *gin.Context) {
	userId := ctx.Query("userId")
	var user *services.UserDto
	var err error
	if len(userId) == 0 {
		user, err = controller.user.FindUserByAccount(context.TODO(), ctx.Query("userName"))
	} else {
		user, err = controller.user.FindUserByID(context.TODO(), ctx.Query("userId"))
	}
	respondWithData(ctx, user, err)
}
