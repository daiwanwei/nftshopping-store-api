package publishers

var publisherInstance *publisher

func GetPublisher() (instance *publisher, err error) {
	if publisherInstance == nil {
		instance, err = newPublisher()
		if err != nil {
			return nil, err
		}
		publisherInstance = instance
	}
	return publisherInstance, nil
}

type publisher struct {
	Item ItemPublisher
}

func newPublisher() (instance *publisher, err error) {
	item, err := NewItemPublisher()
	if err != nil {
		return
	}

	return &publisher{
		Item: item,
	}, nil
}
