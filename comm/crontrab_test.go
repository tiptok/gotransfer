package comm

import (
	"log"
	"testing"
	"time"
)

func TestTimerW(t *testing.T) {
	var itimer ITimerWork
	itimer = NewTimerWork()
	task1 := &Task{
		TaskId:   "task1",
		Interval: 10,
		Work: func(interface{}) {
			log.Printf("Current Task:%s OnWork", "task1")
		},
	}
	task2 := &Task{
		TaskId:   "task2",
		Interval: 20,
		Work: func(interface{}) {
			log.Printf("Currentyyyy Task:%s OnWork", "task2")
		},
	}
	itimer.RegistTask(task1)
	itimer.RegistTask(task2)
	itimer.Start()
	time.Sleep(time.Second * 200)
}
