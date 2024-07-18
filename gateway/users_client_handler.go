package main

import (
	"net/http"

	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/gorilla/mux"
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
}


func (h *UserClientHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	//... implement user registration logic here
    // Create a new user
    // Call the UserServiceClient.Register method
    // Write the response to the client
    // Handle errors as needed
    //...
    // Example:
    // user := &pbUser.User{
    //     Email:    r.FormValue("email"),
    //     Password: r.FormValue("password"),
    // }
    //
    // resp, err := h.client.Register(context.Background(), user)
    // if err!= nil {
    //     http.Error(w, err.Error(), http.StatusInternalServerError)
    //     return
    // }
    //
    // w.Header().Set("Content-Type", "application/json")
    // json.NewEncoder(w).Encode(resp)
}