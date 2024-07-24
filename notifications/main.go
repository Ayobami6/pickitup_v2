package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Ayobami6/common/broker"
	"github.com/Ayobami6/common/config"
	"github.com/Ayobami6/common/utils"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type EmailData struct {
    UserEmail    string `json:"userEmail" bson:"userEmail"`
    Subject      string `json:"subject" bson:"subject"`
    UserUsername string `json:"userUsername" bson:"userUsername"`
    UserMessage  string `json:"userMessage" bson:"userMessage"`
	RiderMessage string `json:"riderMessage" bson:"riderMessage"`
	OrderID  string `json:"orderID" bson:"orderID"`
	RiderName string `json:"riderName" bson:"riderName"`
	RiderEmail string `json:"riderEmail" bson:"riderEmail"`
}


var mongoString = config.GetEnv("MONGO_STRING", "mongodb://user:user@localhost:27017")


func main() {
	ch, conn := broker.ConnectRabbit()
	// connect mongodb
	client, err := connectMongoDB(mongoString)
	if err!= nil {
        log.Fatal(err)
    }
	// Consume messages from the "order_notification" queue
	listenOrderNotification(ch, conn, client)
	
	
}

func listenOrderNotification(ch *amqp091.Channel, conn *amqp091.Connection, mongoClient *mongo.Client) {
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
	collection := mongoClient.Database("notifications").Collection("order_notifications")
	var forever chan struct{}
	go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
			// send the notification email
			var data EmailData
            if err := json.Unmarshal(d.Body, &data); err != nil {
                log.Fatalf("Error unmarshalling JSON: %s", err)
                continue
            }
			go utils.SendMail(data.UserEmail, data.Subject, data.UserUsername, data.UserMessage)
			go utils.SendMail(data.RiderEmail, data.Subject, data.RiderName, data.RiderMessage)
			var doc bson.D
			bsonData, _ := bson.Marshal(d.Body)
    		err = bson.Unmarshal(bsonData, &doc)
			insertResult, err := collection.InsertOne(context.TODO(), doc)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted a single document: ", insertResult.InsertedID)	
        }	
    }()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever // hang until manually closed

}


func connectMongoDB(db string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(db)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err!= nil {
        return nil, err
    }
    err = client.Ping(context.TODO(), nil)
    if err!= nil {
        return nil, err
    }
    log.Println("Connected to MongoDB!")
    return client, nil

}