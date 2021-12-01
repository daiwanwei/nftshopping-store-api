package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	Authorize() gin.HandlerFunc
	Auth() gin.HandlersChain
}

type authMiddleware struct {
	authenticate AuthenticateMiddleware
	authorize    AuthorizeMiddleware
}

func NewAuthMiddleware() (middleware AuthMiddleware, err error) {
	authenticate, err := NewAuthenticateMiddleware()
	if err != nil {
		return
	}
	authorize, err := NewAuthorizeMiddleware()
	if err != nil {
		return
	}
	return &authMiddleware{
		authenticate: authenticate,
		authorize:    authorize,
	}, nil
}

func (middleware *authMiddleware) Authenticate() gin.HandlerFunc {
	return middleware.authenticate.Authenticate()
}

func (middleware *authMiddleware) Authorize() gin.HandlerFunc {
	return middleware.authorize.Authorize()
}

func (middleware *authMiddleware) Auth() gin.HandlersChain {
	return []gin.HandlerFunc{middleware.Authenticate(), middleware.Authorize()}
}
