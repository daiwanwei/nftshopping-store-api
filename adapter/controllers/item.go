package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/business/services"
)

type ItemController interface {
	OrderItem(ctx *gin.Context)
	DeliverItem(ctx *gin.Context)
	FindItem(ctx *gin.Context)
	FindAllItem(ctx *gin.Context)
	GetAmountOfItem(ctx *gin.Context)
}

type itemController struct {
	item services.ItemService
}

func NewItemController() (controller ItemController, err error) {
	service, err := services.GetService()
	if err != nil {
		return
	}
	return &itemController{
		item: service.Item,
	}, nil
}

// OrderItem godoc
// @Summary 商品訂貨
// @Tags item
// @produce application/json
// @Param OrderItemDto body services.OrderItemDto true "商品訂購資料"
// @Success 200 {object}  adapter.NonDataResp "成功後返回的值"
// @Router /api/item/orderItem [post]
// @Security JWT
func (controller *itemController) OrderItem(ctx *gin.Context) {
	order := services.OrderItemDto{}
	if err := ctx.ShouldBindJSON(&order); err != nil {
		respond(ctx, err)
		return
	}
	err := controller.item.OrderItem(context.TODO(), order)
	respond(ctx, err)
}

// DeliverItem godoc
// @Summary 商品交貨
// @Tags item
// @produce application/json
// @Param DeliverItemDto body services.DeliverItemDto true "刊登房屋物件"
// @Success 200 {object}  adapter.DataResp{data=services.ItemDto} "成功後返回的值"
// @Router /api/item/deliverItem [post]
// @Security JWT
func (controller *itemController) DeliverItem(ctx *gin.Context) {
	delivery := services.DeliverItemDto{}
	if err := ctx.ShouldBindJSON(&delivery); err != nil {
		respondWithData(ctx, nil, err)
		return
	}
	item, err := controller.item.DeliverItem(context.TODO(), delivery)
	respondWithData(ctx, item, err)
}

// FindItem godoc
// @Summary 取得商品資訊
// @Tags item
// @produce application/json
// @Param contract path string true "contract"
// @Param token path string true "token"
// @Success 200 {object}  adapter.DataResp{data=services.ItemDto} "成功後返回的值"
// @Router /api/item/{contract}/{token} [get]
func (controller *itemController) FindItem(ctx *gin.Context) {
	contract := ctx.Param("contract")
	token := ctx.Param("token")
	items, err := controller.item.FindItem(context.TODO(), contract, token)
	respondWithData(ctx, items, err)
}

// FindAllItem godoc
// @Summary 取得所有商品資訊
// @Tags item
// @produce application/json
// @Param owner query string false "search by owner"
// @Param brandOwner query string false "search by brandOwner"
// @Param page query string false "search by page"
// @Param size query string false "search by size"
// @Param sort query string false "search by sort"
// @Param order query int false "search by order"
// @Success 200 {object}  adapter.DataResp{data=[]services.ItemDto} "成功後返回的值"
// @Router /api/item/findAllItem [get]
// @Security JWT
func (controller *itemController) FindAllItem(ctx *gin.Context) {
	pageable, err := getPageFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	filter, err := getItemFilterFromQuery(ctx)
	if err != nil {
		respondWithData(ctx, nil, err)
		return
	}

	var items []services.ItemDto

	if pageable.Page < 0 {
		items, err = controller.item.FindAllItemByFilter(context.TODO(), filter)
	} else {
		items, err = controller.item.FindAllItemByFilterAndPage(context.TODO(), filter, *pageable)
	}
	respondWithData(ctx, items, err)
}

// GetAmountOfItem godoc
// @Summary 取得商品總量
// @Tags item
// @produce application/json
// @Param brandId query string true "brandId"
// @Success 200 {object}  adapter.DataResp{data=int64} "成功後返回的值"
// @Router /api/item/getAmountOfItem [get]
func (controller *itemController) GetAmountOfItem(ctx *gin.Context) {
	brandId := ctx.Query("brandId")
	amount, err := controller.item.GetAmountOfItemByBrand(context.TODO(), brandId)
	respondWithData(ctx, amount, err)
}

func getItemFilterFromQuery(ctx *gin.Context) (filter services.ItemFilterDto, err error) {
	if owner := ctx.Query("owner"); len(owner) > 0 {
		filter.Owner = &owner
	}
	if brandOwner := ctx.Query("brandOwner"); len(brandOwner) > 0 {
		filter.BrandOwner = &brandOwner
	}
	return
}
