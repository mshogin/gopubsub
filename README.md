# `gopubsub`, the in-app publish/subscribe system.

[![GoPubSub](https://github.com/mshogin/gopubsub)](https://github.com/mshogin/gopubsub)

`gopubsub` implements interface of `gocloud.dev/pubsub` allows to develop a pubsub based applications
in local environment without interacting with the real pubsub services like kafka, google pubsub, etc

## Quick start

```golang
func Init(conf *config.SystemConfig) error {
	switch conf.MQProvider {
	case inAppPubSub:
		openTopic = gopubsub.OpenTopic
		openSubscription = gopubsub.OpenSubscription
	case googlePubSub:
		gcp.Init(conf)
		openTopic = gcp.OpenTopic
		openSubscription = gcp.OpenSubscription
	case kafkaPubSub:
		kafka.Init(conf)
		openTopic = kafka.OpenTopic
		openSubscription = kafka.OpenSubscription
	default:
		panic("unsupported mq provider: " + conf.MQProvider)
	}
	return nil
}


google pubsub
var openTopicURI string
var subscribeURI string

func Init(conf *config.SystemConfig) {
	openTopicURI = fmt.Sprintf(
		"gcppubsub://projects/%s/topics/",
		conf.ProjectName)
	subscribeURI = fmt.Sprintf(
		"gcppubsub://projects/%s/subscriptions/",
		conf.ProjectName)
}

var OpenTopic = func(ctx context.Context, topic string) (interface{}, error) {
	return pubsub.OpenTopic(ctx, openTopicURI+topic)
}

var OpenSubscription = func(ctx context.Context, topic string) (interface{}, error) {
	return pubsub.OpenSubscription(ctx, subscribeURI+topic+"-sub")
}

kafka
var subscriptionGroups = map[string]uint64{}
var subscriptionGroupsSync sync.Mutex
var openTopicURI = "kafka://"
var subscribeURI = "kafka://group-%d?topic=%s"

func Init(conf *config.SystemConfig) {}

var OpenTopic = func(ctx context.Context, topic string) (interface{}, error) {
	return pubsub.OpenTopic(ctx, openTopicURI+topic)
}

var OpenSubscription = func(ctx context.Context, topic string) (interface{}, error) {
	subscriptionGroupsSync.Lock()
	defer subscriptionGroupsSync.Unlock()
	groupID := subscriptionGroups[topic]
	subscriptionGroups[topic] = groupID + 1

	uri := fmt.Sprintf(subscribeURI, groupID, topic)
	return pubsub.OpenSubscription(ctx, uri)
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
