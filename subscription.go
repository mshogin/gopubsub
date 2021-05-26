package gopubsub

import (
	"context"

	"github.com/google/uuid"
	"gocloud.dev/gcerrors"
	"gocloud.dev/pubsub/driver"
)

const defaultPublishChannelSize = 10000

type subscription struct {
	ID uuid.UUID
	C  chan *driver.Message

	unsubscribe func()
}

func OpenSubscription(ctx context.Context, topicName string) (*subscription, error) {
	t, _ := openTopic(ctx, topicName)
	s := subscription{
		ID: uuid.New(),
		C:  make(chan *driver.Message, defaultPublishChannelSize),
	}
	s.unsubscribe = t.add(s)
	return &s, nil
}

func (m subscription) ReceiveBatch(ctx context.Context, n int) ([]*driver.Message, error) {
	// c, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	// defer cancel()
	for {
		select {
		// case <-c.Done():
		// 	return nil, errors.New("receive message timeout")
		case dm := <-m.C:
			return []*driver.Message{dm}, nil
		}
	}
}

// As implements driver.Subscription.As.
func (s *subscription) As(i interface{}) bool {
	return true
}

// CanNack implements driver.CanNack.
func (s *subscription) CanNack() bool {
	// Nacking a single message doesn't make sense with the way Kafka maintains
	// offsets.
	return false
}

// SendNacks implements driver.Subscription.SendNacks.
func (s *subscription) SendNacks(ctx context.Context, ids []driver.AckID) error {
	panic("unreachable")
}

// Close implements io.Closer.
func (s *subscription) Close() error {
	// Cancel the ctx for the background goroutine and wait until it's done.
	s.unsubscribe()
	return nil
}

// ErrorAs implements driver.Subscription.ErrorAs.
func (s *subscription) ErrorAs(err error, i interface{}) bool {
	return true
}

// ErrorCode implements driver.Subscription.ErrorCode.
func (*subscription) ErrorCode(err error) gcerrors.ErrorCode {
	return gcerrors.Internal
}

// IsRetryable implements driver.Subscription.IsRetryable.
func (*subscription) IsRetryable(error) bool {
	return false
}

// SendAcks implements driver.Subscription.SendAcks.
func (s *subscription) SendAcks(ctx context.Context, ids []driver.AckID) error {
	return nil
}
