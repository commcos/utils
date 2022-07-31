/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package eventbus

import (
	"time"

	"github.com/xsbull/utils/grpool"
	"github.com/xsbull/utils/logger"
)

var (
	//协程重新可用的空闲间隔
	idleTimeDuration = 2 * time.Second
)

type evtBusClient struct {
	processor func(interface{}) interface{}
}

func newEvtBusPool(poolSize int, f func(interface{}) interface{}) *grpool.Pool {

	return grpool.NewWithTimeout(poolSize, func() grpool.Worker {
		client := newEvtBusClient()
		client.processor = f
		return client
	}, idleTimeDuration)
}

func newEvtBusClient() *evtBusClient {
	a := &evtBusClient{}

	return a
}

func (ac *evtBusClient) Process(payload interface{}) interface{} {
	return ac.processor(payload)
}

func (ac *evtBusClient) BlockUntilReady() {

}

func (ac *evtBusClient) Interrupt() {
}

func (ac *evtBusClient) Terminate() {
	logger.Log(logger.InfoLevel, "Terminate")
}
