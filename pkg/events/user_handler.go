package events

import (
	"fmt"

	"github.com/Daniil-1622/todo-app/pkg/kafka"
	"github.com/sirupsen/logrus"
)

type UserEventHandler struct{}

func NewUserEventHandler() *UserEventHandler {
	return &UserEventHandler{}
}

func (h *UserEventHandler) Handle(event kafka.Event) error {
	logrus.Infof("[UserHandler] Processing event: %s, user_id=%d", event.EventType, event.UserId)

	switch event.EventType {
	case kafka.EventUserRegistered:
		return h.handleUserRegistered(event)
	case kafka.EventUserLoggedIn:
		return h.handleUserLoggedIn(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (h *UserEventHandler) handleUserRegistered(event kafka.Event) error {
	logrus.Infof("[UserHandler] User registered: user_id=%d, data=%v",
		event.UserId, event.Data)
	return nil
}

func (h *UserEventHandler) handleUserLoggedIn(event kafka.Event) error {
	logrus.Infof("[UserHandler] User logged in: user_id=%d",
		event.UserId)
	return nil
}
