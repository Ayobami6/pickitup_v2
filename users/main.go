package users

import (
	"log"
	"net"

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
        Port:    50051,
        Check: &api.AgentServiceCheck{
            GRPC:                           "localhost:50051",
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
	l, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()
	Db, err := db.ConnectDb()
	if err!= nil {
        log.Fatal(err)
    }
	userRepo := NewUserRepoImpl(Db)
	NewUsersGrpcHandler(grpcServer, userRepo)
	log.Println("User Service is running... on port 50051")

	if err := grpcServer.Serve(l); err!= nil {
        log.Fatalf("Failed to serve: %v", err)
    }


}

func main() {
	registerUserService()
    startUserService()
}