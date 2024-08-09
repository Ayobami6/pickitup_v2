package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

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

func (h *riderGrpcHandler)GetRiderByID(ctx context.Context, in *riderPb.RiderID) (*riderPb.Rider, error) {
	id := uint(in.RiderId)
    rider, err := h.repo.GetRiderByID(id)
    if err!= nil {
        return nil, err
    }
	// get all reviews 
	reviews, err := h.repo.GetRiderReviews(id)
	var parsedReviews []*riderPb.Review
	if err!= nil {
        return nil, err
    }
	for _, review := range reviews {
		parsedReviews = append(parsedReviews, &riderPb.Review{Rating: float32(review.Rating), Comment: review.Comment, RiderId: int64(review.RiderID)})
	}
	
    return &riderPb.Rider{
        RiderId: strconv.Itoa(int(rider.ID)),
        FirstName: rider.FirstName,
        LastName: rider.LastName,
        Address: rider.Address,
		AvailabilityStatus: string(rider.AvailabilityStatus),
        BikeNumber: rider.BikeNumber,
        Rating: float32(rider.Rating),
		Level: rider.Level,
		SuccessfulRides: strconv.Itoa(int(rider.SuccessfulRides)),
		CurrentLocation: rider.CurrentLocation,
		Reviews: parsedReviews,
		MaximumCharge: float32(rider.MaximumCharge),
		MinimumCharge: float32(rider.MinimumCharge),
		UserId: int64(rider.UserID),

    }, nil
}


func (h *riderGrpcHandler) GetRiderByUserID(ctx context.Context, r *riderPb.RiderUserID) (*riderPb.Rider, error) {
	id := uint(r.UserId)
    rider, err := h.repo.GetRiderByUserID(id)
    if err!= nil {
        return nil, err
    }
	reviews, err := h.repo.GetRiderReviews(rider.ID)
	var parsedReviews []*riderPb.Review
	if err!= nil {
        return nil, err
    }
	for _, review := range reviews {
		parsedReviews = append(parsedReviews, &riderPb.Review{Rating: float32(review.Rating), Comment: review.Comment, RiderId: int64(review.RiderID)})
	}
    return &riderPb.Rider{
        RiderId: strconv.Itoa(int(rider.ID)),
        FirstName: rider.FirstName,
        LastName: rider.LastName,
        Address: rider.Address,
        BikeNumber: rider.BikeNumber,
        Rating: float32(rider.Rating),
        Level: rider.Level,
		Reviews: parsedReviews,
        SuccessfulRides: strconv.Itoa(int(rider.SuccessfulRides)),
        CurrentLocation: rider.CurrentLocation,
        UserId: int64(rider.UserID),
    }, nil
}


func (h *riderGrpcHandler)UpdateRiderSuccessfulRides(ctx context.Context, payload *riderPb.UpdateRiderSuccessfulRidesRequest) (*riderPb.UpdateRiderResponse, error) {
	riderID := payload.RiderId
	err := h.repo.UpdateRiderSuccessfulRides(uint(riderID))
	if err!= nil {
        return nil, err
    }
	return &riderPb.UpdateRiderResponse{}, nil
}

// get riders 
func (h *riderGrpcHandler)GetRiders(ctx context.Context, payload *riderPb.GetRidersRequest)(*riderPb.GetRidersResponse, error ){
	riders, err := h.repo.GetRiders()
    var parsedRiders []*riderPb.Rider
    if err != nil {
		log.Printf("Couldn't fetch riders %v \n", err)
        return nil, err
    }
	for _, rider := range riders {
		parsedRiders = append(parsedRiders, &riderPb.Rider{
			RiderId: strconv.Itoa(int(rider.ID)),
			FirstName: rider.FirstName,
			LastName: rider.LastName,
			Address: rider.Address,
			BikeNumber: rider.BikeNumber,
			Rating: float32(rider.Rating),
			Level: rider.Level,
			SuccessfulRides: strconv.Itoa(int(rider.SuccessfulRides)),
			CurrentLocation: rider.CurrentLocation,
			UserId: int64(rider.UserID),
		})
	}
    return &riderPb.GetRidersResponse{Riders: parsedRiders}, nil

}

// update rider charge

func(h *riderGrpcHandler)UpdateMinAndMaxCharge(ctx context.Context, payload *riderPb.ChargeUpdatePayload) (*riderPb.ResponseMessage, error){
	minCharge := float64(payload.MinimumCharge)
    maxCharge := float64(payload.MaximumCharge)
	userId := payload.UserId
	fmt.Println("I got here", minCharge, maxCharge, userId)
    err := h.repo.UpdateMinAndMaxCharge(minCharge, maxCharge, uint(userId))
    if err!= nil {
		log.Printf("Couldn't update rider's charge range %v \n", err.Error())
        return nil, err
    }

    return &riderPb.ResponseMessage{Message: "Rider's charge range successfully updated!"}, nil
}