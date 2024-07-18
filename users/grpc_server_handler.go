package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Ayobami6/common/auth"
	userPb "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"google.golang.org/grpc"
)

type usersGrpcHandler struct {
	userPb.UnimplementedUserServiceServer
	repo UserRepo
}

func NewUsersGrpcHandler(grpcServer *grpc.Server, repo UserRepo) {
	handler := &usersGrpcHandler{repo: repo}
    userPb.RegisterUserServiceServer(grpcServer, handler)
}



func (h *usersGrpcHandler) RegisterUser(ctx context.Context, req *userPb.UserRegistrationPayload) (*userPb.RegisterMessage, error) {
	// user registration logic
	email := req.Email
	password := req.Password
	username := req.Username
	phone := req.PhoneNumber
	// get 
	_, err := h.repo.GetUserByEmail(email)
    if err != nil {
        return nil, errors.New("User with this email already exists")
    }
	// hash password 
	password, err = auth.HashPassword(password)
	if err!= nil {
        return nil, errors.New("Something went wrong")
    }

	if err := h.repo.CreateUser(&User{
		Email:    email,
        Password: password,
        UserName: username,
        PhoneNumber: phone,
	}); err != nil {
		err := err.Error()
		if strings.Contains(err, "uni_users_phone_number") {
			return nil, errors.New("User with this phone number already exists")
		}
		return nil, errors.New("Something went wrong")
	}
	num, err := utils.GenerateAndCacheVerificationCode(email)
	if err!= nil {
        log.Println("Generate Code Failed: ", err)
    } else {
		// send the email to verify
		msg := fmt.Sprintf("Your verification code is %d\n", num)
		err = utils.SendMail(email, "Email Verification", username, msg)
        if err!= nil {
            log.Printf("Email sending failed due to %v\n", err)
        }
	}
	message := &userPb.RegisterMessage{
		Message: "User successfully created!",
	}

	return message, nil

}