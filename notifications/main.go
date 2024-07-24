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

type OrderStatusData struct {
	OrderID  string `json:"orderID" bson:"orderID"`
	Status      string `json:"status" bson:"status"`
	RiderUsername string `json:"riderUsername" bson:"riderUsername"`
	RiderEmail string `json:"riderEmail" bson:"riderEmail"`
	Message string `json:"message" bson:"message"`
	Subject string `json:"subject" bson:"subject"`
}


type EmailVerify struct {
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Message string `json:"message" bson:"message"`
	Subject string `json:"subject" bson:"subject"`
}


var mongoString = config.GetEnv("MONGO_STRING", "mongodb://user:user@localhost:27017")


func main() {
	ch, conn := broker.ConnectRabbit()
	defer conn.Close()
    defer ch.Close()
	// connect mongodb
	client, err := connectMongoDB(mongoString)
	if err!= nil {
        log.Fatal(err)
    }
	var forever chan struct{}
	// Consume messages from the "order_notification" queue
	listenOrderNotification(ch, client)
	// Consume messages from the "order_status_change" queue
	listenOrderStatusChange(ch, client)
	listenEmailVerification(ch, client)

	<-forever
	
	
}

func listenOrderNotification(ch *amqp091.Channel,  mongoClient *mongo.Client) {
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
	// var forever chan struct{}
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
			// var doc bson.D
			// bsonData, _ := bson.Marshal(d.Body)
    		// err = bson.Unmarshal(bsonData, &doc)
			insertResult, err := collection.InsertOne(context.TODO(), data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted a single document: ", insertResult.InsertedID)	
        }	
    }()
	log.Printf(" [*] Waiting for messages from %s queue. To exit press CTRL+C \n", q.Name)
	// <-forever

}

func listenOrderStatusChange(ch *amqp091.Channel, mongoClient *mongo.Client){
    q, err := ch.QueueDeclare(
        "delivery_status_change",
        false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare order status change: %v", err)
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
		log.Fatalf("Fail to consume %s: %v", q.Name, err)
	}
	collection := mongoClient.Database("notifications").Collection("order_status_notifications")
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var data OrderStatusData
			if err := json.Unmarshal(d.Body, &data); err!= nil {
                log.Fatalf("Error unmarshalling JSON: %s", err)
                continue
            }
			// send mail 
			go utils.SendMail(data.RiderEmail, data.Subject, data.RiderUsername, data.Message)
			insertResult, err := collection.InsertOne(context.TODO(), data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Inserted a single document: ", insertResult.InsertedID)	

		}
	}()
	log.Printf(" [*] Waiting for messages from %s queue. To exit press CTRL+C \n", q.Name)
}

// listen to email_verification
func listenEmailVerification(ch *amqp091.Channel, mongoClient *mongo.Client) {
	q, err := ch.QueueDeclare(
        "email_verification",
        false,
        false,
        false,
        false,
        nil,
    )
    if err!= nil {
        log.Fatalf("Failed to declare email verification: %v", err)
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
    if err!= nil {
        log.Fatalf("Failed to consume %s: %v", q.Name, err)
    }
    collection := mongoClient.Database("notifications").Collection("email_verification_notifications")
	go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
            var data EmailVerify
            if err := json.Unmarshal(d.Body, &data); err!= nil {
                log.Fatalf("Error unmarshalling JSON: %s", err)
                continue
            }
            // send mail
            go utils.SendMail(data.Email, data.Subject, data.Username, data.Message)
            insertResult, err := collection.InsertOne(context.TODO(), data)
            if err!= nil {
                log.Fatal(err)
            }
            fmt.Println("Inserted a single document: ", insertResult.InsertedID)    
        }
    }()
    log.Printf(" [*] Waiting for messages from %s queue. To exit press CTRL+C \n", q.Name)

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