package main

import (
	"log"
	"net"

	"github.com/Ayobami6/common/config"
	"github.com/Ayobami6/common/db"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

// register service to consul start service

// register consul service

func registerUserService() {
	conf := api.DefaultConfig()
    client, err := api.NewClient(conf)
    if err != nil {
        log.Fatalf("Failed to create Consul client: %v", err)
    }

    registration := &api.AgentServiceRegistration{
        ID:      "user-service-1",
        Name:    "user-service",
        Address: "localhost",
        Port:    5005,
        Check: &api.AgentServiceCheck{
            GRPC:                           "localhost:5005",
            Interval:                       "10s",
            DeregisterCriticalServiceAfter: "1m",
        },
    }

    err = client.Agent().ServiceRegister(registration)
    if err != nil {
        log.Fatalf("Failed to register user service: %v", err)
    }

}

func startUserService() {
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", "localhost:5005")
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
	userRepo := NewUserRepoImpl(Db)
	NewUsersGrpcHandler(grpcServer, userRepo)
	log.Println("User Service is running... on port 5005")

	if err := grpcServer.Serve(l); err!= nil {
        log.Fatalf("Failed to serve: %v", err)
    }


}

func main() {
	registerUserService()
    startUserService()
}