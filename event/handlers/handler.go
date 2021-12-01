package handlers

var handlerInstance *handler

func GetHandler() (instance *handler, err error) {
	if handlerInstance == nil {
		instance, err = newHandler()
		if err != nil {
			return nil, err
		}
		handlerInstance = instance
	}
	return handlerInstance, nil
}

type handler struct {
	Item ItemHandler
}

func newHandler() (instance *handler, err error) {
	item, err := NewItemHandler()
	if err != nil {
		return
	}

	return &handler{
		Item: item,
	}, nil
}
