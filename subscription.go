package gopubsub

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gocloud.dev/pubsub"
)

const defaultPublishChannelSize = 10000

type subscription struct {
	ID uuid.UUID
	C  chan *pubsub.Message

	unsubscribe func()
}

func OpenSubscription(ctx context.Context, topicName string) (interface{}, error) {
	t, _ := OpenTopic(ctx, topicName)
	s := subscription{
		ID: uuid.New(),
		C:  make(chan *pubsub.Message, defaultPublishChannelSize),
	}
	s.unsubscribe = t.(*topic).add(s)
	return s, nil
}

func (m subscription) Receive(ctx context.Context) (*pubsub.Message, error) {
	c, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	for {
		select {
		case <-c.Done():
			return nil, errors.New("receive message timeout")
		case m := <-m.C:
			return m, nil
		}
	}
}

func (m subscription) Shutdown(context.Context) error {
	m.unsubscribe()
	return nil
}
