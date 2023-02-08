package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return errors.New("invalid worker number")
	}

	wg := sync.WaitGroup{}
	var errorsCount int32
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
				if m > 0 {
					if atomic.LoadInt32(&errorsCount) >= int32(m) {
						er = ErrErrorsLimitExceeded
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
					atomic.AddInt32(&errorsCount, 1)
				}
			}
			fmt.Println("Gorutine", i, "exit. Empty tasks")
		}(i)
	}

	wg.Wait()
	return er
}
