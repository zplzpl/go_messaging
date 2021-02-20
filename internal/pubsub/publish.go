package pubsub

type Publisher interface {
	Publish(msgKey string, buf []byte) error
}

var _defaultPublisher Publisher

func InitDefaultPublisher(p Publisher) {
	_defaultPublisher = p
}

func GetPublisher() Publisher {

	if _defaultPublisher != nil {
		return _defaultPublisher
	}

	return &emptyPublisher{}
}

type emptyPublisher struct {
}

func (e *emptyPublisher) Publish(msgKey string, buf []byte) error {
	return nil
}
