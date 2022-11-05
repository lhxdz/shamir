package taskgroup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTask struct {
	name   string
	finish chan struct{}
}

func newTestTask(name string) *testTask {
	return &testTask{
		name:   name,
		finish: make(chan struct{}, 1),
	}
}

func (t *testTask) Name() string {
	return t.name
}

func (t *testTask) Start() error {
	<-t.finish
	return nil
}

func (t *testTask) Stop() {
	t.finish <- struct{}{}
}

func TestTaskGroup(t *testing.T) {
	tg := NewTaskGroup()

	taskNames := []string{"task1", "task2", "task3"}
	for _, name := range taskNames {
		ok := tg.Add(newTestTask(name))
		assert.True(t, ok)
	}

	tg.StartAll()

	taskName := "task_latest"
	ok := tg.Add(newTestTask(taskName))
	assert.True(t, ok)
	assert.True(t, tg.IsExist(taskName))
	status, ok := tg.GetStatus(taskName)
	assert.Equal(t, status, Ready)
	assert.True(t, ok)
	assert.Contains(t, tg.TaskNames(), taskName)

	notExistName := "notExistName"
	ok = tg.Start(taskName)
	assert.True(t, ok)
	ok = tg.Start(notExistName)
	assert.False(t, ok)

	ok = tg.Stop(taskName)
	assert.True(t, ok)
	ok = tg.Stop(notExistName)
	assert.False(t, ok)

	tg.StopAll()

	ok = tg.Delete(taskName)
	assert.True(t, ok)
	ok = tg.Delete(notExistName)
	assert.True(t, ok)

	tg.DeleteAll()
}
