package main

import (
	"log"
	"net/http"

	"github.com/Ayobami6/common/config"
	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var httpAddr = config.GetEnv("GATEWAY_PORT", ":2330")

func StartUserServiceClient() {
	serviceAddr, err := utils.GetServiceAddress("user-service")
	if err != nil {
			log.Fatalf("Error getting service address: %v", err)
	}
	conn, err := grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		    log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Println("Dailing user service at ", serviceAddr)
	c := pbUser.NewUserServiceClient(conn)
	handler := NewUserClientHandler(c)
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v2").Subrouter()
	handler.RegisterRoutes(subrouter)

	log.Printf("Server is listening on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, subrouter); err != nil {
		log.Fatal(err)
	}

}