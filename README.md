# `gopubsub`, the in-app publish/subscribe system.

[GoPubSub](https://github.com/mshogin/gopubsub)

`gopubsub` implements interface of `gocloud.dev/pubsub` allows to develop a pubsub based applications
in local environment without interacting with the real pubsub services like kafka, google pubsub, etc

## Quick start

A complete example how to use it could be found [trading/pubsub](https://github.com/mshogin/trading/blob/master/pkg/pubsub/init.go)

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
