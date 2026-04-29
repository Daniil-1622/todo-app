package kafka

import "time"

type EventType string

const (
	TopicTaskEvents = "task-events"
	TopicListEvents = "list-events"
	TopicUserEvents = "user-events"
	// Событие задач
	EventTaskCreated   EventType = "TaskCreated"
	EventTaskUpdated   EventType = "TaskUpdated"
	EventTaskDeleted   EventType = "TaskDeleted"
	EventTaskCompleted EventType = "TaskCompleted"
	// Событие списков
	EventListCreated EventType = "ListCreated"
	EventListUpdated EventType = "ListUpdated"
	EventListDeleted EventType = "ListDeleted"
	// Событие пользователя
	EventUserRegistered EventType = "UserRegistered"
	EventUserLoggedIn   EventType = "UserLoggedIn"
)

type Event struct {
	EventId   string                 `json:"event_id"`
	EventType EventType              `json:"event_type"`
	TimeStamp time.Time              `json:"timestamp"`
	UserId    int                    `json:"user_id"`
	ListId    int                    `json:"list_id,omitempty"`
	TaskId    int                    `json:"task_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
}

/*
EventType string — создаём свой тип на основе обычной строки. Смысл в том, что теперь компилятор Go не даст передать случайную строку туда где ожидается EventType. Только наши константы.
Константы — все возможные события в системе. Когда в коде будешь отправлять событие, напишешь не строку "TaskCreated" а константу kafka.EventTaskCreated. Если опечатаешься — компилятор сразу скажет об ошибке.
Event — главная структура. Это одно сообщение которое полетит в Kafka в виде JSON. Разберём каждое поле:

EventId — уникальный ID события. Будем генерировать через uuid. Нужен чтобы не обработать одно и то же событие дважды
EventType — что произошло: TaskCreated, ListDeleted и т.д.
TimeStamp — точное время события
UserId — кто совершил действие. Важно: это будет ключом сообщения в Kafka, что гарантирует порядок событий одного пользователя
ListId / TaskId — что именно затронуто
Data — гибкое поле map[string]interface{}. Сюда пойдут любые доп. данные. Для TaskCreated это будет {"title": "Buy milk", "description": "...", "done": false}
*/
