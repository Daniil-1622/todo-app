package main

import (
	"log"

	todo_app "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/handler"
	"github.com/Daniil-1622/todo-app/pkg/repository"
	"github.com/Daniil-1622/todo-app/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("init config err: %s", err.Error())
	}
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo_app.Server)                                                 // Создаем новый экземпляр структуры Server
	if err := srv.Run(viper.GetString(""), handlers.InitRoutes()); err != nil { // Запускаемся на порту 8000
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
