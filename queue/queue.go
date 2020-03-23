package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Connect (*amqp.Channel, error) {
// retorna um canal do rabbitmq
func Connect() *amqp.Channel {
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")

	// estabelece conexão com o rabbitmq
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatalf("%s: Failed to connect to RabbitMQ.", err)
		// return nil, errors.New(fmt.Sprintf("%s: Failed to connect to RabbitMQ.", err))
	}

	// cria um canal através da conexão com o rabbitmq
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: Failed to open a channel.", err)
		// return nil, errors.New(fmt.Sprintf("%s: Failed to open a channel", err))
	}

	// return ch, nil
	return ch
}

// StartConsuming consuming a queue
func StartConsuming(in chan []byte, ch *amqp.Channel) {

	// declara uma fila
	queue, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_CONSUMER_QUEUE"), // name
		true,                                 // durable
		false,                                // delete when used
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)
	if err != nil {
		log.Fatalf("%s: Failed to declare queue", err)
	}

	// fica aguardando messagens que serão armazenadas em 'msgs'
	msgs, err := ch.Consume(
		queue.Name,  // queue
		"go-worker", // cosumer
		true,        // auto-ack - recebe a mensagem e esvazia/sai da fila
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("%s: Failed to register a consumer", err)
	}

	// verifica assincronamente a chegada de uma nova mensagem
	// e despacha para o canal 'in'
	go func() {
		for d := range msgs {
			in <- []byte(d.Body)
		}
		close(in)
	}()
}

// Notify publica a mensagem em forma de json
func Notify(payload string, ch *amqp.Channel) {
	err := ch.Publish(
		os.Getenv("RABBITMQ_DESTINATION_POSITIONS_EXCHANGE"), // exchange
		os.Getenv("RABBITMQ_DESTINATION_ROUTING_KEY"),        // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		},
	)
	if err != nil {
		log.Fatalf("%s: Failed to publish a message", err)
	}
	fmt.Println("Message sent: ", payload)
}
