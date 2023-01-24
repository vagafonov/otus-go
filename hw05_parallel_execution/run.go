package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return errors.New("invalid worker number")
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	errorsCount := 0
	wg.Add(n)
	var er error

	ch := make(chan Task)

	go func() {
		defer close(ch)
		for _, val := range tasks {
			ch <- val
		}
	}()

	for i := 0; i < n; i++ {
		go func(i int) {
			fmt.Println("Gorutine", i, "start")

			defer wg.Done()
			for task := range ch {
				mu.Lock()
				if m > 0 {
					errorsCountSafe := errorsCount
					if errorsCountSafe >= m {
						er = ErrErrorsLimitExceeded
						mu.Unlock()
						fmt.Println("Gorutine", i, "exit")
						_, ok := <-ch
						if !ok {
							close(ch)
						}
						return
					}
				}

				fmt.Println("Gorutine", i, "handle task")
				taskResult := task()
				if taskResult != nil {
					errorsCount++
				}
				mu.Unlock()
			}
			fmt.Println("Gorutine", i, "exit. Empty tasks")
		}(i)
	}

	wg.Wait()
	return er
}
