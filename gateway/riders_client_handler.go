package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	pbRider "github.com/Ayobami6/common/proto/riders"
	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
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

func getDomainURL(r *http.Request) string {
    scheme := "http"
    if r.TLS != nil {
        scheme = "https"
    }
    return scheme + "://" + r.Host
}