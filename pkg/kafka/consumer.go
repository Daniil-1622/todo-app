package kafka

import (
	"context"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer interface {
	Read(ctx context.Context, topic string) (<-chan []byte, error)
	Close() error
}

type KafkaConsumer struct {
	brokers []string
	groupId string
	readers map[string]*kafkago.Reader
}

// NewConsumer — создаёт нового consumer'а
func NewKafkaConsumer(brokers []string, groupId string) (*KafkaConsumer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers list is empty")
	}

	if len(groupId) == 0 {
		return nil, fmt.Errorf("kafka groupId is empty")
	}

	return &KafkaConsumer{
		brokers: brokers,
		groupId: groupId,
		readers: make(map[string]*kafkago.Reader),
	}, nil
}

// getReader — возвращает reader для topic'а.
// Если reader ещё не создан — создаёт новый (lazy initialization)
func (c *KafkaConsumer) getReader(topic string) *kafkago.Reader {
	if reader, ok := c.readers[topic]; ok {
		return reader
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        c.brokers,
		Topic:          topic,
		GroupID:        c.groupId,   // группа — Kafka запомнит что мы уже прочитали
		MinBytes:       1,           // минимум байт для чтения
		MaxBytes:       10e6,        // максимум 10MB за раз
		MaxWait:        time.Second, // ждать максимум 1 секунду новых сообщений
		CommitInterval: time.Second, // автоматически коммитить каждую секунду
	})

	c.readers[topic] = reader
	return reader
}

// Read — читает сообщения из topic'а и отправляет их в канал.
// Работает в фоновой goroutine, возвращает канал с сырыми JSON байтами.
func (c *KafkaConsumer) Read(ctx context.Context, topic string) (<-chan []byte, error) {
	reader := c.getReader(topic)

	// Создаём канал — через него будем передавать сообщения обработчику
	msgChan := make(chan []byte, 100) // буфер 100 сообщений

	go func() {
		defer close(msgChan) // закрываем канал когда goroutine завершается

		for {
			// Проверяем не отменён ли контекст (graceful shutdown)
			select {
			case <-ctx.Done():
				logrus.Infof("[Kafka Consumer] Stopping reader for topic: %s", topic)
				return
			default:
			}

			// Читаем следующее сообщение из Kafka
			// Блокируется пока не придёт новое сообщение
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					// Контекст отменён — нормальное завершение
					return
				}
				logrus.Errorf("[Kafka Consumer] Error reading message from topic %s: %s",
					topic, err.Error())
				// Небольшая пауза перед повторной попыткой
				time.Sleep(time.Second)
				continue
			}

			logrus.Infof("[Kafka Consumer] Message consumed: topic=%s, offset=%d",
				topic, msg.Offset)

			// Отправляем сообщение в канал с таймаутом
			select {
			case msgChan <- msg.Value:
				// сообщение успешно отправлено в канал
			case <-ctx.Done():
				return
			case <-time.After(30 * time.Second):
				logrus.Warnf("[Kafka Consumer] Message processing timeout, topic=%s", topic)
			}
		}
	}()

	return msgChan, nil
}

// Close — закрывает все reader'ы
func (c *KafkaConsumer) Close() error {
	for topic, reader := range c.readers {
		if err := reader.Close(); err != nil {
			logrus.Errorf("[Kafka Consumer] Error closing reader for topic %s: %s",
				topic, err.Error())
		}
	}
	return nil
}
