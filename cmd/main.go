package main

import (
	"os"

	todo_app "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/handler"
	"github.com/Daniil-1622/todo-app/pkg/repository"
	"github.com/Daniil-1622/todo-app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("init config err: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("load .env file err: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to init db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo_app.Server)                                                     // Создаем новый экземпляр структуры Server
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil { // Запускаемся на порту 8000
		// Если возникает ошибка, выводим ее
		logrus.Fatalf("error occured while running server: %s", err.Error())
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
