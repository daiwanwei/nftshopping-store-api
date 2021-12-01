package services

import (
	"context"
	"github.com/jinzhu/copier"
	"nftshopping-store-api/persistence/repositories"
	"nftshopping-store-api/pkg/utils"
	"time"
)

type BrandService interface {
	Exist(ctx context.Context, brandId string) (isExisted bool, err error)
	FindBrandById(ctx context.Context, id string) (brandDto *BrandDto, err error)
	FindAllBrandByFilter(ctx context.Context, dto BrandFilterDto) (brandsDto []BrandDto, err error)
	FindAllBrandByFilterAndPage(ctx context.Context, dto BrandFilterDto, pageable utils.Pageable) (brandsDto []BrandDto, err error)
	PostBrand(ctx context.Context, dto PostBrandDto) (brandsDto *BrandDto, err error)
	UpdateBrand(ctx context.Context, dto UpdateBrandDto) (err error)
	DeleteBrand(ctx context.Context, brandId string) (err error)
}

type brandService struct {
	brand repositories.BrandDao
}

func NewBrandService() (service BrandService, err error) {
	dao, err := repositories.GetRepository()
	if err != nil {
		return nil, err
	}
	return &brandService{
		brand: dao.Brand,
	}, nil
}

func (service *brandService) Exist(ctx context.Context, id string) (isExisted bool, err error) {
	isExisted, err = service.brand.Exist(ctx, id)
	return
}

func (service *brandService) FindBrandById(ctx context.Context, id string) (brandDto *BrandDto, err error) {
	brand, err := service.brand.Find(ctx, id)
	if err != nil || brand == nil {
		return
	}
	brandDto = &BrandDto{}
	if err = copier.Copy(brandDto, brand); err != nil {
		return nil, err
	}
	return
}

func (service *brandService) FindAllBrandByFilter(ctx context.Context, dto BrandFilterDto) (brandsDto []BrandDto, err error) {
	selector := repositories.BrandSelector{
		Name:         dto.Name,
		CreateAfter:  dto.CreateAfter,
		CreateBefore: dto.CreateBefore,
	}
	brands, err := service.brand.FindAllByFilter(ctx, repositories.SelectorOfBrand(selector))
	if err != nil {
		return
	}
	if err = copier.Copy(&brandsDto, &brands); err != nil {
		return nil, err
	}
	return
}

func (service *brandService) FindAllBrandByFilterAndPage(
	ctx context.Context, dto BrandFilterDto, pageable utils.Pageable,
) (brandsDto []BrandDto, err error) {
	selector := repositories.BrandSelector{
		Name:         dto.Name,
		CreateAfter:  dto.CreateAfter,
		CreateBefore: dto.CreateBefore,
	}
	page, err := service.brand.FindAllByFilterAndPage(ctx, repositories.SelectorOfBrand(selector), pageable)
	if err != nil || page == nil {
		return
	}
	if brands, ok := page.Content.([]repositories.Brand); !ok {
		return nil, utils.ErrCovertContent
	} else {
		err = copier.Copy(&brandsDto, &brands)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (service *brandService) PostBrand(ctx context.Context, dto PostBrandDto) (brandDto *BrandDto, err error) {
	brand := &repositories.Brand{}
	if err = copier.Copy(brand, &dto); err != nil {
		return
	}
	brand.ID = dto.Name
	brand.CreateAt = time.Now()
	err = service.brand.Create(ctx, brand)
	if err != nil {
		return
	}
	brandDto = &BrandDto{}
	err = copier.Copy(brandDto, brand)
	if err != nil {
		return
	}
	return
}

func (service *brandService) UpdateBrand(ctx context.Context, dto UpdateBrandDto) (err error) {
	brand, err := service.brand.Find(ctx, dto.BrandID)
	if err != nil {
		return
	}

	if brand == nil {
		return NewBrandServiceError(BrandNotFound)
	}
	err = copier.Copy(brand, &dto)
	if err != nil {
		return
	}
	err = service.brand.Save(ctx, brand)
	if err != nil {
		return
	}
	return
}

func (service *brandService) DeleteBrand(ctx context.Context, brandId string) (err error) {
	err = service.brand.Delete(ctx, brandId)
	if err != nil {
		return
	}
	return
}

type BrandDto struct {
	BrandID     string    `json:"brandId"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"imageUrl"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"createAt"`
}

func (dto *BrandDto) ID(id string) {
	dto.BrandID = id
}

type BrandFilterDto struct {
	Name         *string    `json:"name"`
	CreateAfter  *time.Time `json:"createAfter"`
	CreateBefore *time.Time `json:"createBefore"`
}

type PostBrandDto struct {
	Name        string `json:"name"`
	ImageURL    string `json:"imageUrl"`
	Description string `json:"description"`
}

type UpdateBrandDto struct {
	BrandID  string `json:"brandId"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

type BrandServiceError struct {
	ServiceError
}

func NewBrandServiceError(e ServiceEvent) error {
	return &BrandServiceError{ServiceError{ServiceName: "BrandService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}
