package comm

import (
	"log"
	"sync"
	"time"
)

type ITimerWork interface {
	Start()
	RegistTask(task *Task)
	Remove(taskid string)
}

type TimerWork struct {
	TaskList map[string]*Task
	mutex    sync.RWMutex
	wg       sync.WaitGroup
}

func NewTimerWork() *TimerWork {
	return &TimerWork{
		TaskList: make(map[string]*Task),
	}
}

/*
   启动注册的Task定时任务
*/
func (t *TimerWork) Start() {
	for _, v := range t.TaskList {
		t.wg.Add(1)
		/*使用range 时 创建goroutine时 参数要用值拷贝传递到函数内部 不然参数会一直时最后一个*/
		go func(task *Task) {
			t.wg.Done()
			//task.Stopflag = 1
			for {
				if task.Stopflag == 0 { //结束任务
					log.Printf("TimerWork OnStopTask:%v Interval:%ds", v.TaskId, v.Interval)
					return
				}
				task.Work(task.Param)
				time.Sleep(task.Interval * time.Second)
			}
		}(v)
		log.Printf("TimerWork Start Task:%v Interval:%ds", v.TaskId, v.Interval)
	}
	t.wg.Wait()
}

/*
   注册定时任务
*/
func (t *TimerWork) RegistTask(task *Task) {
	t.mutex.Lock()
	if _, exists := t.TaskList[task.TaskId]; !exists {
		t.TaskList[task.TaskId] = task
		t.TaskList[task.TaskId].Stopflag = 1 //开启
	}
	t.mutex.Unlock()
}

/*
   根据TaskId  移除定时任务
*/
func (t *TimerWork) Remove(taskid string) {
	t.mutex.Lock()
	if _, exists := t.TaskList[taskid]; exists {
		t.TaskList[taskid].Stopflag = 0 //关闭
		delete(t.TaskList, taskid)
	}
	t.mutex.Unlock()
}

//移除TimerWork 里面的所有task
func (t *TimerWork) removeAll() {
	t.mutex.Lock()
	for taskid, _ := range t.TaskList {
		t.TaskList[taskid].Stopflag = 0 //关闭
		delete(t.TaskList, taskid)
	}
	t.mutex.Unlock()
}

//Stop  停止 TimerWork
func (t *TimerWork) Stop() {
	t.removeAll()
}

/*
   定时任务Task
*/
type Task struct {
	Interval time.Duration
	TaskId   string
	Work     func(interface{})
	Param    interface{}
	Stopflag int //停止标志  1:开启 0:关闭
}
