package pubsubs

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
)

var (
	subInstance *amqp.Subscriber
	pubInstance *amqp.Publisher
)

func GetSub() (instance *amqp.Subscriber, err error) {

	if subInstance == nil {
		instance, err = newSub()
		if err != nil {
			return nil, err
		}
		subInstance = instance
	}
	return subInstance, nil
}

func newSub() (instance *amqp.Subscriber, err error) {
	amqpConfig := amqp.NewDurableQueueConfig("amqp://ann:1213@localhost:5672/")
	instance, err = amqp.NewSubscriber(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return
	}
	return
}

func GetPub() (instance *amqp.Publisher, err error) {
	if pubInstance == nil {
		instance, err = newPub()
		if err != nil {
			return nil, err
		}
		pubInstance = instance
	}
	return pubInstance, nil
}

func newPub() (instance *amqp.Publisher, err error) {
	amqpConfig := amqp.NewDurableQueueConfig("amqp://ann:1213@localhost:5672/")
	instance, err = amqp.NewPublisher(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return
	}
	return
}
