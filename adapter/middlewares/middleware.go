package middlewares

var middlewareInstance *middleware

func GetMiddleware() (instance *middleware, err error) {
	if middlewareInstance == nil {
		instance, err = NewMiddleware()
		if err != nil {
			return nil, err
		}
		middlewareInstance = instance
	}
	return middlewareInstance, nil
}

type middleware struct {
	Cors CorsMiddleware
}

func NewMiddleware() (instance *middleware, err error) {
	cors, err := NewCorsMiddleware()
	if err != nil {
		return nil, err
	}
	return &middleware{
		Cors: cors,
	}, nil
}
