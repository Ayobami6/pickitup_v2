package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/Ayobami6/common/config"
	"github.com/Ayobami6/common/db"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

const (
	checkID = "alive"
	ttl     = time.Second * 8
)
var servicePort = config.GetEnv("RIDER_SERVICE_PORT", "5001")

func registerRiderService() {
	conf := api.DefaultConfig()
	client, err := api.NewClient(conf)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}
	port, err := strconv.Atoi(servicePort)
	if err != nil {
		log.Fatal("Invalid port number")
	}

	registration := &api.AgentServiceRegistration{
		ID:      "rider-service-1",
		Name:    "rider-service",
		Address: "localhost",
		Port:	port,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "5m",
			CheckID:       checkID,
			TLSSkipVerify: true,
			TTL:           ttl.String(),
		},
	}

	err = client.Agent().ServiceRegister(registration)
	go updateHealthCheck(client)
	if err != nil {
		log.Fatalf("Failed to register user service: %v", err)
	}

}

func updateHealthCheck(client *api.Client) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		err := client.Agent().UpdateTTL(checkID, "active", api.HealthPassing)
		if err != nil {
			log.Fatal(err)
		}
		<-ticker.C
	}
}

func startRiderService(){
	grpcServer := grpc.NewServer()
	servicePort, err := strconv.Atoi(servicePort)
	if err != nil {
		log.Fatal("Invalid port number")
	}
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	user := config.GetEnv("DB_USER", "ayo")
	pwd := config.GetEnv("DB_PWD", "password")
	dbName := config.GetEnv("DB_NAME", "pickitup_db")
	Db, err := db.ConnectDb(host, port, user, pwd, dbName)
	if err!= nil {
        log.Fatal(err)
    }
	riderRepo := NewRiderRepoImpl(Db)
	NewRiderGrpcHandler(grpcServer, riderRepo)
	log.Printf("Rider Service is running... on port %d\n", int(servicePort))

	if err := grpcServer.Serve(l); err!= nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

func main(){
	registerRiderService()
    startRiderService()
}