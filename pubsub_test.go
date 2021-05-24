package gopubsub_test

import (
	"context"
	"testing"

	_ "github.com/mshogin/gopubsub"
	"github.com/stretchr/testify/require"
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
	s := require.New(t)

	topicName := "testtopic"

	top, err := pubsub.OpenTopic(context.Background(), "local://"+topicName)
	s.NoError(err)
	s.NotNil(top)

	s.NoError(top.Send(context.Background(), &pubsub.Message{Body: []byte("test message")}))

	sub, err := pubsub.OpenSubscription(context.Background(), "local://"+topicName)
	s.NoError(err)
	s.NotNil(sub)

	cnt := 3
	for i := 0; i < cnt; i++ {
		s.NoError(top.Send(context.Background(), &pubsub.Message{Body: []byte("test message")}))
	}

	for i := 0; i < cnt; i++ {
		msg, err := sub.Receive(context.Background())
		s.NoError(err)
		s.NotNil(msg)
		s.Equal([]byte("test message"), msg.Body)
	}

	// no messages and timeout
	_, err = sub.Receive(context.Background())
	s.Error(err)
}
