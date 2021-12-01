package services

import (
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nftshopping-store-api/persistence/repositories"
)

type UserService interface {
	ExistByID(ctx context.Context, userId string) (isExisted bool, err error)
	ExistByAccount(ctx context.Context, account string) (isExisted bool, err error)
	Register(ctx context.Context, dto RegisterUserDto) (userDto *UserDto, err error)
	DeleteUser(ctx context.Context, userId string) (err error)
	FindUserByID(ctx context.Context, userId string) (userDto *UserDto, err error)
	FindUserByAccount(ctx context.Context, userName string) (userDto *UserDto, err error)
}

type userService struct {
	user repositories.UserDao
}

func NewUserService() (service UserService, err error) {
	dao, err := repositories.GetRepository()
	if err != nil {
		return nil, err
	}
	return &userService{
		user: dao.User,
	}, nil
}

func (service *userService) ExistByID(ctx context.Context, userId string) (isExisted bool, err error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}
	isExisted, err = service.user.ExistByID(ctx, id)
	return
}

func (service *userService) ExistByAccount(ctx context.Context, account string) (isExisted bool, err error) {
	isExisted, err = service.user.ExistByAccount(ctx, account)
	return
}

func (service *userService) Register(ctx context.Context, dto RegisterUserDto) (userDto *UserDto, err error) {
	user, err := service.user.FindByAccount(ctx, dto.EtherAccount)
	if err != nil {
		return
	}
	if user != nil {
		return nil, NewUserServiceError(UserRegistered)
	}
	user = &repositories.User{
		Account: dto.EtherAccount,
	}
	user.ID = primitive.NewObjectID()
	err = service.user.Create(ctx, user)
	if err != nil {
		return
	}
	userDto = &UserDto{}
	if err = copier.Copy(userDto, user); err != nil {
		return
	}
	return
}

func (service *userService) DeleteUser(ctx context.Context, userId string) (err error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}
	err = service.user.Delete(ctx, id)
	if err != nil {
		if err != repositories.UserNotFound {
			return nil
		}
	}
	return
}

func (service *userService) FindUserByID(ctx context.Context, userId string) (userDto *UserDto, err error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}
	user, err := service.user.FindByID(ctx, id)
	if err != nil || user == nil {
		return
	}

	userDto = &UserDto{}
	if err = copier.Copy(userDto, user); err != nil {
		return
	}
	return
}

func (service *userService) FindUserByAccount(ctx context.Context, account string) (userDto *UserDto, err error) {
	user, err := service.user.FindByAccount(ctx, account)
	if err != nil || user == nil {
		return
	}

	userDto = &UserDto{}
	if err = copier.Copy(userDto, user); err != nil {
		return
	}
	return
}

type RegisterUserDto struct {
	EtherAccount string `json:"etherAccount"`
}

type UserDto struct {
	UserID  string `json:"userId"`
	Account string `json:"account"`
}

func (dto *UserDto) ID(id primitive.ObjectID) {
	dto.UserID = id.Hex()
}

type UserServiceError struct {
	ServiceError
}

func NewUserServiceError(e ServiceEvent) error {
	return &UserServiceError{ServiceError{ServiceName: "UserService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}
