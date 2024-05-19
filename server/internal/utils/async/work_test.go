package async

import (
	"sync"
	"testing"
	"time"
)

func TestNewWork(t *testing.T) {

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			go func(j int) {
				time.Sleep(time.Duration(j) * time.Second)
				println(j)
				if j == 2 {
					wait.Done()
				}
			}(i)
		}

	}()

	wait.Wait()
	time.Sleep(5 * time.Second)
	println("done")

}
