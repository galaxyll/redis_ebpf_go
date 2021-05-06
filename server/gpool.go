package main

type Task struct {
	f func() error
}

func NewTask(f func() error) *Task {
	t := Task{
		f: f,
	}

	return &t
}

func (t *Task) Execute() {
	t.f()
}

type Pool struct {
	EntryChannel chan *Task

	worker_num int

	JobsChannel chan *Task
}

func NewPool(cap int) *Pool {
	p := Pool{
		EntryChannel: make(chan *Task),
		worker_num:   cap,
		JobsChannel:  make(chan *Task),
	}

	return &p
}

func (p *Pool) worker(work_ID int) {
	for task := range p.JobsChannel {
		task.Execute()
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.worker_num; i++ {
		go p.worker(i)
	}

	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}

	close(p.JobsChannel)

	close(p.EntryChannel)
}
