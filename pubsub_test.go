package gopubsub_test

import (
	"context"
	"testing"

	"github.com/mshogin/gopubsub"
	"github.com/stretchr/testify/assert"
	"gocloud.dev/pubsub"
)

type (
	sender interface {
		Send(ctx context.Context, msg *pubsub.Message) error
	}

	receiver interface {
		Receive(context.Context) (*pubsub.Message, error)
	}
)

func TestPubsub(t *testing.T) {
	s := assert.New(t)
	gopubsub.Init()

	topicName := "test topic"

	top, err := gopubsub.OpenTopic(context.Background(), topicName)
	s.NoError(err)
	s.NotNil(top)

	s.NoError(top.(sender).Send(context.Background(), &pubsub.Message{Body: []byte("test message")}))

	sub, err := gopubsub.OpenSubscription(context.Background(), topicName)
	s.NoError(err)
	s.NotNil(sub)

	cnt := 3
	for i := 0; i < cnt; i++ {
		s.NoError(top.(sender).Send(context.Background(), &pubsub.Message{Body: []byte("test message")}))
	}

	for i := 0; i < cnt; i++ {
		msg, err := sub.(receiver).Receive(context.Background())
		s.NoError(err)
		s.NotNil(msg)
		s.Equal([]byte("test message"), msg.Body)
	}

	// no messages and timeout
	_, err = sub.(receiver).Receive(context.Background())
	s.Error(err)
}
