package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Ayobami6/common/auth"
	"github.com/Ayobami6/common/config"
	userPb "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type usersGrpcHandler struct {
	userPb.UnimplementedUserServiceServer
	repo UserRepo
	db *gorm.DB
	ch *amqp.Channel
}

func NewUsersGrpcHandler(grpcServer *grpc.Server, repo UserRepo, db *gorm.DB, ch *amqp.Channel) {
	handler := &usersGrpcHandler{repo: repo, db: db, ch: ch}
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
		ch := h.ch
		// declare email verification queue
		q, err := ch.QueueDeclare(
            "email_verification",         
            false,               
            false,               
            false,               
            false,               
            nil,                 
        )
        if err!= nil {
            log.Fatalf("Failed to declare queue: %v", err)
        }
		mailData := map[string]string{
			"username": username,
            "message": msg,
			"email": email,
			"subject": "Email Verification",
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		body, err := json.Marshal(mailData)
		if err != nil {
			log.Println("Failed to marshal data", err)
		}
		err = ch.PublishWithContext(ctx,
			"",     
			q.Name,
			false,
			false,
			amqp.Publishing {
			  ContentType: "application/json",
			  Body:        body,
			})
		if err!= nil {
			log.Println("Failed to publish a message", err)
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

func (h *usersGrpcHandler)CreditUserWallet(ctx context.Context, payload *userPb.ChargeRequest)(*userPb.ChargeResponse, error ) {
	userId := payload.UserId
    charge := payload.Charge
    // get the user 
    user, err := h.repo.GetUserByID(uint(userId))
    if err!= nil {
        return nil, errors.New("user not found")
    }
    err = user.Credit(h.db, float64(charge))
    if err!= nil {
		log.Printf("Failed to credit user wallet: %v", err)
        return nil, errors.New("failed to credit user wallet")
    }
    message := &userPb.ChargeResponse{}

    return message, nil
}

// #8
func (h *usersGrpcHandler)VerifyOTP(ctx context.Context, payload *userPb.OTPVerifyPayload) (*userPb.OTPVerifyResponse, error) {
	email := payload.Email
    otp := payload.Otp
    user, err := h.repo.GetUserByEmail(email)
    if err!= nil {
        return nil, errors.New("user not found")
    }
//    get cached otp 
    cachedOtp, err := utils.GetCachedVerificationCode(email)
    if err!= nil {
        return nil, errors.New("failed to get cached otp")
    }
	castedOtp := strconv.Itoa(cachedOtp)
	log.Printf("Cached OTP: %s, Provided OTP: %s \n", castedOtp, otp)
    if castedOtp != otp {
        return nil, errors.New("invalid otp")
    }
	// update verification status
	user.Verified = true
    err = h.repo.UpdateUser(user)
	if err!= nil {
        return nil, errors.New("failed to update user verification status")
    }
	// credit user wallet #1000
	err = user.Credit(h.db, float64(1000))
    if err!= nil {
        log.Println(err.Error())
    }

    message := &userPb.OTPVerifyResponse{
        Message: "User verified successfully!",
    }

    return message, nil
}