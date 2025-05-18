# Практическое задание 12: RabbitMQ

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