package services

import "encoding/json"

var serviceInstance *service

func GetService() (instance *service, err error) {
	if serviceInstance == nil {
		instance, err = newService()
		if err != nil {
			return nil, err
		}
		serviceInstance = instance
	}
	return serviceInstance, nil
}

type service struct {
	Auth       AuthService
	User       UserService
	Creation   CreationService
	Item       ItemService
	Collection CollectionService
	Trade      TradeService
	Brand      BrandService
	Stock      StockService
}

func newService() (instance *service, err error) {
	auth, err := NewAuthService()
	if err != nil {
		return
	}
	user, err := NewUserService()
	if err != nil {
		return
	}
	brand, err := NewBrandService()
	if err != nil {
		return
	}
	creation, err := NewCreationService(brand)
	if err != nil {
		return
	}
	item, err := NewItemService(creation, brand, user)
	if err != nil {
		return
	}
	collection, err := NewCollectionService(user, creation)
	if err != nil {
		return
	}
	trade, err := NewTradeService(user, creation)
	if err != nil {
		return
	}
	stock, err := NewStockService(brand, creation)
	if err != nil {
		return
	}

	return &service{
		Auth:       auth,
		User:       user,
		Creation:   creation,
		Item:       item,
		Collection: collection,
		Trade:      trade,
		Brand:      brand,
		Stock:      stock,
	}, nil
}

func generateKeyOfCache(components ...interface{}) (key string, err error) {
	key = "query"
	for i := range components {
		strByte, err := json.Marshal(components[i])
		if err != nil {
			return "", err
		}
		key = key + "+" + string(strByte)
	}
	return key, nil
}
