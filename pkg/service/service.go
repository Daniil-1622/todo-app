package service

import (
	todo "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userid int) ([]todo.TodoList, error)
	GetById(userid, listId int) (todo.TodoList, error)
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		TodoList:      NewTodoListService(repo.TodoList),
	}
}
