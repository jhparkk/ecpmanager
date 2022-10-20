package worker

import (
	"fmt"
	"testing"
	"time"
)

type TestWorker struct {
	worker *Worker
	id     int
	loop   int
}

func NewTestWorker(id int) *TestWorker {
	w := &TestWorker{}
	w.worker = NewWorker("testworker", w)
	w.id = id
	return w
}

func (w *TestWorker) In() (int, error) {
	fmt.Printf("TestWorker[%d] Start\n", w.id)
	w.loop = 0
	return 0, nil
}

func (w *TestWorker) Run() (int, error) {
	w.loop++
	fmt.Printf("TestWorker[%d] Run[%d]\n", w.id, w.loop)
	if w.loop > w.id {
		return 1, nil
	}
	w.worker.Sleep(1 * time.Second)
	return 0, nil
}

func (w *TestWorker) Out() (int, error) {
	fmt.Printf("TestWorker[%d] Out\n", w.id)
	return 0, nil
}

func TestRunWorker(t *testing.T) {
	for i := 1; i < 10; i++ {
		tworker := NewTestWorker(i)
		tworker.worker.Start()
	}

	time.Sleep(time.Second)
	fmt.Println("StopAllWorkers")
	go StopAllWorkers()
	WaitWorkers()
}
