package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Producer — интерфейс. Описывает что умеет продюсер.
// Благодаря интерфейсу можно подменить реальный Kafka на mock в тестах.
type Producer interface {
	Publish(topic string, event Event) error
	Close() error
}

// KafkaProducer — реальная реализация Producer
type KafkaProducer struct {
	brokers []string
	writers map[string]*kafkago.Writer // один writer на каждый topic
}

// NewProducer — создаёт нового продюсера и подключается к Kafka
func NewProducer(brokers []string) (*KafkaProducer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers list is empty")
	}

	return &KafkaProducer{
		brokers: brokers,
		writers: make(map[string]*kafkago.Writer),
	}, nil
}

// getWriter — возвращает writer для topic'а.
// Если writer ещё не создан — создаёт новый (lazy initialization)
func (p *KafkaProducer) getWriter(topic string) *kafkago.Writer {
	if writer, ok := p.writers[topic]; ok {
		return writer
	}

	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(p.brokers...),
		Topic:        topic,
		Balancer:     &kafkago.Hash{}, // сообщения одного user_id идут в одну партицию
		MaxAttempts:  3,               // 3 попытки при ошибке
		BatchTimeout: 10 * time.Millisecond,
		Compression:  kafkago.Gzip, // сжимаем сообщения
	}

	p.writers[topic] = writer
	return writer
}

// Publish — отправляет событие в указанный topic
func (p *KafkaProducer) Publish(topic string, event Event) error {
	// Превращаем структуру Event в JSON байты
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Ключ сообщения = user_id. Гарантирует что события одного
	// пользователя всегда попадают в одну партицию (порядок сохраняется)
	key := []byte(fmt.Sprintf("%d", event.UserId))

	msg := kafkago.Message{
		Key:   key,
		Value: data,
	}

	// Пробуем отправить 3 раза с нарастающей задержкой
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = p.getWriter(topic).WriteMessages(ctx, msg)
		cancel()

		if err == nil {
			logrus.Infof("[Kafka Producer] Event published: topic=%s, event_type=%s, user_id=%d",
				topic, event.EventType, event.UserId)
			return nil
		}

		lastErr = err
		logrus.Warnf("[Kafka Producer] Publish retry: topic=%s, attempt=%d, error=%s",
			topic, attempt, err.Error())

		// Задержка: 100ms, 200ms, 300ms
		time.Sleep(time.Duration(100*attempt) * time.Millisecond)
	}

	logrus.Errorf("[Kafka Producer] Publish failed: topic=%s, error=%s", topic, lastErr.Error())
	return fmt.Errorf("failed to publish after 3 attempts: %w", lastErr)
}

// Close — закрывает все writer'ы. Вызывается при завершении приложения.
func (p *KafkaProducer) Close() error {
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			logrus.Errorf("[Kafka Producer] Error closing writer for topic %s: %s", topic, err.Error())
		}
	}
	return nil
}
