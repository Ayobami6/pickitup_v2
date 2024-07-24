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

// get user orders
func (h * orderGrpcHandler)GetOrders(ctx context.Context, payload *pb.AllOderRequest) (*pb.AllOrderReponse, error) {
	orders, err  := h.repo.GetOrders(uint(payload.UserId))
	if err!= nil {
        log.Println("Error getting orders: ", err)
        return nil, err
    }
	var parsedOrders []*pb.OrderResponse
	for _, order := range orders {
		parsedOrders = append(parsedOrders, &pb.OrderResponse{
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
			Id: int64(order.ID),
        })
	}
	return &pb.AllOrderReponse{Orders: parsedOrders}, nil

}

// get order 
func (h *orderGrpcHandler)GetOrder(ctx context.Context, payload *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := h.repo.GetOrderByID(uint(payload.Id))
    if err!= nil {
        log.Println("Error getting order: ", err)
        return nil, err
    }
    return &pb.OrderResponse{
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
		Id: int64(order.ID),
    }, nil
}

func (h *orderGrpcHandler)UpdateDeliveryStatus(ctx context.Context, payload *pb.UpdateDeliveryStatusRequest) (*pb.UpdateResponse, error){
	orderId := payload.Id
	status := payload.Status

	err := h.repo.UpdateDeliveryStatus(uint(orderId), StatusType(status))
	if err!= nil {
        log.Println("Error updating delivery status: ", err)
        return nil, err
    }
	return &pb.UpdateResponse{Message: "Delivery status updated successfully"}, nil

}

func (h *orderGrpcHandler) UpdateAcknowledgement(ctx context.Context, payload *pb.UpdateAcknowledgementRequest) (*pb.UpdateResponse, error){
	orderId := payload.Id
	err := h.repo.UpdateAcknowledgeStatus(uint(orderId))
	if err!= nil {
        log.Println("Error updating acknowledgement status: ", err)
        return nil, err
    }
	return &pb.UpdateResponse{Message: "Acknowledgement status updated successfully"}, nil

}