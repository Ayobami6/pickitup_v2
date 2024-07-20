package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbit() (conn *amqp.Connection){
	fmt.Println("Connecting to RabbitMQ...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err!= nil {
        panic(err)
    }
  	defer conn.Close()
	return conn
}