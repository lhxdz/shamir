package taskgroup

import (
	"sync/atomic"

	"shamir/pkg/utils/log"
)

type Status int32

const (
	Ready   Status = 0
	Running Status = 1
	Stopped Status = 2
)

type Task interface {
	Name() string
	Start() error
	Stop()
}

type task struct {
	Task
	status atomic.Int32
	err    error
}

func newTask(t Task) *task {
	newT := &task{
		Task:   t,
		status: atomic.Int32{},
	}
	newT.status.Store(int32(Ready))

	return newT
}

func (t *task) Start() {
	if Status(t.status.Load()) == Running {
		log.Warnf("task(%q) is running", t.Name())
		return
	}
	t.err = t.Task.Start()
	if t.Error() != nil {
		log.Errorf("task%q start failed: %v", t.Name(), t.Error())
	}
	t.status.Store(int32(Stopped))
	return
}

func (t *task) Stop() {
	if Status(t.status.Load()) != Running {
		log.Warnf("task(%q) is not running", t.Name())
		return
	}

	t.Task.Stop()
}

func (t *task) Status() Status {
	return Status(t.status.Load())
}

func (t *task) Error() error {
	return t.err
}
