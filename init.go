package gopubsub

import (
	"context"
	"fmt"
	"net/url"

	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/batcher"
)

func init() {
	opener := new(URLOpener)
	pubsub.DefaultURLMux().RegisterTopic(Scheme, opener)
	pubsub.DefaultURLMux().RegisterSubscription(Scheme, opener)
}

// Scheme is the URL scheme that kafkapubsub registers its URLOpeners under on pubsub.DefaultMux.
const Scheme = "local"

// URLOpener opens Kafka URLs like "kafka://mytopic" for topics and
// "kafka://group?topic=mytopic" for subscriptions.
//
// For topics, the URL's host+path is used as the topic name.
//
// For subscriptions, the URL's host+path is used as the group name,
// and the "topic" query parameter(s) are used as the set of topics to
// subscribe to.
type URLOpener struct{}

var sendBatcherOpts = &batcher.Options{
	MaxBatchSize: 100,
	MaxHandlers:  2,
}

// OpenTopicURL opens a pubsub.Topic.
func (o *URLOpener) OpenTopic(ctx context.Context, topicName string) (*pubsub.Topic, error) {
	t, err := openTopic(ctx, topicName)
	if err != nil {
		return nil, fmt.Errorf("cannot open local topic: %w", err)
	}
	return pubsub.NewTopic(t, sendBatcherOpts), nil
}

func (o *URLOpener) OpenTopicURL(ctx context.Context, u *url.URL) (*pubsub.Topic, error) {
	return o.OpenTopic(ctx, u.Host)
}

var recvBatcherOpts = &batcher.Options{
	MaxBatchSize: 1,
	MaxHandlers:  1,
}

// OpenSubscriptionURL opens a pubsub.Subscription.
func (o *URLOpener) OpenSubscription(ctx context.Context, topicName string) (*pubsub.Subscription, error) {
	s, err := OpenSubscription(ctx, topicName)
	if err != nil {
		return nil, fmt.Errorf("cannot open local subscription: %w", err)
	}
	return pubsub.NewSubscription(s, recvBatcherOpts, nil), nil
}

func (o *URLOpener) OpenSubscriptionURL(ctx context.Context, u *url.URL) (*pubsub.Subscription, error) {
	return o.OpenSubscription(ctx, u.Host)
}
