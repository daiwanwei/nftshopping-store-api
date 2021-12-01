package subscribers

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"nftshopping-store-api/event"
	"nftshopping-store-api/event/handlers"
	"nftshopping-store-api/pkg/pubsubs"
)

type ItemSubscriber interface {
	Run() error
	Stop() error
}

type itemSubscriber struct {
	orderItemCh   <-chan *message.Message
	deliverItemCh <-chan *message.Message
	stopCh        chan struct{}
	item          handlers.ItemHandler
}

func NewItemSubscriber() (subscriber ItemSubscriber, err error) {
	handler, err := handlers.GetHandler()
	if err != nil {
		return nil, err
	}
	sub, err := pubsubs.GetSub()
	if err != nil {
		return nil, err
	}
	deliverItemCh, err := sub.Subscribe(context.Background(), event.DeliverItem)
	if err != nil {
		return
	}
	orderItemCh, err := sub.Subscribe(context.Background(), event.OrderItem)
	if err != nil {
		return
	}
	return &itemSubscriber{
		item:          handler.Item,
		deliverItemCh: deliverItemCh,
		orderItemCh:   orderItemCh,
		stopCh:        make(chan struct{}, 20),
	}, nil
}

func (subscriber *itemSubscriber) Run() error {
	for {
		select {
		case msg := <-subscriber.deliverItemCh:
			err := subscriber.item.ListenDeliverItem(msg)
			if err != nil {
				msg.Nack()
				continue
			}
			msg.Ack()
		case msg := <-subscriber.orderItemCh:
			err := subscriber.item.ListenOrderItem(msg)
			if err != nil {
				msg.Nack()
				continue
			}
			msg.Ack()
		case <-subscriber.stopCh:
			return nil
		}
	}

}

func (subscriber *itemSubscriber) Stop() error {
	close(subscriber.stopCh)
	return nil
}
