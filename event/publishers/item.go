package publishers

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"nftshopping-store-api/event"
	"nftshopping-store-api/event/messages"
	"nftshopping-store-api/pkg/pubsubs"
)

type ItemPublisher interface {
	PublishToOrderItem(msg messages.OrderItemMessage) (err error)
	PublishToDeliverItem(msg messages.DeliverItemMessage) (err error)
}

type itemPublisher struct {
	pub *amqp.Publisher
}

func NewItemPublisher() (ItemPublisher, error) {
	pub, err := pubsubs.GetPub()
	if err != nil {
		return nil, err
	}
	return &itemPublisher{
		pub: pub,
	}, nil
}

func (publisher *itemPublisher) PublishToOrderItem(msg messages.OrderItemMessage) (err error) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return
	}
	m := message.NewMessage(watermill.NewUUID(), msgByte)
	if err := publisher.pub.Publish(event.OrderItem, m); err != nil {
		return err
	}
	return
}

func (publisher *itemPublisher) PublishToDeliverItem(msg messages.DeliverItemMessage) (err error) {

	msgByte, err := json.Marshal(msg)
	if err != nil {
		return
	}
	m := message.NewMessage(watermill.NewUUID(), msgByte)
	if err := publisher.pub.Publish(event.DeliverItem, m); err != nil {
		return err
	}
	return
}
