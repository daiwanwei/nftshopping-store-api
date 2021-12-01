package routers

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"nftshopping-store-api/event"
	"nftshopping-store-api/event/handlers"
	"nftshopping-store-api/pkg/pubsubs"
)

func InitItemRouter(router *message.Router) (err error) {
	sub, err := pubsubs.GetSub()
	if err != nil {
		return
	}
	handler, err := handlers.GetHandler()
	if err != nil {
		return
	}
	router.AddNoPublisherHandler(
		"DeliverItem",
		event.DeliverItem,
		sub,
		handler.Item.ListenDeliverItem,
	)
	router.AddNoPublisherHandler(
		"OrderItem",
		event.OrderItem,
		sub,
		handler.Item.ListenOrderItem,
	)
	return
}
