package async

import (
	"testing"
)

func TestNewPool(t *testing.T) {

	newPool := NewWorkPool("test", 2, 5)
	for i := 0; i < 21; i++ {

		newPool.Submit(createFunc(i, t))
	}
	newPool.Close()
	newPool.Wait()

}

func createFunc(i int, t *testing.T) func() {
	cur := i
	return func() {
		t.Logf("current: %d", cur)
	}
}
