package gopubsub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gocloud.dev/pubsub"
)

type topic struct {
	name     string
	subs     map[uuid.UUID]subscription
	subsSync sync.RWMutex
}

var topics = make(map[string]*topic)
var topicsSync sync.Mutex

func OpenTopic(ctx context.Context, topicName string) (interface{}, error) {
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

func (m *topic) Send(ctx context.Context, msg *pubsub.Message) (err error) {
	m.subsSync.RLock()
	defer m.subsSync.RUnlock()

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
				case s.C <- msg:
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
