package events

import (
	"fmt"

	"github.com/Daniil-1622/todo-app/pkg/kafka"
	"github.com/sirupsen/logrus"
)

type ListEventHandler struct{}

func NewListEventHandler() *ListEventHandler {
	return &ListEventHandler{}
}

func (h *ListEventHandler) Handle(event kafka.Event) error {
	logrus.Infof("[ListHandler] Processing event: %s, list_id=%d", event.EventType, event.ListId)

	switch event.EventType {
	case kafka.EventListCreated:
		return h.handleListCreated(event)
	case kafka.EventListUpdated:
		return h.handleListUpdated(event)
	case kafka.EventListDeleted:
		return h.handleListDeleted(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (h *ListEventHandler) handleListCreated(event kafka.Event) error {
	logrus.Infof("[ListHandler] List created: list_id=%d, user_id=%d, data=%v",
		event.ListId, event.UserId, event.Data)
	return nil
}

func (h *ListEventHandler) handleListUpdated(event kafka.Event) error {
	logrus.Infof("[ListHandler] List updated: list_id=%d, user_id=%d, data=%v",
		event.ListId, event.UserId, event.Data)
	return nil
}

func (h *ListEventHandler) handleListDeleted(event kafka.Event) error {
	logrus.Infof("[ListHandler] List deleted: list_id=%d, user_id=%d",
		event.ListId, event.UserId)
	return nil
}
