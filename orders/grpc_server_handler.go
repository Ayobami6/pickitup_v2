package main

import (
	"context"
	"log"

	pb "github.com/Ayobami6/common/proto/orders"
	"google.golang.org/grpc"
)

type orderGrpcHandler struct {
	pb.UnimplementedOrderServiceServer
	repo OrderRepo
}

func NewOrderGrpcHandler(grpcServer *grpc.Server, repo OrderRepo) {
    handler := &orderGrpcHandler{repo: repo}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}


func (h *orderGrpcHandler) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest)(*pb.OrderResponse, error) {
	order := Order{
		RiderID: uint(in.RiderId),
		UserID: uint(in.UserId),
		Quantity: int(in.Quantity),
		Charge: in.Charge,
		Item: in.Item,
		PickUpAddress: in.PickupAddress,
		DropOffAddress: in.DropOffAddress,
	}
	err := h.repo.CreateOrder(&order)
	if err!= nil {
        log.Println("Error creating order: ", err)
        return nil, err
    }
	response := &pb.OrderResponse{
		RefId: order.RefID,
		Status: string(order.Status),
		CreatedAt: order.CreatedAt.String(),
		Charge: order.Charge,
        Item: order.Item,
		Acknowledge: order.Acknowledge,
		RiderId: int64(order.RiderID),
		UserId: int64(order.UserID),
		PaymentStatus: string(order.PaymentStatus),
		PickupAddress: order.PickUpAddress,
		DropOffAddress: order.DropOffAddress,

	}
	return response, nil
    
}