package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsMiddleware interface {
	Cors() gin.HandlerFunc
}

type corsMiddleware struct {
	Config cors.Config
}

func NewCorsMiddleware() (middleware CorsMiddleware, err error) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"*"}
	return &corsMiddleware{Config: config}, nil
}

func (middleware *corsMiddleware) Cors() gin.HandlerFunc {
	return cors.New(middleware.Config)
}
