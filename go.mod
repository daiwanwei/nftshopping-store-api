module nftshopping-store-api

go 1.16

require (
	github.com/ThreeDotsLabs/watermill v1.2.0-rc.7
	github.com/ThreeDotsLabs/watermill-amqp v1.1.4
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/casbin/casbin v1.9.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/sentry-go v0.11.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.1
	github.com/go-redis/redis/v8 v8.11.1
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/golang-lru v0.5.1
	github.com/jinzhu/copier v0.3.0
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	go.mongodb.org/mongo-driver v1.5.2
	go.uber.org/zap v1.16.0
	golang.org/x/mod v0.3.1-0.20200828183125-ce943fd02449 // indirect
	items v0.0.1
)

replace items v0.0.1 => ./items
