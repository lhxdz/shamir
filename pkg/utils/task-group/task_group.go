package taskgroup

import (
	"sync"

	"shamir/pkg/utils/log"
)

type TaskGroup struct {
	taskMap *sync.Map
}

func NewTaskGroup() *TaskGroup {
	return &TaskGroup{
		taskMap: &sync.Map{},
	}
}

func (tg *TaskGroup) Add(t Task) bool {
	newT := newTask(t)

	_, ok := tg.taskMap.LoadOrStore(newT.Name(), newT)
	if ok {
		// 已经存在，则加入失败
		log.Warnf("add task%q failed, exist", t.Name())
		return false
	}

	log.Infof("added task%q", t.Name())
	return true
}

func (tg *TaskGroup) Start(name string) bool {
	t, ok := tg.taskMap.Load(name)
	if !ok {
		// 不存在，则启动失败
		log.Errorf("starting task%q failed, not exist", name)
		return false
	}

	log.Infof("starting task%q", name)
	go t.(*task).Start()
	return true
}

func (tg *TaskGroup) Stop(name string) bool {
	t, ok := tg.taskMap.Load(name)
	if !ok {
		// 不存在，则停止失败
		log.Errorf("stopping task%q failed, not exist", name)
		return false
	}

	log.Infof("stopping task%q", name)
	t.(*task).Stop()
	return true
}

func (tg *TaskGroup) Delete(name string) bool {
	t, ok := tg.taskMap.Load(name)
	if !ok {
		// 不存在，则不用删除，表示成功
		log.Warnf("task%q not exist", name)
		return true
	}

	if t.(*task).Status() == Running {
		log.Errorf("task%q delete failed: task is running", name)
		return false
	}

	tg.taskMap.Delete(name)
	log.Infof("deleted task%q", name)
	return true
}

func (tg *TaskGroup) StartAll() {
	log.Info("starting all task...")
	tg.taskMap.Range(func(key, value any) bool {
		t := value.(*task)
		if t.Status() == Running {
			log.Warnf("task%q is running", t.Name())
			return true
		}

		log.Infof("starting task%q", t.Name())
		go t.Start()
		return true
	})
}

func (tg *TaskGroup) StopAll() {
	log.Info("stopping all task...")
	tg.taskMap.Range(func(key, value any) bool {
		t := value.(*task)
		if t.Status() != Running {
			log.Warnf("task%q is not running", t.Name())
			return true
		}

		log.Infof("stopping task%q", t.Name())
		t.Stop()
		return true
	})
}

// DeleteAll 将会删除所有已经停止的任务
func (tg *TaskGroup) DeleteAll() {
	log.Info("deleting all task...")
	tg.taskMap.Range(func(key, value any) bool {
		if value.(*task).Status() == Running {
			log.Warnf("task%q delete failed: task is running", key)
			return true
		}

		tg.taskMap.Delete(key)
		log.Infof("deleted task %q", key)
		return true
	})
}

func (tg *TaskGroup) IsExist(name string) bool {
	_, ok := tg.taskMap.Load(name)
	return ok
}

func (tg *TaskGroup) GetStatus(name string) (Status, bool) {
	t, ok := tg.taskMap.Load(name)
	if !ok {
		return 0, false
	}

	return t.(*task).Status(), true
}

func (tg *TaskGroup) TaskNames() []string {
	tasks := make([]string, 0)
	tg.taskMap.Range(func(key, value any) bool {
		name, ok := key.(string)
		if ok {
			tasks = append(tasks, name)
		}
		return true
	})
	return tasks
}
