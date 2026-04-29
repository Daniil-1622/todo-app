package events

import (
	"context"
	"encoding/json"

	"github.com/Daniil-1622/todo-app/pkg/kafka"
	"github.com/sirupsen/logrus"
)

// EventProcessor — главный менеджер. Запускает чтение из всех топиков
// и направляет каждое сообщение в нужный обработчик.
type EventProcessor struct {
	consumer    kafka.Consumer
	taskHandler *TaskEventHandler
	listHandler *ListEventHandler
	userHandler *UserEventHandler
}

// NewEventProcessor — создаёт новый процессор
func NewEventProcessor(
	consumer kafka.Consumer,
	taskHandler *TaskEventHandler,
	listHandler *ListEventHandler,
	userHandler *UserEventHandler,
) *EventProcessor {
	return &EventProcessor{
		consumer:    consumer,
		taskHandler: taskHandler,
		listHandler: listHandler,
		userHandler: userHandler,
	}
}

// Start — запускает чтение из всех топиков.
// Вызывается один раз при старте приложения в отдельной goroutine.
func (p *EventProcessor) Start(ctx context.Context, topics kafka.TopicsConfig) {
	// Запускаем чтение каждого топика в своей goroutine
	go p.processTopics(ctx, topics.TaskEvents, p.taskHandler.Handle)
	go p.processTopics(ctx, topics.ListEvents, p.listHandler.Handle)
	go p.processTopics(ctx, topics.UserEvents, p.userHandler.Handle)

	logrus.Info("[EventProcessor] Started processing topics")

	// Ждём сигнала остановки
	<-ctx.Done()
	logrus.Info("[EventProcessor] Stopped")
}

// processTopics — читает сообщения из одного топика и вызывает handler
func (p *EventProcessor) processTopics(
	ctx context.Context,
	topic string,
	handler func(event kafka.Event) error,
) {
	// Получаем канал с сообщениями из Kafka
	msgChan, err := p.consumer.Read(ctx, topic)
	if err != nil {
		logrus.Errorf("[EventProcessor] Failed to read topic %s: %s", topic, err.Error())
		return
	}

	logrus.Infof("[EventProcessor] Listening topic: %s", topic)

	// Читаем сообщения из канала пока он не закроется
	for msg := range msgChan {
		// Превращаем JSON байты обратно в структуру Event
		var event kafka.Event
		if err := json.Unmarshal(msg, &event); err != nil {
			logrus.Errorf("[EventProcessor] Failed to unmarshal event: %s", err.Error())
			continue // пропускаем битое сообщение, читаем следующее
		}

		// Передаём событие в нужный обработчик
		if err := handler(event); err != nil {
			logrus.Errorf("[EventProcessor] Handler error for event %s: %s",
				event.EventId, err.Error())
			continue
		}
	}
}
