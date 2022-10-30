package taskgroup

import (
	"fmt"
	"sync"
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
		return false
	}

	return true
}

func (tg *TaskGroup) Start(name string) bool {
	fmt.Printf("starting task(%q)...\n", name)
	t, ok := tg.taskMap.Load(name)
	if !ok {
		// 不存在，则启动失败
		fmt.Printf("starting task(%q) failed, not exist\n", name)
		return false
	}

	go t.(*task).Start()

	return true
}

func (tg *TaskGroup) Stop(name string) bool {
	fmt.Printf("stopping task(%q)...\n", name)
	t, ok := tg.taskMap.Load(name)
	if !ok {
		// 不存在，则启动失败
		fmt.Printf("stopping task(%q) failed, not exist\n", name)
		return false
	}

	t.(*task).Stop()
	return true
}

func (tg *TaskGroup) Delete(name string) bool {
	fmt.Printf("delete task(%q)...\n", name)
	t, ok := tg.taskMap.Load(name)
	if !ok {
		// 不存在，则不用删除，表示成功
		fmt.Printf("task(%q) not exist\n", name)
		return true
	}

	if t.(*task).Status() == Running {
		fmt.Printf("task(%q) delete failed: task is running\n", name)
		return false
	}

	tg.taskMap.Delete(name)
	return true
}

func (tg *TaskGroup) StartAll() {
	fmt.Println("starting all task...")
	tg.taskMap.Range(func(key, value any) bool {
		t := value.(*task)
		if t.Status() == Running {
			fmt.Printf("task(%q) is running\n", t.Name())
			return true
		}

		fmt.Printf("starting task(%q)...\n", t.Name())
		go t.Start()
		return true
	})
}

func (tg *TaskGroup) StopAll() {
	fmt.Println("stopping all task...")
	tg.taskMap.Range(func(key, value any) bool {
		t := value.(*task)
		if t.Status() != Running {
			fmt.Printf("task(%q) is not running\n", t.Name())
			return true
		}

		fmt.Printf("stopping task(%q)...\n", t.Name())
		t.Stop()
		return true
	})
}

func (tg *TaskGroup) DeleteAll() {
	fmt.Println("deleting all task...")
	tg.taskMap.Range(func(key, value any) bool {
		if value.(*task).Status() == Running {
			fmt.Printf("task(%q) delete failed: task is running\n", key)
			return true
		}

		tg.taskMap.Delete(key)
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
