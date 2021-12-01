package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/business/services"
	"time"
)

type BrandController interface {
	FindBrand(ctx *gin.Context)
	FindAllBrand(ctx *gin.Context)
	PostBrand(ctx *gin.Context)
	UpdateBrand(ctx *gin.Context)
	DeleteBrand(ctx *gin.Context)
}

type brandController struct {
	brand services.BrandService
}

func NewBrandController() (controller BrandController, err error) {
	service, err := services.GetService()
	if err != nil {
		return
	}
	return &brandController{
		brand: service.Brand,
	}, nil
}

// FindBrand godoc
// @Summary 取得品牌
// @Tags brand
// @produce application/json
// @Param brandId query string false "search by brandId"
// @Success 200 {object}  adapter.DataResp{data=services.BrandDto} "成功後返回的值"
// @Router /api/brand/findBrand [get]
// @Security JWT
func (controller *brandController) FindBrand(ctx *gin.Context) {
	id := ctx.Query("brandId")
	brand, err := controller.brand.FindBrandById(context.TODO(), id)
	respondWithData(ctx, brand, err)
}

// FindAllBrand godoc
// @Summary 取得所有品牌
// @Tags brand
// @produce application/json
// @Param createAfter query string false "search by createAfter"
// @Param createBefore query string false "search by createBefore"
// @Param name query string false "search by name"
// @Param page query string false "search by page"
// @Param size query string false "search by size"
// @Param sort query string false "search by sort"
// @Param order query int false "search by order"
// @Success 200 {object}  adapter.DataResp{data=[]services.BrandDto} "成功後返回的值"
// @Router /api/brand/findAllBrand [get]
// @Security JWT
func (controller *brandController) FindAllBrand(ctx *gin.Context) {

	pageable, err := getPageFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	filter, err := getBrandFilterFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	var brands []services.BrandDto
	if pageable.Page < 0 {
		brands, err = controller.brand.FindAllBrandByFilter(context.TODO(), filter)
	} else {
		brands, err = controller.brand.FindAllBrandByFilterAndPage(context.TODO(), filter, *pageable)
	}
	respondWithData(ctx, brands, err)
}

// PostBrand godoc
// @Summary 刊登品牌
// @Tags brand
// @produce application/json
// @Param PostBrand body services.PostBrandDto true "建立品牌"
// @Success 200 {object}  adapter.DataResp{data=services.BrandDto} "成功後返回的值"
// @Router /api/brand/postBrand [post]
// @Security JWT
func (controller *brandController) PostBrand(ctx *gin.Context) {
	post := services.PostBrandDto{}
	if err := ctx.ShouldBindJSON(&post); err != nil {
		respondWithData(ctx, nil, err)
		return
	}
	brand, err := controller.brand.PostBrand(context.TODO(), post)
	respondWithData(ctx, brand, err)
}

// UpdateBrand godoc
// @Summary 更新品牌
// @Tags brand
// @produce application/json
// @Param PostBrand body services.UpdateBrandDto true "編輯品牌"
// @Success 200 {object}  adapter.NonDataResp "成功後返回的值"
// @Router /api/brand/updateBrand [post]
// @Security JWT
func (controller *brandController) UpdateBrand(ctx *gin.Context) {
	post := services.UpdateBrandDto{}
	if err := ctx.ShouldBindJSON(&post); err != nil {
		respond(ctx, err)
		return
	}
	err := controller.brand.UpdateBrand(context.TODO(), post)
	respond(ctx, err)
}

// DeleteBrand godoc
// @Summary DeleteBrand
// @Tags brand
// @produce application/json
// @Param brandId query string false "search by brandId"
// @Success 200 {object}  adapter.NonDataResp "成功後返回的值"
// @Router /api/brand/deleteBrand [delete]
// @Security JWT
func (controller *brandController) DeleteBrand(ctx *gin.Context) {
	id := ctx.Query("brandId")
	err := controller.brand.DeleteBrand(context.TODO(), id)
	respond(ctx, err)
}

func getBrandFilterFromQuery(ctx *gin.Context) (filter services.BrandFilterDto, err error) {
	if name := ctx.Query("name"); len(name) > 0 {
		filter.Name = &name
	}

	if after := ctx.Query("createAfter"); len(after) > 0 {
		createAfter, err := time.Parse(time.RFC3339, after)
		if err != nil {
			return filter, err
		}
		filter.CreateAfter = &createAfter
	}

	if before := ctx.Query("createBefore"); len(before) > 0 {
		createBefore, err := time.Parse(time.RFC3339, before)
		if err != nil {
			return filter, err
		}
		filter.CreateBefore = &createBefore
	}
	return
}
