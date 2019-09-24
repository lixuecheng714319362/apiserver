package redismodel

import (
	"fmt"
)


//创建工作对象
type Score struct {
	Num int
}

//定义工作内容
func (s *Score) Do() {
	fmt.Println("num:", s.Num)
	//time.Sleep(10 * time.Millisecond)
}


// -------------------------工人--------------------------------
type Job interface {
	Do()
}

//工人维护一个工作队列
type Worker struct {
	JobList chan Job
}

//创建工人
func NewWorker() Worker {
	return Worker{JobList: make(chan Job)}
}

//工人执行工作 接收一个
func (w Worker) Run(wq chan chan Job) {
	go func() {
		for {
			//把自己工作列表放入工人队列中
			wq <- w.JobList

			//从工作队列中去除任务完成
			select {
			case job := <-w.JobList:
				job.Do()
			}
		}
	}()
}
// ---------------------工作池------------------------------------
type Pool struct {
	WorkerNum   int           //工人个数
	JobQueue    chan Job      //工作工作池中所有消息队列
	WorkerQueue chan chan Job //工人队列
}

func NewPool(workerNum, jobQueueCap int) *Pool {
	return &Pool{
		WorkerNum:   workerNum,
		JobQueue:    make(chan Job, jobQueueCap),
		WorkerQueue: make(chan chan Job, workerNum),
	}
}

//启动工作池
func (wp *Pool) Start() {
	fmt.Println("初始化worker")
	//初始化worker
	for i := 0; i < wp.WorkerNum; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue)
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				jobList := <-wp.WorkerQueue
				jobList <- job
			}
		}
	}()
}
// ---------------------------------------------------------