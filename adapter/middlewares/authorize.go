package middlewares

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"nftshopping-store-api/pkg/casbins"
	"nftshopping-store-api/pkg/security"
)

type AuthorizeMiddleware interface {
	Authorize() gin.HandlerFunc
}

type authorizeMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizeMiddleware() (middleware AuthorizeMiddleware, err error) {
	enforcer, err := casbins.GetEnforcer()
	if err != nil {
		return nil, err
	}
	return &authorizeMiddleware{enforcer: enforcer}, nil
}

func (middleware *authorizeMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		authentication, isExist := c.Get("Authentication")
		if !isExist {
			c.AbortWithStatus(403)
			return
		}
		auth, ok := authentication.(security.Authentication)
		if !ok {
			c.AbortWithStatus(500)
			return
		}
		isAuthorized := false
		for _, role := range auth.GetAuthorities() {
			ok, err := middleware.enforcer.EnforceSafe(role, c.Request.URL.Path, c.Request.Method)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}
			if ok {
				isAuthorized = true
				break
			}
		}
		if !isAuthorized {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}
