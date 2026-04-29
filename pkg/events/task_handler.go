package events

import (
	"fmt"

	"github.com/Daniil-1622/todo-app/pkg/kafka"
	"github.com/sirupsen/logrus"
)

type TaskEventHandler struct{}

func NewTaskEventHandler() *TaskEventHandler {
	return &TaskEventHandler{}
}

func (h *TaskEventHandler) Handle(event kafka.Event) error {
	logrus.Infof("[TaskHandler] Processing event: %s, task_id=%d", event.EventType, event.TaskId)

	switch event.EventType {
	case kafka.EventTaskCreated:
		return h.handleTaskCreated(event)
	case kafka.EventTaskUpdated:
		return h.handleTaskUpdated(event)
	case kafka.EventTaskDeleted:
		return h.handleTaskDeleted(event)
	case kafka.EventTaskCompleted:
		return h.handleTaskCompleted(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (h *TaskEventHandler) handleTaskCreated(event kafka.Event) error {
	logrus.Infof("[TaskHandler] Task created: task_id=%d, user_id=%d, data=%v",
		event.TaskId, event.UserId, event.Data)
	return nil
}

func (h *TaskEventHandler) handleTaskUpdated(event kafka.Event) error {
	logrus.Infof("[TaskHandler] Task updated: task_id=%d, user_id=%d, data=%v",
		event.TaskId, event.UserId, event.Data)
	return nil
}

func (h *TaskEventHandler) handleTaskDeleted(event kafka.Event) error {
	logrus.Infof("[TaskHandler] Task deleted: task_id=%d, user_id=%d",
		event.TaskId, event.UserId)
	return nil
}

func (h *TaskEventHandler) handleTaskCompleted(event kafka.Event) error {
	logrus.Infof("[TaskHandler] Task completed: task_id=%d, user_id=%d",
		event.TaskId, event.UserId)
	return nil
}
