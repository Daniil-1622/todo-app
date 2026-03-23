package main

import (
	"log"

	todo_app "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(todo_app.Server)                                    // Создаем новый экземпляр структуры Server
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil { // Запускаемся на порту 8000
		// Если возникает ошибка, выводим ее
		log.Fatalf("error occured while running server: %s", err.Error())
	}
	/*
			После добавления:
			handlers := new(handler.Handler)
			srv.Run("8000", handlers.InitRoutes())
		Сервер теперь знает что делать и по каким маршрутам двигаться!
	*/
}
