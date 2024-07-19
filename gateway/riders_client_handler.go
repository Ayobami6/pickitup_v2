package main

import (
	"net/http"

	pbRider "github.com/Ayobami6/common/proto/riders"
	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"github.com/gorilla/mux"
)

type RiderClientHandler struct {
	client pbRider.RiderServiceClient
	userClient pbUser.UserServiceClient
}

func NewRiderClientHandler(client pbRider.RiderServiceClient) *RiderClientHandler {
    return &RiderClientHandler{client: client}
}

func (h *RiderClientHandler)RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register/rider", h.HandleRiderRegister).Methods("POST")
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

	res, err := h.userClient.RegisterUser(r.Context(), &pbUser.UserRegistrationPayload{
		Email: email,
		Username: username,
        Password:  password,
        PhoneNumber: phone_number,
	})

	// TODO: handle some error types
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

	// create rider
	userId := res.UserID

	message, newErr := h.client.CreateRider(r.Context(), &pbRider.CreateRiderPayload{
		UserId: userId,
		FirstName: riderRegisterPayload.FirstName,
		LastName: riderRegisterPayload.LastName,
		Address: riderRegisterPayload.Address,
		NextOfKinName: riderRegisterPayload.NextOfKinName,
		NextOfKinPhone: riderRegisterPayload.NextOfKinPhone,
		NextOfKinAddress: riderRegisterPayload.NextOfKinAddress,
		DriverLicenseNumber: riderRegisterPayload.DriverLicenseNumber,
		BikeNumber: riderRegisterPayload.BikeNumber,
	})

	if newErr != nil {
		utils.WriteError(w, http.StatusInternalServerError, newErr.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "success", nil, message.Message)

}