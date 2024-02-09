package main

import (
	"sync"
	"time"
)

type Job func()

type Pool struct {
	workQueue chan Job
	wg        sync.WaitGroup
}

func NewPool(limit int) *Pool {
	pool := &Pool{
		workQueue: make(chan Job),
	}

	pool.wg.Add(limit)

	for i := 0; i < limit; i++ {
		go func() {
			defer pool.wg.Done()

			for job := range pool.workQueue {
				job()
			}
		}()
	}

	return pool
}

func (pool *Pool) AddJob(job Job) {
	pool.workQueue <- job
}

func (pool *Pool) Wait() {
	close(pool.workQueue)
}

func main() {
	pool := NewPool(2)

	for i := 0; i < 30; i++ {
		job := func() {
			time.Sleep(1 * time.Second)
			println("Job completed", i, "\n")
		}

		pool.AddJob(job)
	}

	pool.Wait()
}
