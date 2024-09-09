/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package eventbus

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/commcos/utils/grpool"
	"github.com/commcos/utils/logger"
)

// EventBus - box for handlers and callbacks.
type EventBus struct {
	handlers     map[EventTopic][]*eventHandler
	lock         sync.Mutex // a lock for the map
	wg           sync.WaitGroup
	evtAsyncPool *grpool.Pool
}

type evtAsyncPayload struct {
	args    []interface{}
	handler *eventHandler
	topic   EventTopic
}

type eventHandler struct {
	callBack      reflect.Value
	flagOnce      bool
	async         bool
	transactional bool
	sync.Mutex    // lock for an event handler - useful for running async callbacks serially
}

// New returns new EventBus with empty handlers.
func New(asyncPoolSize int) Interface {
	b := &EventBus{
		handlers: make(map[EventTopic][]*eventHandler),
		lock:     sync.Mutex{},
		wg:       sync.WaitGroup{},
	}

	b.evtAsyncPool = newEvtBusPool(asyncPoolSize, b.doEvtAsyncPoolExecFunc)

	return Interface(b)
}

func (bus *EventBus) doEvtAsyncPoolExecFunc(payload interface{}) interface{} {
	param, ok := payload.(*evtAsyncPayload)
	if !ok {
		logger.Log(logger.ErrorLevel, "payload type = ", reflect.TypeOf(payload))
		return fmt.Errorf("解析payload类型失败")
	}

	bus.doPublishAsync(param.handler, param.topic, param.args...)

	return nil
}

// doSubscribe handles the subscription logic and is utilized by the public Subscribe functions
func (bus *EventBus) doSubscribe(topic EventTopic, fn interface{}, handler *eventHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	bus.handlers[topic] = append(bus.handlers[topic], handler)
	return nil
}

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (bus *EventBus) Subscribe(topic EventTopic, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), false, false, false, sync.Mutex{},
	})
}

// SubscribeAsync subscribes to a topic with an asynchronous callback
// Transactional determines whether subsequent callbacks for a topic are
// run serially (true) or concurrently (false)
// Returns error if `fn` is not a function.
func (bus *EventBus) SubscribeAsyncPool(topic EventTopic, fn interface{}, transactional bool) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), false, true, transactional, sync.Mutex{},
	})
}

// SubscribeOnce subscribes to a topic once. Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *EventBus) SubscribeOnce(topic EventTopic, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), true, false, false, sync.Mutex{},
	})
}

// SubscribeOnceAsync subscribes to a topic once with an asynchronous callback
// Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *EventBus) SubscribeOnceAsync(topic EventTopic, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), true, true, false, sync.Mutex{},
	})
}

// HasCallback returns true if exists any callback subscribed to the topic.
func (bus *EventBus) HasCallback(topic EventTopic) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	_, ok := bus.handlers[topic]
	if ok {
		return len(bus.handlers[topic]) > 0
	}
	return false
}

// Unsubscribe removes callback defined for a topic.
// Returns error if there are no callbacks subscribed to the topic.
func (bus *EventBus) Unsubscribe(topic EventTopic, handler interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if _, ok := bus.handlers[topic]; ok && len(bus.handlers[topic]) > 0 {
		bus.removeHandler(topic, bus.findHandlerIdx(topic, reflect.ValueOf(handler)))
		return nil
	}
	return fmt.Errorf("topic %s doesn't exist", topic)
}

// Publish executes callback defined for a topic. Any additional argument will be transferred to the callback.
func (bus *EventBus) Publish(topic EventTopic, args ...interface{}) {
	bus.lock.Lock() // will unlock if handler is not found or always after setUpPublish
	defer bus.lock.Unlock()
	if handlers, ok := bus.handlers[topic]; ok && 0 < len(handlers) {
		// Handlers slice may be changed by removeHandler and Unsubscribe during iteration,
		// so make a copy and iterate the copied slice.
		copyHandlers := make([]*eventHandler, len(handlers))
		copy(copyHandlers, handlers)
		for i, handler := range copyHandlers {
			if handler.flagOnce {
				bus.removeHandler(topic, i)
			}
			if !handler.async {
				bus.doPublish(handler, topic, args...)
			} else {
				bus.wg.Add(1)
				if handler.transactional {
					bus.lock.Unlock()
					handler.Lock()
					bus.lock.Lock()
				}
				payload := &evtAsyncPayload{
					handler: handler,
					topic:   topic,
					args:    args,
				}
				bus.evtAsyncPool.AsyncProcess(payload)
			}
		}
	}
}

func (bus *EventBus) doPublish(handler *eventHandler, topic EventTopic, args ...interface{}) {
	passedArguments := bus.setUpPublish(handler, args...)
	logger.Log(logger.DebugLevel, "call user function with args-len(%v)", len(passedArguments))
	handler.callBack.Call(passedArguments)
}

func (bus *EventBus) doPublishAsync(handler *eventHandler, topic EventTopic, args ...interface{}) {
	defer bus.wg.Done()
	if handler.transactional {
		defer handler.Unlock()
	}
	bus.doPublish(handler, topic, args...)
}

func (bus *EventBus) removeHandler(topic EventTopic, idx int) {
	if _, ok := bus.handlers[topic]; !ok {
		return
	}
	l := len(bus.handlers[topic])

	if !(0 <= idx && idx < l) {
		return
	}

	copy(bus.handlers[topic][idx:], bus.handlers[topic][idx+1:])
	bus.handlers[topic][l-1] = nil // or the zero value of T
	bus.handlers[topic] = bus.handlers[topic][:l-1]
}

func (bus *EventBus) findHandlerIdx(topic EventTopic, callback reflect.Value) int {
	if _, ok := bus.handlers[topic]; ok {
		for idx, handler := range bus.handlers[topic] {
			if handler.callBack.Type() == callback.Type() &&
				handler.callBack.Pointer() == callback.Pointer() {
				return idx
			}
		}
	}
	return -1
}

func (bus *EventBus) setUpPublish(callback *eventHandler, args ...interface{}) []reflect.Value {
	funcType := callback.callBack.Type()
	passedArguments := make([]reflect.Value, len(args))
	passedArgumentsIndex := 0

	logger.Log(logger.DebugLevel, "setup publish %v", len(args))

	for i, v := range args {
		if v == nil {
			passedArguments[passedArgumentsIndex] = reflect.New(funcType.In(i)).Elem()
		} else {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Slice:
				if refVal, ok := v.([]reflect.Value); ok {

					passedRefletSliceArguments := make([]reflect.Value, len(refVal))
					for aj, refValdeep := range refVal {
						logger.Log(logger.DebugLevel, "setup publish array %v", aj)
						passedRefletSliceArguments[aj] = refValdeep
					}

					reflectSliceSize := len(passedRefletSliceArguments)
					if reflectSliceSize > 0 {
						passedArguments = append([]reflect.Value(nil), passedRefletSliceArguments...)
						passedArgumentsIndex += reflectSliceSize
					}
				} else {
					logger.Log(logger.WarnLevel, "setup parameter input a slice but not a reflect value. ignore?")
				}
			default:
				passedArguments[passedArgumentsIndex] = reflect.ValueOf(v)
			}
		}
		passedArgumentsIndex++
	}

	return passedArguments
}

// WaitAsync waits for all async callbacks to complete
func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}
