/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package grpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/commcos/utils/uuid"
)

//------------------------------------------------------------------------------

// Errors that are used throughout the Tunny API.
var (
	ErrPoolNotRunning = errors.New("the pool is not running")
	ErrJobNotFunc     = errors.New("generic worker not given a func()")
	ErrWorkerClosed   = errors.New("worker was closed")
	ErrJobTimedOut    = errors.New("job request timed out")
	defaultIdleTime   = 10 * time.Minute
)

// Worker is an interface representing a Tunny working agent. It will be used to
// block a calling goroutine until ready to process a job, process that job
// synchronously, interrupt its own process call when jobs are abandoned, and
// clean up its resources when being removed from the pool.
//
// Each of these duties are implemented as a single method and can be averted
// when not needed by simply implementing an empty func.
type Worker interface {
	// Process will synchronously perform a job and return the result.
	Process(interface{}) interface{}

	// BlockUntilReady is called before each job is processed and must block the
	// calling goroutine until the Worker is ready to process the next job.
	BlockUntilReady()

	// Interrupt is called when a job is cancelled. The worker is responsible
	// for unblocking the Process implementation.
	Interrupt()

	// Terminate is called when a Worker is removed from the processing pool
	// and is responsible for cleaning up any held resources.
	Terminate()
}

//------------------------------------------------------------------------------

// closureWorker is a minimal Worker implementation that simply wraps a
// func(interface{}) interface{}
type closureWorker struct {
	processor func(interface{}) interface{}
}

func (w *closureWorker) Process(payload interface{}) interface{} {
	return w.processor(payload)
}

func (w *closureWorker) BlockUntilReady() {}
func (w *closureWorker) Interrupt()       {}
func (w *closureWorker) Terminate()       {}

//------------------------------------------------------------------------------

// callbackWorker is a minimal Worker implementation that attempts to cast
// each job into func() and either calls it if successful or returns
// ErrJobNotFunc.
type callbackWorker struct{}

func (w *callbackWorker) Process(payload interface{}) interface{} {
	f, ok := payload.(func())
	if !ok {
		return ErrJobNotFunc
	}
	f()
	return nil
}

func (w *callbackWorker) BlockUntilReady() {}
func (w *callbackWorker) Interrupt()       {}
func (w *callbackWorker) Terminate()       {}

//------------------------------------------------------------------------------

// Pool is a struct that manages a collection of workers, each with their own
// goroutine. The Pool can initialize, expand, compress and close the workers,
// as well as processing jobs with the workers synchronously.
type Pool struct {
	queuedJobs int64

	ctor         func() Worker
	workers      map[string]*workerWrapper
	reqChan      chan workRequest
	initSize     int
	idleTime     time.Duration
	idleWorkChan chan string
	workerMut    sync.Mutex
}

// New creates a new Pool of workers that starts with n workers. You must
// provide a constructor function that creates new Worker types and when you
// change the size of the pool the constructor will be called to create each new
// Worker.
func New(n int, ctor func() Worker) *Pool {
	p := &Pool{
		ctor:         ctor,
		reqChan:      make(chan workRequest),
		initSize:     n,
		idleTime:     defaultIdleTime,
		idleWorkChan: make(chan string),
	}
	p.SetSize(n)
	go p.run()
	return p
}

// NewWithTimeout creates a new Pool of workers that starts with n workers. You must
// provide a constructor function that creates new Worker types and when you
// change the size of the pool the constructor will be called to create each new
// Worker.
func NewWithTimeout(n int, ctor func() Worker, idleTime time.Duration) *Pool {
	p := &Pool{
		ctor:         ctor,
		reqChan:      make(chan workRequest),
		initSize:     n,
		idleTime:     idleTime,
		idleWorkChan: make(chan string),
	}
	p.SetSize(n)
	go p.run()
	return p
}

// NewFunc creates a new Pool of workers where each worker will process using
// the provided func.
func NewFunc(n int, f func(interface{}) interface{}) *Pool {
	return New(n, func() Worker {
		return &closureWorker{
			processor: f,
		}
	})
}

// NewCallback creates a new Pool of workers where workers cast the job payload
// into a func() and runs it, or returns ErrNotFunc if the cast failed.
func NewCallback(n int) *Pool {
	return New(n, func() Worker {
		return &callbackWorker{}
	})
}

//------------------------------------------------------------------------------

func (p *Pool) run() {
	for {
		uuid := <-p.idleWorkChan
		p.removeWorkerByUUID(uuid)
	}
}

