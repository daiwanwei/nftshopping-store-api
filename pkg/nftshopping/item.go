package nftshopping

import (
	"items"
	"nftshopping-store-api/pkg/config"
)

var itemInstance *items.Client

func GetItem() (instance *items.Client, err error) {
	if itemInstance == nil {
		instance, err = newItem()
		if err != nil {
			return nil, err
		}
		itemInstance = instance
	}
	return itemInstance, nil
}

func newItem() (instance *items.Client, err error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	itemConfig := c.Item
	domain := itemConfig.Domain
	return items.NewClient(domain), nil
}
