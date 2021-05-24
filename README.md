`in-memory` driver for the gocloud.dev/pubsub

[GoPubSub](https://github.com/mshogin/gopubsub)

`gopubsub` implements `gocloud.dev/pubsub/driver` interface and
provides the ability to use in-memory pubsub system.
That allows quickly switch from prod systems like kafka, gcppubsub, etc to the in-memory system. That brings a productivity boost in the development process.

## Documentation

See documentation about pubsub on https://gocloud.dev/howto/pubsub/

## Driver: Quick start

### Sending the message
```golang
import (
    "context"

    "gocloud.dev/pubsub"
    _ "github.com/mshogin/gopubsub" // driver
)
...
ctx := context.Background()
topic, err := pubsub.OpenTopic(ctx, "local://topic-name")
if err != nil {
    return fmt.Errorf("could not open topic: %v", err)
}
defer topic.Shutdown(ctx) // topic is a *pubsub.Topic;

err := topic.Send(ctx, &pubsub.Message{
	Body: []byte("Hello, World!\n"),
	// Metadata is optional and can be nil.
	Metadata: map[string]string{
		// These are examples of metadata.
		// There is nothing special about the key names.
		"language":   "en",
		"importance": "high",
	},
})
if err != nil {
	return err
}
```
### Receiving the message
```golang
import (
    "context"

    "gocloud.dev/pubsub"
    _ "github.com/mshogin/gopubsub" // driver
)
...
ctx := context.Background()
subs, err := pubsub.OpenSubscription(ctx, "local://topic-name")
if err != nil {
    return fmt.Errorf("could not open topic subscription: %v", err)
}
defer subs.Shutdown(ctx) // subs is a *pubsub.Subscription;

// Loop on received messages.
for {
	msg, err := subscription.Receive(ctx)
	if err != nil {
		// Errors from Receive indicate that Receive will no longer succeed.
		log.Printf("Receiving message: %v", err)
		break
	}
	// Do work based on the message, for example:
	fmt.Printf("Got message: %q\n", msg.Body)
	// Messages must always be acknowledged with Ack.
	msg.Ack()
}
```

### How to add support of several pubsub systems
A complete example how to use it could be found [trading/pubsub](https://github.com/mshogin/trading/blob/master/pkg/pubsub/init.go)

```golang
const (
	googlePubSub = "gcppubsub"
	kafkaPubSub  = "kafka"
	inAppPubSub  = "in-app"
)

type (
	sender interface {
		Send(context.Context, *pubsub.Message) error
	}

	receiver interface {
		Receive(context.Context) (*pubsub.Message, error)
	}

	shutdowner interface {
		Shutdown(context.Context) error
	}
)

var (
	openTopic        func(ctx context.Context, topic string) (interface{}, error)
	openSubscription func(ctx context.Context, topic string) (interface{}, error)
)

func Init(conf *config.SystemConfig) error {
	switch conf.MQProvider {
	case googlePubSub:
		Debug("using google pubsub")
		gcp.Init(conf)
		openTopic = gcp.OpenTopic
		openSubscription = gcp.OpenSubscription
	case kafkaPubSub:
		Debug("using kafka pubsub")
		kafka.Init(conf)
		openTopic = kafka.OpenTopic
		openSubscription = kafka.OpenSubscription
	case inAppPubSub:
		Debug("using in-app pubsub")
		openTopic = local.OpenTopic
		openSubscription = local.OpenSubscription
	default:
		panic("unsupported mq provider: " + conf.MQProvider)
	}
	return nil
}

```

## Installation

```sh
go get github.com/mshogin/gopubsub
```

## Contribution

Please feel free to submit any pull requests.

## Contributor List


|User|
|--|
| [mshogin](https://github.com/mshogin) |
