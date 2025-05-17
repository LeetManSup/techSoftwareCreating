# Практическое задание 7: Кеширование и работа с очередями

## 1. Redis как кеш

Код проекта доступен [по ссылке](https://github.com/LeetManSup/mirea_finance_tracker/).

Установка зависимости:

```bash
go get github.com/redis/go-redis/v9
```

Определение контейнера с Redis в `docker-compose.yml`:

```yml
...

  redis:
    image: redis:7
    container_name: redis
    restart: always
    ports:
      - "6380:6379"

...
```

Redis-клиент `internal/redis/redis.go`:

```go
package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	Client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
}
```

Адрес доступа к контейнеризированному redis в файле `.env`:
```
...

REDIS_ADDR=redis:6379
```

Переработанный обработчик `GetAccount` файла `account_handler.go` с кешированием результатов на 10 минут:

```go
func (h *AccountHandler) GetAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	accountID := c.Param("id")
	cacheKey := fmt.Sprintf("account:%s:user:%s", accountID, userID.(string))

	val, err := redis.Client.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		var acc model.Account
		if err := json.Unmarshal([]byte(val), &acc); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"id":              acc.ID,
				"name":            acc.Name,
				"currency_code":   acc.CurrencyCode,
				"initial_balance": acc.InitialBalance,
				"created_at":      acc.CreatedAt,
			})
			return
		}
	}

	account, err := h.accountService.GetAccountByID(userID.(string), accountID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	data, _ := json.Marshal(account)
	redis.Client.Set(redis.Ctx, cacheKey, data, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"id":              account.ID,
		"name":            account.Name,
		"currency_code":   account.CurrencyCode,
		"initial_balance": account.InitialBalance,
		"created_at":      account.CreatedAt,
	})
}
```

Инициализация Redis в функции `main`:
```go
func main() {
	redis.InitRedis()
    ...
}
```

### Подтверждение работы

Несколько раз подряд выполнялись запросы на просмотр одного и того же счёта пользователя:
![alt text](images/image1.png)

Отчёты о выполнении запросов:
![alt text](images/image2.png)

Можно видеть, что время выполнения первых двух запросов, исчисляется в миллисекундах, они приводили к запросу в БД. Время выполнения последующих запросов уменьшилось на несколько порядков и исчиляется уже в микросекундах, связано это с тем, что результаты были быстро взяты из кеша redis, без выполнения запроса к БД.

Иллюстрация кеша, хранимого в redis, который был использован для быстрых ответов на запросы:
![alt text](images/image3.png)

Выполним запрос спустя время (когда кеш должен исчезнуть) и ознакомимся с отчётом исполнения:
![alt text](images/image4.png)

Запрос снова обрабатывался миллисекунды.


**Дальнейшие реализации не нашли себя в проекте трекинга личных финансов, поэтому будут реализованы в виде демонстрационных проектов**

## 2. Redis pub/sub

Код проекта доступен [по ссылке](https://github.com/LeetManSup/tsc-p7-redis-pubsub/).

Структура проекта:

![alt text](images/image5.png)

Пример образа публикатора `publisher/Dockerfile` (образ читателя построен по аналогии):
```Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./publisher .

RUN go build -o publisher .

CMD ["./publisher"]
```

Публикатор пишет 5 сообщений в топик `notifications` с интервалом в 2 секунды и завершает работу:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	topic := "notifications"

	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Hello from publisher! Message #%d", i)
		if err := rdb.Publish(ctx, topic, message).Err(); err != nil {
			fmt.Printf("Publish error: %v\n", err)
		} else {
			fmt.Printf("Published: %s\n", message)
		}
		time.Sleep(2 * time.Second)
	}
}
```

Читатель подписан на этот топик и выводит сообщения, которые читает из него:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	topic := "notifications"
	sub := rdb.Subscribe(ctx, topic)

	ch := sub.Channel()

	fmt.Println("Subscribed to", topic)

	for msg := range ch {
		fmt.Printf("Received message: %s\n", msg.Payload)
	}
}
```

Подтверждение работы:

![alt text](images/image6.png)

Вывод недетерминирован, поскольку процессы параллельны.

## 3. RabbitMQ

