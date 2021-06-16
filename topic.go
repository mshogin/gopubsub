package gopubsub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gocloud.dev/gcerrors"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/driver"
)

type ErrorCode = gcerrors.ErrorCode

const (
	PubsubLocal = 1001
)

type topic struct {
	name     string
	subs     map[uuid.UUID]subscription
	subsSync sync.RWMutex
}

var topics = make(map[string]*topic)
var topicsSync sync.Mutex

func openTopic(ctx context.Context, topicName string) (*topic, error) {
	topicsSync.Lock()
	defer topicsSync.Unlock()
	if topics[topicName] == nil {
		topics[topicName] = &topic{
			name: topicName,
			subs: make(map[uuid.UUID]subscription),
		}
	}
	return topics[topicName], nil
}

func (m *topic) SendBatch(ctx context.Context, msgs []*driver.Message) error {
	m.subsSync.RLock()
	defer m.subsSync.RUnlock()

	for idx, msg := range msgs {
		if err := m.send(ctx, &pubsub.Message{
			Body: msg.Body,
		}); err != nil {
			return fmt.Errorf("cannot send message #%d; messages count: %d: %w", idx, len(msgs), err)
		}
	}
	return nil
}

func (m *topic) send(ctx context.Context, msg *pubsub.Message) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(m.subs))
	for _, s := range m.subs {
		s := s
		go func() {
			c, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
			defer cancel()
			running := true
			for running {
				select {
				case s.C <- &driver.Message{Body: msg.Body}:
					running = false
				case <-c.Done():
					running = false
					if err == nil {
						err = fmt.Errorf("cannot send message to subscription: %v", s)
					} else {
						err = fmt.Errorf("cannot send message to subscription: %v: %w", s, err)
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return
}

// Close implements io.Closer.
func (t *topic) Close() error {
	return nil
}

// IsRetryable implements driver.Topic.IsRetryable.
func (t *topic) IsRetryable(error) bool {
	return false
}

// As implements driver.Topic.As.
func (t *topic) As(i interface{}) bool {
	panic("not sure what for it is")
	// if p, ok := i.(*sarama.SyncProducer); ok {
	// 	*p = t.producer
	// 	return true
	// }
	// return false
}

// ErrorAs implements driver.Topic.ErrorAs.
func (t *topic) ErrorAs(err error, i interface{}) bool {
	panic("not sure what for it is")
	// return errorAs(err, i)
}

// ErrorCode implements driver.Topic.ErrorCode.
func (t *topic) ErrorCode(err error) gcerrors.ErrorCode {
	fmt.Printf("pubsub local: %s", err)
	return PubsubLocal
}

func (m *topic) add(s subscription) func() {
	m.subsSync.Lock()
	defer m.subsSync.Unlock()
	m.subs[s.ID] = s
	return func() {
		m.subsSync.Lock()
		defer m.subsSync.Unlock()
		delete(m.subs, s.ID)
	}
}
