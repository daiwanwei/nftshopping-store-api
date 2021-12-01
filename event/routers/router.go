package routers

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"time"
)

var (
	routerInstance *message.Router
	logger         = watermill.NewStdLogger(false, false)
)

func GetRouter() (instance *message.Router, err error) {
	if routerInstance == nil {
		instance, err = newRouter()
		if err != nil {
			return nil, err
		}
		routerInstance = instance
	}
	return routerInstance, nil
}

func newRouter() (router *message.Router, err error) {
	router, err = message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return
	}
	router.AddPlugin(plugin.SignalsHandler)

	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)
	err = InitItemRouter(router)
	if err != nil {
		return
	}
	return router, nil
}