Код проекта доступен [по ссылке](https://github.com/LeetManSup/tsc-p7-rabbitmq-pubwork/).

Структура проекта:

![alt text](images/image7.png)

Образы Dockerfile публикатора и воркера построены по аналогии с образами из предыдущего раздела.

В связи с долгим запуском rabbit, в публикаторе и воркере реализована функция ожидания его запуска.

Публикатор публикует задачи отправки электронных писем в очередь `tasks` канала соединения, имитируемые текстом в body, и завершает работу:
```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	conn := waitForRabbitMQ(amqpURL, 10, 3*time.Second)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks", // name
		true,    // durable
		false,   // auto-delete
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	for i := 1; i <= 5; i++ {
		body := fmt.Sprintf("Task #%d: Send email to user%d@example.com", i, i)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			log.Printf("Published: %s", body)
		}
		time.Sleep(1 * time.Second)
	}
}

func waitForRabbitMQ(amqpURL string, retries int, delay time.Duration) *amqp.Connection {
	for i := 1; i <= retries; i++ {
		conn, err := amqp.Dial(amqpURL)
		if err == nil {
			return conn
		}
		log.Printf("RabbitMQ not ready (attempt %d/%d): %v", i, retries, err)
		time.Sleep(delay)
	}
	log.Fatal("RabbitMQ not available after retries")
	return nil
}
```

Воркер читает из очереди `tasks` задачи и имитирует их выполнение путём вывода информационных сообщений в консоль:
```go
package main

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	conn := waitForRabbitMQ(amqpURL, 10, 3*time.Second)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Worker started, waiting for tasks...")

	for msg := range msgs {
		log.Printf("Processing: %s", msg.Body)
		time.Sleep(2 * time.Second) // имитация работы
		log.Println("Done.")
	}
}

func waitForRabbitMQ(amqpURL string, retries int, delay time.Duration) *amqp.Connection {
	for i := 1; i <= retries; i++ {
		conn, err := amqp.Dial(amqpURL)
		if err == nil {
			return conn
		}
		log.Printf("RabbitMQ not ready (attempt %d/%d): %v", i, retries, err)
		time.Sleep(delay)
	}
	log.Fatal("RabbitMQ not available after retries")
	return nil
}
```

Подтверждение работы:

![alt text](images/image8.png)

Вывод недетерминирован, поскольку процессы параллельны.

## 4. CQRS

Код проекта доступен [по ссылке](https://github.com/LeetManSup/tsc-p7-cqrs/).

Структура проекта:

![alt text](images/image9.png)

`docker-compose.yml`:
```yml
services:
  cqrs:
    build: .
    ports:
      - "8081:8080"
```

Образ Dockerfile построен по аналогии с образами из предыдущих разделов.

Файл main.go реализует демонстрационный микросервис на основе архитектурного подхода CQRS с event sourcing. В нём команды (создать, оплатить, отменить заказ) порождают события, которые сохраняются в лог, а текущее состояние заказов (read model) восстанавливается при запуске путём последовательного применения этих событий. Взаимодействие с системой осуществляется через REST API, реализованный с использованием пакета gorilla/mux, без внешних зависимостей и хранилищ — все данные хранятся в оперативной памяти:
```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderStatus string

const (
	StatusPending  OrderStatus = "PENDING"
	StatusPaid     OrderStatus = "PAID"
	StatusCanceled OrderStatus = "CANCELED"
)

// --- Events ---
type EventType string

const (
	EventOrderCreated  EventType = "OrderCreated"
	EventOrderPaid     EventType = "OrderPaid"
	EventOrderCanceled EventType = "OrderCanceled"
)

type Event struct {
	Type      EventType       `json:"type"`
	OrderID   string          `json:"order_id"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// --- Read model (in-memory) ---
type Order struct {
	ID     string      `json:"id"`
	Status OrderStatus `json:"status"`
}

var (
	eventLog []Event              // упрощённый event store
	orders   = map[string]Order{} // read model
	mutex    sync.Mutex
)

// --- Command Handlers ---
func createOrder(w http.ResponseWriter, r *http.Request) {
	orderID := uuid.New().String()
	event := Event{
		Type:      EventOrderCreated,
		OrderID:   orderID,
		Timestamp: time.Now(),
		Data:      json.RawMessage(`{}`),
	}
	appendEvent(event)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"order_id": orderID})
}

func payOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	event := Event{
		Type:      EventOrderPaid,
		OrderID:   orderID,
		Timestamp: time.Now(),
		Data:      json.RawMessage(`{}`),
	}
	appendEvent(event)
	w.WriteHeader(http.StatusNoContent)
}

func cancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	event := Event{
		Type:      EventOrderCanceled,
		OrderID:   orderID,
		Timestamp: time.Now(),
		Data:      json.RawMessage(`{}`),
	}
	appendEvent(event)
	w.WriteHeader(http.StatusNoContent)
}

// --- Query Handlers ---
func getOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	mutex.Lock()
	order, ok := orders[orderID]
	mutex.Unlock()

	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(eventLog)
}

// --- Event Store & Projection ---
func appendEvent(e Event) {
	mutex.Lock()
	defer mutex.Unlock()
	eventLog = append(eventLog, e)
	applyEvent(e)
}

func applyEvent(e Event) {
	switch e.Type {
	case EventOrderCreated:
		orders[e.OrderID] = Order{ID: e.OrderID, Status: StatusPending}
	case EventOrderPaid:
		if o, ok := orders[e.OrderID]; ok {
			o.Status = StatusPaid
			orders[e.OrderID] = o
		}
	case EventOrderCanceled:
		if o, ok := orders[e.OrderID]; ok {
			o.Status = StatusCanceled
			orders[e.OrderID] = o
		}
	}
}

// --- Init ---
func rebuildState() {
	for _, e := range eventLog {
		applyEvent(e)
	}
}

func main() {
	rebuildState()
	r := mux.NewRouter()

	// Команды
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}/pay", payOrder).Methods("POST")
	r.HandleFunc("/orders/{id}/cancel", cancelOrder).Methods("POST")

	// Запросы
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/events", getAllEvents).Methods("GET")

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
```

Подтверждение работы:

Создание заказа:
![alt text](images/image10.png)
Оплата заказа:
![alt text](images/image11.png)
Отмена заказа:
![alt text](images/image12.png)
Получение текущего состояния заказа:
![alt text](images/image13.png)
Просмотр лога событий:
![alt text](images/image14.png)
