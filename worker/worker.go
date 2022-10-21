package worker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type WorkerProcedure interface {
	In() (int, error)
	Run() (int, error)
	Out() (int, error)
}

const (
	StatNone    uint8 = 1
	StatStart         = 2
	StatIn            = 3
	StatRun           = 4
	StatOut           = 5
	StatDone          = 6
	StatSuspend       = 7
)

const (
	CmdNone    uint8 = 0
	CmdStop          = 1
	CmdSuspend       = 2
	CmdResume        = 3
)

type Worker struct {
	work      WorkerProcedure
	wg        sync.WaitGroup
	wid       int
	name      string
	status    uint8
	cmd       uint8
	wakeUpCh  chan uint8
	isAliveCh bool
	chLock    sync.Mutex
}

// static variables
// private
var workerTable = make(map[int]*Worker)
var wtLock = &sync.Mutex{}
var widSeq = 0

// public
var Logger *log.Logger

func WaitWorkers() {
	var numWorkes int
	for {
		numWorkes = len(workerTable)
		if numWorkes > 0 {
			time.Sleep(100 * time.Millisecond)
		} else {
			return
		}
	}
}

func StopAllWorkers() {
	for _, w := range workerTable {
		w.Stop()
	}
}

func (w *Worker) Stop() {
	w.cmd = CmdStop
}

func NewWorker(name string, wp WorkerProcedure) *Worker {
	wtLock.Lock()
	widSeq++
	wtLock.Unlock()
	w := Worker{work: wp, wid: widSeq, name: name, status: StatNone}

	workerTable[widSeq] = &w
	return &w
}

func (w *Worker) Start() (int, error) {
	if w.status != StatNone {
		return -1, errors.New("worker already started")
	}

	w.status = StatStart
	w.wg.Add(1)
	go w.entry()

	return 0, nil
}

func (w *Worker) entry() {
	w.status = StatIn
	ret, err := w.work.In()
	if err != nil {
		if Logger != nil {
			Logger.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - In() failed : ", err)
		} else {
			fmt.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - In() failed : ", err)
		}
		w.wg.Done()
		return
	}
	if ret < 0 {
		if Logger != nil {
			Logger.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - In() failed : invalid ret:", ret)
		} else {
			fmt.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - In() failed : invalid ret:", ret)
		}
	}

	for {
		if w.cmd == CmdStop {
			break
		}
		if w.cmd == CmdSuspend {
			w.status = StatSuspend
			time.Sleep(time.Millisecond)
			continue
		}

		w.status = StatRun
		ret, err = w.work.Run()
		if err != nil {
			if Logger != nil {
				Logger.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - Run() failed : ", err)
			} else {
				fmt.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - Run() failed : ", err)
			}
			break
		}
		if ret != 0 {
			break
		}
	}

	w.status = StatOut
	ret, err = w.work.Out()
	if err != nil {
		if Logger != nil {
			Logger.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - Out() failed : ", err)
		} else {
			fmt.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - Out() failed : ", err)
		}
	}

	if ret < 0 {
		if Logger != nil {
			Logger.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - Out() failed : invalid ret:", ret)
		} else {
			fmt.Println("[", os.Getpid(), "]"+w.name+"_"+strconv.Itoa(w.wid), " - Out() failed : invalid ret:", ret)
		}
	}

	w.wg.Done()
	w.status = StatDone

	delete(workerTable, w.wid)
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func (w *Worker) SetCommand(cmd uint8) {
	w.cmd = cmd
}

func (w *Worker) Sleep(d time.Duration) {
	isSleepOn := true
	w.wakeUpCh = make(chan uint8)
	w.isAliveCh = true
	go func() {
		time.Sleep(d)
		w.chLock.Lock()
		if w.isAliveCh && isSleepOn {
			// send channel
			w.wakeUpCh <- 1
		}
		w.chLock.Unlock()
	}()
	// recv channel
	<-w.wakeUpCh

	w.chLock.Lock()
	w.isAliveCh = false
	close(w.wakeUpCh)
	w.chLock.Unlock()
	isSleepOn = false
}

func (w *Worker) WakeUp() {
	w.chLock.Lock()
	if w.isAliveCh {
		// send channel
		w.wakeUpCh <- 1
	}
	w.chLock.Unlock()
}

func (w *Worker) Wid() int {
	return w.wid
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) Status() uint8 {
	return w.status
}
