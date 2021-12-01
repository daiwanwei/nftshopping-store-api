package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"nftshopping-store-api/business/services"
	"nftshopping-store-api/event/messages"
	"nftshopping-store-api/pkg/log"
)

type ItemHandler interface {
	ListenDeliverItem(msg *message.Message) (err error)
	ListenOrderItem(msg *message.Message) (err error)
}

type itemHandler struct {
	logger log.Logger
	item   services.ItemService
}

func NewItemHandler() (handler ItemHandler, err error) {
	service, err := services.GetService()
	if err != nil {
		return nil, err
	}
	logger, err := log.GetLog()
	if err != nil {
		return nil, err
	}
	return &itemHandler{
		item:   service.Item,
		logger: logger,
	}, nil
}

func (handler *itemHandler) ListenDeliverItem(msg *message.Message) (err error) {
	fmt.Printf("received message: %s, payload: %s \n", msg.UUID, string(msg.Payload))
	var m messages.DeliverItemMessage
	err = json.Unmarshal(msg.Payload, &m)
	if err != nil {
		return
	}
	item, err := handler.item.DeliverItem(context.Background(), services.DeliverItemDto{
		Contract:   m.Contract,
		Token:      m.Token,
		CreationId: m.CreationId,
	})
	if err != nil {
		handler.logger.Error(err)
		return nil
	}
	fmt.Printf("received item:\n contract(%s),\n token(%s),\n", item.Contract, item.Token)
	return
}

func (handler *itemHandler) ListenOrderItem(msg *message.Message) (err error) {
	fmt.Printf("received message: %s, payload: %s \n", msg.UUID, string(msg.Payload))
	var m messages.OrderItemMessage
	err = json.Unmarshal(msg.Payload, &m)
	if err != nil {
		return
	}
	fmt.Printf("order item:\n creationId(%s),\n amount(%d),\n", m.CreationId, m.Amount)
	return
}
