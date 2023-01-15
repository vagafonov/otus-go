package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task, 100)
	for _, val := range tasks {
		ch <- val
	}
	close(ch)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	errorsCount := 0
	wg.Add(n)
	var er error

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range ch {
				taskResult := task()
				if taskResult != nil {
					mu.Lock()
					errorsCount++
					errorsCountSafe := errorsCount
					mu.Unlock()
					if m > 0 {
						if errorsCountSafe >= m {
							mu.Lock()
							er = ErrErrorsLimitExceeded
							mu.Unlock()
							return
						}
					}
				}
			}
		}()
	}

	wg.Wait()
	return er
}
