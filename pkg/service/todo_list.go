package service

import (
	"time"

	todo "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/kafka"
	"github.com/Daniil-1622/todo-app/pkg/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TodoListService struct {
	repo     repository.TodoList
	producer kafka.Producer
}

func NewTodoListService(repo repository.TodoList, producer kafka.Producer) *TodoListService {
	return &TodoListService{repo: repo, producer: producer}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	id, err := s.repo.Create(userId, list)
	if err != nil {
		return 0, err
	}

	go func() {
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: kafka.EventListCreated,
			TimeStamp: time.Now().UTC(),
			UserId:    userId,
			ListId:    id,
			Data: map[string]interface{}{
				"id":          list.Id,
				"title":       list.Title,
				"description": list.Description,
			},
		}
		if err := s.producer.Publish(kafka.TopicListEvents, event); err != nil {
			logrus.Errorf("[TodoListService] Failed to publish LIST_CREATED: %s", err.Error())
		}
	}()
	return id, nil
}

func (s *TodoListService) GetAll(userid int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userid)
}

func (s *TodoListService) GetById(userid, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userid, listId)
}

func (s *TodoListService) Delete(userid, listId int) error {
	err := s.repo.Delete(userid, listId)
	if err != nil {
		return err
	}

	go func() {
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: kafka.EventListDeleted,
			TimeStamp: time.Now().UTC(),
			UserId:    userid,
			ListId:    listId,
		}
		if err := s.producer.Publish(kafka.TopicListEvents, event); err != nil {
			logrus.Errorf("[TodoListService] Failed to publish LIST_DELETED: %s", err.Error())
		}
	}()

	return nil
}

func (s *TodoListService) Update(userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	err := s.repo.Update(userId, listId, input)
	if err != nil {
		return err
	}
	go func() {
		eventType := kafka.EventListUpdated
		event := kafka.Event{
			EventId:   uuid.New().String(),
			EventType: eventType,
			TimeStamp: time.Now().UTC(),
			UserId:    userId,
			ListId:    listId,
			Data: map[string]interface{}{
				"title":       input.Title,
				"description": input.Description,
			},
		}
		if err := s.producer.Publish(kafka.TopicListEvents, event); err != nil {
			logrus.Errorf("[TodoListService] Failed to publish %s: %s", eventType, err.Error())
		}
	}()

	return nil
}
