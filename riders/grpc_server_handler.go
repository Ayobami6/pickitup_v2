package main

import (
	"context"

	riderPb "github.com/Ayobami6/common/proto/riders"
	"google.golang.org/grpc"
)



type riderGrpcHandler struct {
	riderPb.UnimplementedRiderServiceServer
	repo RiderRepo
}

func NewRiderGrpcHandler(grpcServer *grpc.Server, repo RiderRepo){
	handler := &riderGrpcHandler{repo: repo}
	riderPb.RegisterRiderServiceServer(grpcServer, handler)
}


func (h *riderGrpcHandler)CreateRider(ctx context.Context, body *riderPb.CreateRiderPayload) (*riderPb.CreateRiderResponse, error){
	rider := Rider{
		UserID: uint(body.UserId),
		FirstName: body.FirstName,
		LastName: body.LastName,
		Address: body.Address,
		NextOfKinName: body.NextOfKinName,
		NextOfKinPhone: body.NextOfKinPhone,
		NextOfKinAddress: body.NextOfKinAddress,
		DriverLicenseNumber: body.DriverLicenseNumber,
		BikeNumber: body.BikeNumber,
	}
	
	err := h.repo.CreateRider(&rider)
	if err != nil {
		return nil, err
	}

	response := &riderPb.CreateRiderResponse{
		Message: "Rider Successfully Created!",
	}
	return response, nil

}

