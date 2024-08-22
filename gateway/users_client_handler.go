package main

import (
	"log"
	"net/http"

	"github.com/Ayobami6/common/auth"
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
    router.HandleFunc("/users/details", auth.Auth(h.HandleGetUserDetails, h.client)).Methods("GET")
    router.HandleFunc("/users/otp/verify", h.HandleVerifyOTP).Methods("POST")
    router.HandleFunc("/users/otp/resend", h.ResendOTP).Methods("POST")
    router.HandleFunc("/users/wallet/balance", auth.Auth(h.GetWalletBalance, h.client)).Methods("GET")
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

func (h *UserClientHandler) HandleGetUserDetails(w http.ResponseWriter, r *http.Request){
    userID := auth.GetUserIDFromContext(r.Context())
    log.Println(userID)
    if userID == -1 {
        auth.Forbidden(w)
        return
    }
    res, err := h.client.GetUserByID(r.Context(), &pbUser.UserIDMessage{
        UserId: int64(userID),
    })
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get user details")
        return
    }
    utils.WriteJSON(w, http.StatusOK, "success", res, "User Details Retrieved Successfully")
}

func (h *UserClientHandler)HandleVerifyOTP(w http.ResponseWriter, r *http.Request){
    // 
    var verifyPayload pbUser.OTPVerifyPayload
    err := utils.ParseJSON(r, &verifyPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }
    res, err := h.client.VerifyOTP(r.Context(), &verifyPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }
    utils.WriteJSON(w, http.StatusOK, "success", res, "OTP verified successfully")
}

func (h *UserClientHandler)ResendOTP(w http.ResponseWriter, r *http.Request) {
    var OTPResendPayload pbUser.OTPResendPayload
    err := utils.ParseJSON(r, &OTPResendPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }
    res, err := h.client.ResendOTP(r.Context(), &OTPResendPayload)
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }
    utils.WriteJSON(w, http.StatusOK, "success", res, "OTP resent successfully")
}


func (h *UserClientHandler)GetWalletBalance(w http.ResponseWriter, r *http.Request){
    // get user id from request context
    userID := auth.GetUserIDFromContext(r.Context())
    if userID == -1 {
        auth.Forbidden(w)
        return
    }
    res, err := h.client.GetWalletBalance(r.Context(), &pbUser.WalletBalanceRequest{
        UserId: int64(userID),
    })
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get wallet balance")
        return
    }
    utils.WriteJSON(w, http.StatusOK, "success", res, "Wallet Balance Retrieved Successfully")
}