package kafka

type Config struct {
	Brokers []string

	Producer ProducerConfig
	Consumer ConsumerConfig
	Topics   TopicsConfig
}

type ProducerConfig struct {
	Compression string
	RetryMax    int
	BatchSize   int
}

type ConsumerConfig struct {
	GroupID          string
	SessionTimeoutMs int
	MaxPollRecords   int
}

type TopicsConfig struct {
	TaskEvents string
	ListEvents string
	UserEvents string
}

/*
Config — главная структура. Она содержит всё необходимое для подключения к Kafka. Когда приложение стартует, оно читает config.yml и заполняет эту структуру.
Brokers — список адресов серверов Kafka. Kafka может работать на нескольких серверах сразу (для надёжности), поэтому это слайс. В твоём случае будет просто один: localhost:9092.
ProducerConfig — настройки для отправщика сообщений:

Compression — сжимать ли сообщения перед отправкой (gzip). Экономит трафик.
RetryMax — если Kafka временно недоступна, сколько раз попробовать снова (будет 3).
BatchSize — отправлять не по одному сообщению, а пачками по 100. Быстрее.

ConsumerConfig — настройки для читателя сообщений:

GroupID — имя группы потребителей. Kafka запоминает что группа уже прочитала, чтобы при перезапуске не читать заново с начала.
SessionTimeoutMs — если consumer завис и не отвечает 30 секунд, Kafka считает его мёртвым.
MaxPollRecords — за один раз читать не более 500 сообщений.

TopicsConfig — названия очередей. Три очереди: для задач, для списков, для пользователей. Выносим в конфиг чтобы не хардкодить строки по всему коду.
*/
