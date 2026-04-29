package handler

import (
	"net/http"
	"time"

	todo "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Это обработчики (handlers) для маршрутов списков.

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: kafka.EventUserRegistered,
			TimeStamp: time.Now().UTC(),
			UserId:    id,
			Data: map[string]interface{}{
				"username": input.Username,
				"name":     input.Name,
			},
		}
		if err := h.producer.Publish(kafka.TopicUserEvents, event); err != nil {
			logrus.Errorf("[Handler] Failed to publish USER_REGISTERED: %s", err.Error())
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type singInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input singInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: kafka.EventUserLoggedIn,
			TimeStamp: time.Now().UTC(),
			Data: map[string]interface{}{
				"username": input.Username,
			},
		}
		if err := h.producer.Publish(kafka.TopicUserEvents, event); err != nil {
			logrus.Errorf("[Handler] Failed to publish USER_LOGGED_IN: %s", err.Error())
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}
