package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/business/services"
)

type CollectionController interface {
	FindCollection(ctx *gin.Context)
	FindAllCollection(ctx *gin.Context)
}

type collectController struct {
	collection services.CollectService
}

func NewCollectionController() (controller CollectionController, err error) {
	service, err := services.GetService()
	if err != nil {
		return
	}
	return &collectController{
		collection: service.Collection,
	}, nil
}

// FindCollection godoc
// @Summary 取得蒐藏
// @Tags collection
// @produce application/json
// @Param owner query string true "owner"
// @Param creationId query string true "creationId"
// @Success 200 {object}  adapter.DataResp{data=services.CollectDto} "成功後返回的值"
// @Router /api/collection/findCollection [get]
func (controller *collectController) FindCollection(ctx *gin.Context) {
	owner := ctx.Query("owner")
	creationId := ctx.Query("creationId")
	collection, err := controller.collection.FindCollect(context.TODO(), owner, creationId)
	respondWithData(ctx, collection, err)
}

// FindAllCollection godoc
// @Summary 取得所有蒐藏
// @Tags collection
// @produce application/json
// @Param owner query string false "search by owner"
// @Param creationId query string false "search by creationId"
// @Param page query string false "search by page"
// @Param size query string false "search by size"
// @Param sort query string false "search by sort"
// @Param order query int false "search by order"
// @Success 200 {object}  adapter.DataResp{data=services.CollectDto} "成功後返回的值"
// @Router /api/collection/findAllCollection [get]
func (controller *collectController) FindAllCollection(ctx *gin.Context) {
	pageable, err := getPageFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	filter, err := getCollectionFilterFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	var collections []services.CollectDto
	if pageable.Page < 0 {
		collections, err = controller.collection.FindAllCollectByFilter(context.TODO(), filter)
		if err != nil {
			respondWithData(ctx, nil, err)
			return
		}
	} else {
		collections, err = controller.collection.FindAllCollectByFilterAndPage(context.TODO(), filter, *pageable)
		if err != nil {
			respondWithData(ctx, nil, err)
			return
		}
	}
	respondWithData(ctx, collections, err)
}

func getCollectionFilterFromQuery(ctx *gin.Context) (filter services.CollectFilterDto, err error) {
	if owner := ctx.Query("owner"); len(owner) > 0 {
		filter.Owner = &owner
	}
	if creationId := ctx.Query("creationId"); len(creationId) > 0 {
		filter.CreationID = &creationId
	}
	return
}
