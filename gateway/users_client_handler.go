package main

import (
	"net/http"

	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserClientHandler struct {
	client pbUser.UserServiceClient
}

func NewUserClientHandler(client pbUser.UserServiceClient) *UserClientHandler {
    return &UserClientHandler{client: client}
}

// Implement register routes
func (h *UserClientHandler) RegisterRoutes(router *mux.Router) {
	// Register your routes here
    router.HandleFunc("/register", h.HandleRegister).Methods("POST")
    router.HandleFunc("/login", h.HandleLoginUser).Methods("POST")
}


func (h *UserClientHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
    var registerPayload *pbUser.UserRegistrationPayload

    err := utils.ParseJSON(r, &registerPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }
    // TODO: Validate Payload with validate
    res, err := h.client.RegisterUser(r.Context(), registerPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Message())
            return
		}
	}
    message := res.Message

    utils.WriteJSON(w, http.StatusCreated, "success", nil, message)
}

func (h *UserClientHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
    var loginPayload *pbUser.UserLoginPayload
    err := utils.ParseJSON(r, &loginPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    res, err := h.client.LoginUser(r.Context(), loginPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			utils.WriteError(w, http.StatusInternalServerError, rStatus.Message())
            return
		}
	}

    utils.WriteJSON(w, http.StatusOK, "success", res, "Login Successful" )
}