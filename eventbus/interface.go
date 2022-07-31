/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package eventbus

type EventTopic string

//EventBusSubscriber defines subscription-related bus behavior
type EventBusSubscriber interface {
	Subscribe(topic EventTopic, fn interface{}) error
	//transactional 为true，代表顺序化执行callback
	SubscribeAsyncPool(topic EventTopic, fn interface{}, transactional bool) error
	SubscribeOnce(topic EventTopic, fn interface{}) error
	SubscribeOnceAsync(topic EventTopic, fn interface{}) error
	Unsubscribe(topic EventTopic, handler interface{}) error
}

//BusPublisher defines publishing-related bus behavior
type EventBusPublisher interface {
	Publish(topic EventTopic, args ...interface{})
}

//EventBus englobes global (subscribe, publish, control) bus behavior
type Interface interface {
	EventBusSubscriber
	EventBusPublisher

	HasCallback(topic EventTopic) bool
	WaitAsync()
}
