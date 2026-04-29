package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	todo_app "github.com/Daniil-1622/todo-app"
	"github.com/Daniil-1622/todo-app/pkg/events"
	"github.com/Daniil-1622/todo-app/pkg/handler"
	"github.com/Daniil-1622/todo-app/pkg/kafka"
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

	// [НОВОЕ] Инициализация Kafka Producer
	producer, err := kafka.NewProducer(
		viper.GetStringSlice("kafka.brokers"),
	)
	if err != nil {
		logrus.Fatalf("failed to init kafka producer: %s", err.Error())
	}

	// [НОВОЕ] Инициализация Kafka Consumer
	consumer, err := kafka.NewKafkaConsumer(
		viper.GetStringSlice("kafka.brokers"),
		viper.GetString("kafka.consumer.group_id"),
	)
	if err != nil {
		logrus.Fatalf("failed to init kafka consumer: %s", err.Error())
	}

	// [НОВОЕ] Инициализация EventProcessor и запуск в горутине
	processor := events.NewEventProcessor(
		consumer,
		events.NewTaskEventHandler(),
		events.NewListEventHandler(),
		events.NewUserEventHandler(),
	)

	topics := kafka.TopicsConfig{
		TaskEvents: viper.GetString("kafka.topics.task_events"),
		ListEvents: viper.GetString("kafka.topics.list_events"),
		UserEvents: viper.GetString("kafka.topics.user_events"),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go processor.Start(ctx, topics)

	// [ИЗМЕНЕНО] NewService теперь принимает producer
	repos := repository.NewRepository(db)
	services := service.NewService(repos, producer)

	// [ИЗМЕНЕНО] NewHandler теперь принимает producer
	handlers := handler.NewHandler(services, producer)

	srv := new(todo_app.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp stopping...")

	// [НОВОЕ] Останавливаем EventProcessor через контекст
	cancel()

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured while shutting down server: %s", err.Error())
	}

	// [НОВОЕ] Закрываем Producer
	if err := producer.Close(); err != nil {
		logrus.Errorf("error occured while closing producer: %s", err.Error())
	}

	// [НОВОЕ] Закрываем Consumer
	if err := consumer.Close(); err != nil {
		logrus.Errorf("error occured while closing consumer: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured while closing db: %s", err.Error())
	}

	logrus.Print("TodoApp stopped")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
