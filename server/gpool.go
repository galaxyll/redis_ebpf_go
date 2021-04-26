package main

type Task struct {
	f func() error
}

func NewTask(f func() error) *Task {
	return &Task{f}
}

func (t *Task) Execute() {
	t.f()
}

type Pool struct {
	EntryTaskPool chan *Task
	worker_num    int
	JobsChannel   chan *Task
}

func NewPool(cap int) *Pool {
	return &Pool{
		EntryTaskPool: make(chan *Task),
		worker_num:    cap,
		JobsChannel:   make(chan *Task),
	}
}
