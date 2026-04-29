package service

import (
	"time"

	"github.com/google/uuid"

	todo "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/kafka"
	"github.com/Daniil-1622/todo-app/pkg/repository"
	"github.com/sirupsen/logrus"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
	producer kafka.Producer
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList, producer kafka.Producer) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo, producer: producer}
}

func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(listId, item)
	if err != nil {
		return 0, err
	}

	go func() {
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: kafka.EventTaskCreated,
			TimeStamp: time.Now().UTC(),
			UserId:    userId,
			ListId:    listId,
			TaskId:    id,
			Data: map[string]interface{}{
				"title":       item.Title,
				"description": item.Description,
				"done":        item.Done,
			},
		}
		if err := s.producer.Publish(kafka.TopicTaskEvents, event); err != nil {
			logrus.Errorf("[TodoItemService] Failed to publish TASK_CREATED: %s", err.Error())
		}
	}()

	return id, nil
}

func (s *TodoItemService) GetAll(userid, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userid, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	err := s.repo.Delete(userId, itemId)
	if err != nil {
		return err
	}

	go func() {
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: kafka.EventTaskDeleted,
			TimeStamp: time.Now().UTC(),
			UserId:    userId,
			TaskId:    itemId,
		}
		if err := s.producer.Publish(kafka.TopicTaskEvents, event); err != nil {
			logrus.Errorf("[TodoItemService] Failed to publish TASK_DELETED: %s", err.Error())
		}
	}()

	return nil
}

func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	err := s.repo.Update(userId, itemId, input)
	if err != nil {
		return err
	}

	go func() {
		eventType := kafka.EventTaskUpdated
		if input.Done != nil && *input.Done {
			eventType = kafka.EventTaskCompleted
		}
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: eventType,
			TimeStamp: time.Now().UTC(),
			UserId:    userId,
			TaskId:    itemId,
			Data: map[string]interface{}{
				"title":       input.Title,
				"description": input.Description,
				"done":        input.Done,
			},
		}
		if err := s.producer.Publish(kafka.TopicTaskEvents, event); err != nil {
			logrus.Errorf("[TodoItemService] Failed to publish %s: %s", eventType, err.Error())
		}
	}()

	return nil
}
