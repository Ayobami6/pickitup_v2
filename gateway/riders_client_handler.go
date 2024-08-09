package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Ayobami6/common/auth"
	pbRider "github.com/Ayobami6/common/proto/riders"
	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"github.com/Ayobami6/gateway/dto"
	"github.com/gorilla/mux"
)

type RiderClientHandler struct {
	client pbRider.RiderServiceClient
	userClient pbUser.UserServiceClient
}

func NewRiderClientHandler(client pbRider.RiderServiceClient, userClient pbUser.UserServiceClient) *RiderClientHandler {
    return &RiderClientHandler{client: client, userClient: userClient}
}

func (h *RiderClientHandler)RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register/rider", h.HandleRiderRegister).Methods("POST")
	router.HandleFunc("/riders/{id}", h.HandleGetRider).Methods("GET")
	router.HandleFunc("/riders", h.HandleGetRiders).Methods("GET")
	router.HandleFunc("/riders/{id}/charge-update",  auth.RiderAuth(h.HandleUpdateRiderCharge, h.client)).Methods("PUT")
	router.HandleFunc("/riders/{id}/status-update",  auth.RiderAuth(h.HandleUpdateStatus, h.client)).Methods("PATCH")
}

// handler rider register

func (h *RiderClientHandler) HandleRiderRegister(w http.ResponseWriter, r *http.Request){
	var riderRegisterPayload pbRider.CreateRiderPayload

	err := utils.ParseJSON(r, &riderRegisterPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	// validate payload
	email := riderRegisterPayload.Email
	username := riderRegisterPayload.Username
	password := riderRegisterPayload.Password
	phone_number := riderRegisterPayload.PhoneNumber
	// create user 
	res, err := h.userClient.RegisterUser(r.Context(), &pbUser.UserRegistrationPayload{
        Email:    email,
        Password: password,
        Username: username,
        PhoneNumber: phone_number,
    })
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

	// TODO: handle some error types

	// create rider
	userID := res.UserID

	message, newErr := h.client.CreateRider(r.Context(), &pbRider.CreateRiderPayload{
		UserId: int64(userID),
		FirstName: riderRegisterPayload.FirstName,
		LastName: riderRegisterPayload.LastName,
		Address: riderRegisterPayload.Address,
		NextOfKinName: riderRegisterPayload.NextOfKinName,
		NextOfKinPhone: riderRegisterPayload.NextOfKinPhone,
		NextOfKinAddress: riderRegisterPayload.NextOfKinAddress,
		DriverLicenseNumber: riderRegisterPayload.DriverLicenseNumber,
		BikeNumber: riderRegisterPayload.BikeNumber,
		Email: email,
		Password: password,
		PhoneNumber: phone_number,
		Username: username,
	})
	log.Println("Couldnt get here")

	if newErr != nil {
		utils.WriteError(w, http.StatusInternalServerError, newErr.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "success", nil, message.Message)

}

func (h *RiderClientHandler) HandleGetRider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ID, err := strconv.Atoi(id)
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
	rider, err := h.client.GetRiderByID(r.Context(), &pbRider.RiderID{RiderId: int64(ID)})
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
    }
	domain := getDomainURL(r)
	var selfUrl = fmt.Sprintf("%s/api/v2/riders/%s", domain, rider.RiderId)
	rider.SelfUrl = selfUrl

    utils.WriteJSON(w, http.StatusOK, "success", rider, "Fetch Successful")
}


func (h *RiderClientHandler)HandleGetRiders(w http.ResponseWriter, r *http.Request) {
	riders, err := h.client.GetRiders(r.Context(), &pbRider.GetRidersRequest{})
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Something went wrong")
		return
    }
	domain := getDomainURL(r)
	for _, r := range riders.Riders {
        var selfUrl = fmt.Sprintf("%s/api/v2/riders/%s", domain, r.RiderId)
        r.SelfUrl = selfUrl
    }
    utils.WriteJSON(w, http.StatusOK, "success", riders, "Fetch Successful")
    // TODO: handle pagination and sorting
    // TODO: add authentication to fetch riders by user ID
    // TODO: add rate limiting to prevent abuse
}

func(h *RiderClientHandler)HandleUpdateRiderCharge(w http.ResponseWriter, r *http.Request){
	// get user Id from auth context
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	fmt.Println(id)
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
	ctx := r.Context()
	userId := auth.GetUserIDFromContext(ctx)
	fmt.Println("This is user userId: ", userId)
	fmt.Println("This is id: ", id)
	if userId == -1 {
		fmt.Println("Got here")
        auth.Forbidden(w)
        return
    }
	if int64(id) != int64(userId) {
		auth.Forbidden(w)
        return
	}
	var payload dto.UpdateChargeDTO
	err = utils.ParseJSON(r, &payload)
	if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Payload")
        return
    }
	fmt.Println(payload)
	response, err := h.client.UpdateMinAndMaxCharge(ctx, &pbRider.ChargeUpdatePayload{
		MaximumCharge: float32(payload.MaximumCharge),
		MinimumCharge: float32(payload.MinimumCharge),
		UserId: int64(userId),
	})
	if err!= nil {
		log.Println(err)
        utils.WriteError(w, http.StatusInternalServerError, "Failed to update charges")
        return
    }
	utils.WriteJSON(w, http.StatusOK, "success", response, "Charge updated successfully")

}

func(h *RiderClientHandler)HandleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid ID")
        return
    }
    ctx := r.Context()
    riderId := auth.GetRiderIDFromContext(ctx)
    if riderId == -1 {
        auth.Forbidden(w)
        return
    }
    if int64(id)!= int64(riderId) {
        auth.Forbidden(w)
        return
    }
    var payload dto.UpdateStatusDTO
    err = utils.ParseJSON(r, &payload)
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid Payload")
        return
    }
	// TODO Add status type validation
	statusMap := map[string]bool{
		"Available": true,
		"Unavailable": true,
		"On Break": true,
		"Busy": true,
	}
	if !statusMap[payload.Status] {
		utils.WriteError(w, http.StatusBadRequest, "Invalid availability status")
		return
	}
	// update status
	res, err := h.client.UpdateAvailabilityStatus(ctx, &pbRider.UpdateAvailabiltyStatusPayLoad{
		RiderId: int64(riderId),
		Status: payload.Status,
	})
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to update status")
        return
    }
	utils.WriteJSON(w, http.StatusOK, "success", res, "Status updated successfully")
}

func getDomainURL(r *http.Request) string {
    scheme := "http"
    if r.TLS != nil {
        scheme = "https"
    }
    return scheme + "://" + r.Host
}