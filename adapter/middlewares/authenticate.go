package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/business/services"
	"nftshopping-store-api/pkg/security"
	"strings"
)

type AuthenticateMiddleware interface {
	Authenticate() gin.HandlerFunc
}

type authenticateMiddleware struct {
	auth AuthService
}

func NewAuthenticateMiddleware() (middleware AuthenticateMiddleware, err error) {
	service, err := services.GetService()
	if err != nil {
		return nil, err
	}
	return &authenticateMiddleware{auth: service.Auth}, nil
}

func (middleware *authenticateMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.AbortWithStatus(401)
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.AbortWithStatus(401)
			return
		}
		name, err := security.ExtractUserName(clientToken)
		if err != nil {
			return
		}
		user, err := middleware.auth.FindAuthByName(context.Background(), name)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		if user == nil {
			c.AbortWithStatus(401)
			return
		}

		isValid, err := security.ValidateToken(clientToken, user)
		if err != nil {
			return
		}

		if !isValid {
			c.AbortWithStatus(401)
			return
		}

		authentication := &authentication{user.GetName(), user.GetAuthorities()}
		c.Set("Authentication", authentication)
		c.Next()
	}
}

type AuthService interface {
	FindAuthByName(ctx context.Context, userName string) (security.Authentication, error)
}

type authentication struct {
	UserName    string
	Authorities []string
}

func (auth *authentication) GetName() string {
	return auth.UserName
}

func (auth *authentication) GetAuthorities() []string {
	return auth.Authorities
}
