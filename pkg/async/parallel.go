package async

import (
	"context"
	"sync"
)

// Result holds the result of an async operation
type Result struct {
	Value interface{}
	Error error
}

// Task represents an async task to execute
type Task func(ctx context.Context) (interface{}, error)

// RunParallel executes multiple tasks in parallel and waits for all to complete
// Returns results in the same order as the input tasks
func RunParallel(ctx context.Context, tasks ...Task) []Result {
	results := make([]Result, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t Task) {
			defer wg.Done()

			// Check context before starting
			select {
			case <-ctx.Done():
				results[idx] = Result{Error: ctx.Err()}
				return
			default:
			}

			value, err := t(ctx)
			results[idx] = Result{Value: value, Error: err}
		}(i, task)
	}

	wg.Wait()
	return results
}

// RunParallelWithLimit executes tasks in parallel with a concurrency limit
func RunParallelWithLimit(ctx context.Context, limit int, tasks ...Task) []Result {
	results := make([]Result, len(tasks))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, limit)

	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t Task) {
			defer wg.Done()

			// Acquire semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				results[idx] = Result{Error: ctx.Err()}
				return
			}

			value, err := t(ctx)
			results[idx] = Result{Value: value, Error: err}
		}(i, task)
	}

	wg.Wait()
	return results
}

// WorkerPool manages a pool of workers for batch operations
type WorkerPool struct {
	workers  int
	jobQueue chan func()
	done     chan struct{}
	wg       sync.WaitGroup
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		workers:  workers,
		jobQueue: make(chan func(), workers*2),
		done:     make(chan struct{}),
	}
	pool.start()
	return pool
}

// start initializes the workers
func (p *WorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		go func() {
			for {
				select {
				case job := <-p.jobQueue:
					job()
					p.wg.Done()
				case <-p.done:
					return
				}
			}
		}()
	}
}

// Submit adds a job to the pool
func (p *WorkerPool) Submit(job func()) {
	p.wg.Add(1)
	p.jobQueue <- job
}

// Wait waits for all submitted jobs to complete
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// Stop stops the worker pool
func (p *WorkerPool) Stop() {
	close(p.done)
}

// FanOut executes the same operation on multiple inputs in parallel
func FanOut[T any, R any](ctx context.Context, inputs []T, operation func(T) (R, error)) ([]R, []error) {
	results := make([]R, len(inputs))
	errors := make([]error, len(inputs))
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(idx int, in T) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errors[idx] = ctx.Err()
				return
			default:
			}

			result, err := operation(in)
			results[idx] = result
			errors[idx] = err
		}(i, input)
	}

	wg.Wait()
	return results, errors
}