// remove special worker
func (p *Pool) removeWorkerByUUID(uuid string) {
	currSize := p.GetSize()
	if currSize > p.initSize {
		if currWorker, ok := p.workers[uuid]; ok {
			delete(p.workers, uuid)
			currWorker.stop()
			// Synchronously wait for worker to stop
			currWorker.join()
		}
	}
}

// Process will use the Pool to process a payload and synchronously return the
// result. Process can be called safely by any goroutines, but will panic if the
// Pool has been stopped.
func (p *Pool) Process(payload interface{}) interface{} {
	atomic.AddInt64(&p.queuedJobs, 1)
	defer atomic.AddInt64(&p.queuedJobs, -1)

	request, open := <-p.reqChan
	if !open {
		panic(ErrPoolNotRunning)
	}

	request.jobChan <- jobData{
		payload:      payload,
		asyncProcess: false,
	}

	payload, open = <-request.retChan
	if !open {
		panic(ErrWorkerClosed)
	}

	return payload
}

// AsyncProcess 异步处理process，不等待结果
func (p *Pool) AsyncProcess(
	payload interface{},
) error {
	atomic.AddInt64(&p.queuedJobs, 1)

	request, open := <-p.reqChan
	if !open {
		panic(ErrPoolNotRunning)
	}

	request.jobChan <- jobData{
		payload:      payload,
		asyncProcess: true,
		asyncJobComplete: func(result interface{}) {
			defer atomic.AddInt64(&p.queuedJobs, -1)
		},
	}

	return nil
}

// ProcessTimed will use the Pool to process a payload and synchronously return
// the result. If the timeout occurs before the job has finished the worker will
// be interrupted and ErrJobTimedOut will be returned. ProcessTimed can be
// called safely by any goroutines.
func (p *Pool) ProcessTimed(
	payload interface{},
	timeout time.Duration,
) (interface{}, error) {
	atomic.AddInt64(&p.queuedJobs, 1)
	defer atomic.AddInt64(&p.queuedJobs, -1)

	tout := time.NewTimer(timeout)

	var request workRequest
	var open bool

	select {
	case request, open = <-p.reqChan:
		if !open {
			return nil, ErrPoolNotRunning
		}
	case <-tout.C:
		return nil, ErrJobTimedOut
	}

	select {
	case request.jobChan <- jobData{
		payload:      payload,
		asyncProcess: false,
	}:
	case <-tout.C:
		request.interruptFunc()
		return nil, ErrJobTimedOut
	}

	select {
	case payload, open = <-request.retChan:
		if !open {
			return nil, ErrWorkerClosed
		}
	case <-tout.C:
		request.interruptFunc()
		return nil, ErrJobTimedOut
	}

	tout.Stop()
	return payload, nil
}

// QueueLength returns the current count of pending queued jobs.
func (p *Pool) QueueLength() int64 {
	return atomic.LoadInt64(&p.queuedJobs)
}

// SetSize changes the total number of workers in the Pool. This can be called
// by any goroutine at any time unless the Pool has been stopped, in which case
// a panic will occur.
func (p *Pool) SetSize(n int) {
	p.workerMut.Lock()
	defer p.workerMut.Unlock()

	lWorkers := len(p.workers)
	if lWorkers == n {
		return
	}
	if p.workers == nil {
		p.workers = make(map[string]*workerWrapper)
	}

	// Add extra workers if N > len(workers)
	for i := lWorkers; i < n; i++ {
		uuidStr := string(uuid.NewUUID())
		p.workers[uuidStr] = newWorkerWrapper(p.reqChan, p.ctor(), p.idleTime, uuidStr, p.idleWorkChan)
	}

	if n < lWorkers {
		deleteWorkerCount := lWorkers - n
		var deleteWorkers []*workerWrapper
		for uuid, worker := range p.workers {
			deleteWorkers = append(deleteWorkers, worker)
			defer delete(p.workers, uuid)
			deleteWorkerCount--
			if deleteWorkerCount < 0 {
				break
			}
		}
		for _, worker := range deleteWorkers {
			worker.stop()
		}

		for _, worker := range deleteWorkers {
			worker.join()
		}
	}
}

// GetSize returns the current size of the pool.
func (p *Pool) GetSize() int {
	p.workerMut.Lock()
	defer p.workerMut.Unlock()

	return len(p.workers)
}

// Close will terminate all workers and close the job channel of this Pool.
func (p *Pool) Close() {
	p.SetSize(0)
	close(p.reqChan)
}

//------------------------------------------------------------------------------
