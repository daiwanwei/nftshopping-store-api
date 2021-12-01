package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nftshopping-store-api/business/services"
	"strconv"
	"strings"
	"time"
)

type CreationController interface {
	FindCreation(ctx *gin.Context)
	FindAllCreation(ctx *gin.Context)
	PostCreation(ctx *gin.Context)
	DeleteCreation(ctx *gin.Context)
	UpdateCreation(ctx *gin.Context)
}

type creationController struct {
	creation services.CreationService
}

func NewCreationController() (controller CreationController, err error) {
	service, err := services.GetService()
	if err != nil {
		return
	}
	return &creationController{
		creation: service.Creation,
	}, nil
}

// FindCreation godoc
// @Summary 取得藝術品資訊
// @Tags creation
// @produce application/json
// @Param creationId query string false "search by creationId"
// @Success 200 {object}  adapter.DataResp{data=services.CreationDto} "成功後返回的值"
// @Router /api/creation/findCreation [get]
// @Security JWT
func (controller *creationController) FindCreation(ctx *gin.Context) {
	creationId := ctx.Query("creationId")
	creation, err := controller.creation.FindCreationByID(context.TODO(), creationId)
	respondWithData(ctx, creation, err)
}

// FindAllCreation godoc
// @Summary 取得所有藝術品資訊
// @Tags creation
// @produce application/json
// @Param page query string false "search by page"
// @Param size query string false "search by size"
// @Param creationIds query string false "search by creationIds"
// @Param properties query string false "search by properties"
// @Param saleStartBefore query string false "search by saleStartBefore"
// @Param saleStartAfter query string false "search by saleStartAfter"
// @Param maxPrice query int false "search by maxPrice"
// @Param minPrice query int false "search by minPrice"
// @Param creationName query string false "search by creationName"
// @Param creator query string false "search by creator"
// @Param brandId query string false "search by brandId"
// @Param sort query string false "search by sort"
// @Param order query int false "search by order"
// @Success 200 {object}  adapter.DataResp{data=[]services.CreationDto} "成功後返回的值"
// @Router /api/creation/findAllCreation [get]
// @Security JWT
func (controller *creationController) FindAllCreation(ctx *gin.Context) {
	pageable, err := getPageFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	filter, err := getCreationFilterFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	var creations []services.CreationDto

	if pageable.Page < 0 {
		creations, err = controller.creation.FindAllCreationByFilter(context.TODO(), filter)
	} else {
		creations, err = controller.creation.FindAllCreationByFilterAndPage(context.TODO(), filter, *pageable)
	}
	respondWithData(ctx, creations, err)
}

// PostCreation godoc
// @Summary 刊登藝術品
// @Tags creation
// @produce application/json
// @Param postCreation body services.PostCreationDto true "藝術品資料"
// @Success 200 {object}  adapter.DataResp{data=services.CreationDto} "成功後返回的值"
// @Router /api/creation/postCreation [post]
// @Security JWT
func (controller *creationController) PostCreation(ctx *gin.Context) {
	postCreation := services.PostCreationDto{}
	if err := ctx.ShouldBindJSON(&postCreation); err != nil {
		respondWithData(ctx, nil, err)
		return
	}
	creation, err := controller.creation.PostCreation(context.TODO(), postCreation)
	respondWithData(ctx, creation, err)
}

// DeleteCreation godoc
// @Summary 刪除藝術品
// @Tags creation
// @produce application/json
// @Param id query string false "search by id"
// @Success 200 {object}  adapter.NonDataResp "成功後返回的值"
// @Router /api/creation/deleteCreation [delete]
// @Security JWT
func (controller *creationController) DeleteCreation(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Query("id"))
	if err != nil {
		respond(ctx, err)
	}
	err = controller.creation.DeleteCreation(context.TODO(), id)
	respond(ctx, err)
}

// UpdateCreation godoc
// @Summary 更新藝術品
// @Tags creation
// @produce application/json
// @Param postCreation body services.UpdateCreationDto true "編輯藝術品"
// @Success 200 {object}  adapter.NonDataResp "成功後返回的值"
// @Router /api/creation/updateCreation [post]
// @Security JWT
func (controller *creationController) UpdateCreation(ctx *gin.Context) {
	creation := services.UpdateCreationDto{}
	if err := ctx.ShouldBindJSON(&creation); err != nil {
		respond(ctx, err)
		return
	}
	err := controller.creation.UpdateCreation(context.TODO(), creation)
	respond(ctx, err)
}

func getCreationFilterFromQuery(ctx *gin.Context) (filter services.CreationFilterDto, err error) {
	if creationIds := ctx.Query("creationIds"); len(creationIds) > 0 {
		filter.CreationIDs = strings.Split(creationIds, ",")
	}
	if properties := ctx.Query("properties"); len(properties) > 0 {
		filter.Properties = strings.Split(properties, ",")
	}

	if creationName := ctx.Query("creationName"); len(creationName) > 0 {
		filter.CreationName = &creationName
	}

	if creator := ctx.Query("creator"); len(creator) > 0 {
		filter.Creator = &creator
	}

	if brand := ctx.Query("brandId"); len(brand) > 0 {
		filter.BrandID = &brand
	}

	if after := ctx.Query("saleStartAfter"); len(after) > 0 {
		saleStartAfter, err := time.Parse(time.RFC3339, after)
		if err != nil {
			return filter, err
		}
		filter.SaleEndAfter = &saleStartAfter
	}

	if before := ctx.Query("saleStartBefore"); len(before) > 0 {
		saleStartBefore, err := time.Parse(time.RFC3339, before)
		if err != nil {
			return filter, err
		}
		filter.SaleStartBefore = &saleStartBefore
	}

	if min := ctx.Query("minPrice"); len(min) > 0 {
		minPrice, err := strconv.Atoi(min)
		if err != nil {
			return filter, err
		}
		filter.MinPrice = &minPrice
	}

	if max := ctx.Query("maxPrice"); len(max) > 0 {
		maxPrice, err := strconv.Atoi(max)
		if err != nil {
			return filter, err
		}
		filter.MaxPrice = &maxPrice
	}
	return
}
