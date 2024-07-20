package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Ayobami6/common/auth"
	"github.com/Ayobami6/common/config"
	userPb "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type usersGrpcHandler struct {
	userPb.UnimplementedUserServiceServer
	repo UserRepo
	db *gorm.DB
}

func NewUsersGrpcHandler(grpcServer *grpc.Server, repo UserRepo, db *gorm.DB) {
	handler := &usersGrpcHandler{repo: repo, db: db}
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
    if err == nil {
        return nil, errors.New("user with this email already exists")
    }
	// hash password 
	password, err = auth.HashPassword(password)
	if err!= nil {
        return nil, errors.New("something went wrong")
    }
	user := &User{
		Email:    email,
        Password: password,
        UserName: username,
        PhoneNumber: phone,
	}

	if err := h.repo.CreateUser(user); err != nil {
		err := err.Error()
		if strings.Contains(err, "uni_users_phone_number") {
			return nil, errors.New("user with this phone number already exists")
		}
		return nil, errors.New("something went wrong")
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
		UserID: int64(user.ID),
	}

	return message, nil

}

func (h *usersGrpcHandler) LoginUser(ctx context.Context, in *userPb.UserLoginPayload) (*userPb.LoginResponse, error) {
	email := in.Email
    password := in.Password
    user, err := h.repo.GetUserByEmail(email)
    if err!= nil {
        return nil, errors.New("user not found")
    }
    if!auth.CheckPassword(user.Password, []byte(password)) {
        return nil, errors.New("incorrect password")
    }
	secret := []byte(config.GetEnv("JWT_SECRET", "secret"))
    token, err := auth.CreateJWT(secret, int(user.ID))
    if err!= nil {
        return nil, errors.New("failed to generate token")
    }
    message := &userPb.LoginResponse{
        AccessToken: token,
    }
    return message, nil
}

func (h *usersGrpcHandler)GetUserByID(ctx context.Context, in *userPb.UserIDMessage) (*userPb.User, error) {
	id := in.UserId
    user, err := h.repo.GetUserByID(uint(id))
    if err!= nil {
        return &userPb.User{}, errors.New("user not found")
    }
    return &userPb.User{
        Id:          int64(user.ID),
        Email:        user.Email,
        Username:     user.UserName,
        PhoneNumber: user.PhoneNumber,
		WalletBalance: float32(user.WalletBalance),
		RiderId: int64(user.RiderID),
		Verified: user.Verified,
		AccountName: user.AccountName,
		BankName: user.BankName,
		AccountNumber: user.AccountNumber,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
        UpdatedAt: user.UpdatedAt.Format(time.RFC3339),	
    }, nil
}

func (h *usersGrpcHandler)ChargeUserWallet(ctx context.Context, in *userPb.ChargeRequest) (*userPb.ChargeResponse, error){
	userId := in.UserId
	charge := in.Charge
	// get the user 
	user, err := h.repo.GetUserByID(uint(userId))
    if err!= nil {
        return nil, errors.New("user not found")
    }
	err = user.Debit(h.db, float64(charge))
	if err!= nil {
        return nil, errors.New("failed to debit user wallet")
    }
	message := &userPb.ChargeResponse{}

	return message, nil
    
}

