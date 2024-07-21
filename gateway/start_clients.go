package main

import (
	"log"
	"net/http"

	"github.com/Ayobami6/common/broker"
	"github.com/Ayobami6/common/config"
	pbOrder "github.com/Ayobami6/common/proto/orders"
	pbRider "github.com/Ayobami6/common/proto/riders"
	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var httpAddr = config.GetEnv("GATEWAY_PORT", ":2330")

func StartGateway() {
	ch, coon := broker.ConnectRabbit() 
	defer coon.Close()
	defer ch.Close()
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v2").Subrouter()
	userServiceAddr, err := utils.GetServiceAddress("user-service")
	if err != nil {
			log.Fatalf("Error getting service address: %v", err)
	}
	uconn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		    log.Fatalf("Failed to connect: %v", err)
	}
	defer uconn.Close()
	log.Println("Dailing user service at ", userServiceAddr)
	uClient := pbUser.NewUserServiceClient(uconn)
	userHandler := NewUserClientHandler(uClient)
	userHandler.RegisterRoutes(subrouter)

	log.Printf("Server is listening on %s", httpAddr)
	RiderServiceAddr, err := utils.GetServiceAddress("rider-service")
	if err != nil {
		log.Fatalf("Error getting service address: %v", err)
	}
	rConn, err := grpc.Dial(RiderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
			log.Fatalf("Failed to connect: %v", err)
	}
	defer rConn.Close()
	log.Println("Dailing rider service at ", RiderServiceAddr)
	rC := pbRider.NewRiderServiceClient(rConn)
	uC := pbUser.NewUserServiceClient(uconn)
	riderHandler := NewRiderClientHandler(rC, uC)
	riderHandler.RegisterRoutes(subrouter)

	// Order Service
	orderServiceAddr, err := utils.GetServiceAddress("order-service")
	if err!= nil {
        log.Fatalf("Error getting service address: %v", err)
    }
	oConn, err := grpc.Dial(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Fail to transport OrderService")
	}
	defer oConn.Close()
	log.Println("Dailing OrderService at ", orderServiceAddr)
	oC := pbOrder.NewOrderServiceClient(oConn)
	rNC := pbRider.NewRiderServiceClient(rConn)
	uNC := pbUser.NewUserServiceClient(uconn)
	orderHandler := NewOrderClientHandler(oC, uNC, rNC, ch)
	orderHandler.RegisterRoutes(subrouter)
	log.Printf("Server is listening on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, subrouter); err != nil {
		log.Fatal(err)
	}
}
