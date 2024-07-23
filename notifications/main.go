package main

import (
	"encoding/json"
	"log"

	"github.com/Ayobami6/common/broker"
	"github.com/Ayobami6/common/utils"
	"github.com/rabbitmq/amqp091-go"
)


type EmailData struct {
    UserEmail    string `json:"userEmail"`
    Subject      string `json:"subject"`
    UserUsername string `json:"userUsername"`
    UserMessage  string `json:"userMessage"`
	RiderMessage string `json:"riderMessage"`
	OrderID  string `json:"orderID"`
	RiderName string `json:"riderName"`
	RiderEmail string `json:"riderEmail"`
}


func main() {
	ch, conn := broker.ConnectRabbit()
	// Consume messages from the "order_notification" queue
	listenOrderNotification(ch, conn)
	
	
}

func listenOrderNotification(ch *amqp091.Channel, conn *amqp091.Connection) {
	defer conn.Close()
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"order_notification",
		false,
		false,
		false,
		false,
		nil,
	  )
	if err!= nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }
	msgs, err := ch.Consume(
        q.Name,
        "",   
        true, 
        false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
	var forever chan struct{}
	go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
			// send the notification email
			var data EmailData
            if err := json.Unmarshal(d.Body, &data); err != nil {
                log.Printf("Error unmarshalling JSON: %s", err)
                continue
            }
			go utils.SendMail(data.UserEmail, data.Subject, data.UserUsername, data.UserMessage)
			go utils.SendMail(data.RiderEmail, data.Subject, data.RiderName, data.RiderMessage)
        }
    }()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever // hang until manually closed

}