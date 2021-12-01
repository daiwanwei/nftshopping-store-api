package services

import (
	"context"
	"nftshopping-store-api/persistence/repositories"
	"nftshopping-store-api/pkg/security"
)

type AuthService interface {
	FindAuthByName(ctx context.Context, name string) (authDetail security.Authentication, err error)
}

type authService struct {
	auth repositories.AuthDao
}

func NewAuthService() (service AuthService, err error) {
	dao, err := repositories.GetRepository()
	if err != nil {
		return nil, err
	}
	return &authService{
		auth: dao.Auth,
	}, nil
}

func (service *authService) FindAuthByName(ctx context.Context, name string) (auth security.Authentication, err error) {
	auth, err = service.auth.FindByName(ctx, name)
	if err != nil {
		return
	}
	return
}

type AuthServiceError struct {
	ServiceError
}

func NewAuthServiceError(e ServiceEvent) error {
	return &AuthServiceError{ServiceError{ServiceName: "AuthService", Code: e.GetEvent().Code, Msg: e.GetEvent().Msg, Err: nil}}
}
