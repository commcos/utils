/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package grpool

import (
	"time"
)

//------------------------------------------------------------------------------

type jobData struct {
	//send the payload to this worker
	payload interface{}
	//asyncProcess asyncProcess flag
	asyncProcess bool
	//asyncJobComplete 异步处理完成通知回调
	asyncJobComplete func(result interface{})
}

// workRequest is a struct containing context representing a workers intention
// to receive a work payload.
type workRequest struct {
	// jobChan is used to send the payload to this worker.
	jobChan chan<- jobData

	// retChan is used to read the result from this worker.
	retChan <-chan interface{}

	// interruptFunc can be called to cancel a running job. When called it is no
	// longer necessary to read from retChan.
	interruptFunc func()
}

//------------------------------------------------------------------------------

// workerWrapper takes a Worker implementation and wraps it within a goroutine
// and channel arrangement. The workerWrapper is responsible for managing the
// lifetime of both the Worker and the goroutine.
type workerWrapper struct {
	worker       Worker
	idleTime     time.Duration
	uuid         string
	idleWorkChan chan<- string

	interruptChan chan struct{}

	// reqChan is NOT owned by this type, it is used to send requests for work.
	reqChan chan<- workRequest

	// closeChan can be closed in order to cleanly shutdown this worker.
	closeChan chan struct{}

	// closedChan is closed by the run() goroutine when it exits.
	closedChan chan struct{}
}

func newWorkerWrapper(
	reqChan chan<- workRequest,
	worker Worker,
	idleTime time.Duration,
	uuid string,
	idleWorkChan chan<- string,
) *workerWrapper {
	w := workerWrapper{
		worker:        worker,
		idleTime:      idleTime,
		uuid:          uuid,
		idleWorkChan:  idleWorkChan,
		interruptChan: make(chan struct{}),
		reqChan:       reqChan,
		closeChan:     make(chan struct{}),
		closedChan:    make(chan struct{}),
	}

	go w.run()

	return &w
}

//------------------------------------------------------------------------------

func (w *workerWrapper) interrupt() {
	close(w.interruptChan)
	w.worker.Interrupt()
}

func (w *workerWrapper) run() {
	jobChan, retChan := make(chan jobData), make(chan interface{})
	defer func() {
		w.worker.Terminate()
		close(retChan)
		close(w.closedChan)
	}()

	idleTimeTimer := time.NewTimer(w.idleTime)
	defer idleTimeTimer.Stop()
	for {
		// NOTE: Blocking here will prevent the worker from closing down.
		w.worker.BlockUntilReady()
		// timer may be not active and may not fired
		if !idleTimeTimer.Stop() {
			select {
			case <-idleTimeTimer.C: //drain from the channel
			default:
			}
		}
		idleTimeTimer.Reset(w.idleTime)
		select {
		case w.reqChan <- workRequest{
			jobChan:       jobChan,
			retChan:       retChan,
			interruptFunc: w.interrupt,
		}:
			select {
			case jobInData := <-jobChan:
				result := w.worker.Process(jobInData.payload)
				//result must not be nil. otherelse it will block select
				if result == nil {
					result = true
				}

				if jobInData.asyncProcess {
					if jobInData.asyncJobComplete != nil {
						jobInData.asyncJobComplete(result)
					}
				} else {
					select {
					case retChan <- result:
					case <-w.interruptChan:
						w.interruptChan = make(chan struct{})
					}
				}
			case <-w.interruptChan:
				w.interruptChan = make(chan struct{})
			}
		case <-idleTimeTimer.C:
			w.idleWorkChan <- w.uuid
		case <-w.closeChan:
			return
		}
	}
}

//------------------------------------------------------------------------------

func (w *workerWrapper) stop() {
	close(w.closeChan)
}

func (w *workerWrapper) join() {
	<-w.closedChan
}

//------------------------------------------------------------------------------
