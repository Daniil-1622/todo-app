package service

import (
	todo "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userid int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userid)
}

func (s *TodoListService) GetById(userid, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userid, listId)
}
