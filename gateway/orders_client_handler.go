package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ayobami6/common/auth"
	pb "github.com/Ayobami6/common/proto/orders"
	pbRider "github.com/Ayobami6/common/proto/riders"
	pbUser "github.com/Ayobami6/common/proto/users"
	"github.com/Ayobami6/common/utils"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)


type OrderClientHandler struct {
	client pb.OrderServiceClient
	userClient pbUser.UserServiceClient
	riderClient pbRider.RiderServiceClient
	ch *amqp.Channel
}

func NewOrderClientHandler(client pb.OrderServiceClient, userClient pbUser.UserServiceClient, riderClient pbRider.RiderServiceClient, ch *amqp.Channel) *OrderClientHandler {
    return &OrderClientHandler{client: client, userClient: userClient, riderClient: riderClient, ch: ch}
}

func (h *OrderClientHandler)RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders/{rider_id}", auth.Auth(h.HandleCreateOrder, h.userClient)).Methods("POST")
	router.HandleFunc("/orders", auth.Auth(h.HandleGetOrders, h.userClient)).Methods("GET")
	router.HandleFunc("/orders/{id}", auth.Auth(h.HandleGetOrder, h.userClient)).Methods("GET")
}

func (h *OrderClientHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var createOrderPayload pb.CreateOrderRequest
	if err := utils.ParseJSON(r, &createOrderPayload); err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Bad Data!")
        return
    }
	log.Println(createOrderPayload.PickupAddress, createOrderPayload.DropOffAddress)

	params := mux.Vars(r)
    riderID, err := strconv.Atoi(params["rider_id"])
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid rider id")
        return
    }

	var riderId uint = uint(riderID)
	rider, err := h.riderClient.GetRiderByID(r.Context(), &pbRider.RiderID{RiderId: int64(riderId)})
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "It's Us!")
        return
    }
	riderAvailableStatus := rider.AvailabilityStatus
	switch {
		case riderAvailableStatus == "Unavailable":
        case riderAvailableStatus == "OnBreak":
		case riderAvailableStatus == "Busy":
            utils.WriteError(w, http.StatusBadRequest, "Rider is currently unavailable")
            return
        default:
			break  
	}
	minCharge := rider.MinimumCharge
	maxCharge := rider.MaximumCharge
	charge := minCharge + ((maxCharge - minCharge)/2)

	ctx := r.Context()
	// get user Id from context
	userID := auth.GetUserIDFromContext(ctx)
	if userID == -1 {
        auth.Forbidden(w)
        return
    }

	// contact user client
	user, err := h.userClient.GetUserByID(ctx, &pbUser.UserIDMessage{
		UserId: int64(userID),
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get user")
        return
    }

	if !user.Verified {
		utils.WriteError(w, http.StatusBadRequest, "User is not verified")
        return
	}
	// check user balance
	if bal := user.WalletBalance; bal < charge {
		utils.WriteError(w, http.StatusBadRequest, "Insufficient balance")
        return
	}
	// charge user
	_, cErr := h.userClient.ChargeUserWallet(ctx, &pbUser.ChargeRequest{
		UserId: int64(userID),
		Charge: charge,
	})
	if cErr!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to deduct balance")
        return
    }

	res, oErr := h.client.CreateOrder(ctx, &pb.CreateOrderRequest{
		RiderId: int64(riderId),
        UserId: int64(userID),
        Charge: float64(charge),
        Item: createOrderPayload.Item,
		Quantity: createOrderPayload.Quantity,
		PickupAddress: createOrderPayload.PickupAddress,
        DropOffAddress: createOrderPayload.DropOffAddress,

	})

	if oErr!= nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create order")
        return
    }
	riderUser, rErr := h.userClient.GetUserByID(ctx, &pbUser.UserIDMessage{
		UserId: int64(rider.UserId),
	})
	if rErr!= nil {
        log.Println("Error getting rider user")
    }
	respData := map[string]string{
		"ref_id": res.RefId,
	}
	log.Println(res)
	riderMessage := fmt.Sprintf("You have New Pick Up Order with ID %s\n Containing item %s which is to be picked up at %s \n and delivered at %s Please go to your dashboard to accept the order and transit immediately or reject \n", res.RefId, res.Item, res.PickupAddress, res.DropOffAddress)
	userMessage := fmt.Sprintf("Your Order %s has been placed successfully \n Here is your rider phone number %s\n\n", res.RefId, riderUser.PhoneNumber)
	subject := "PickItUp Order Notification"
	mailData := map[string]string{
		"riderName": riderUser.Username,
        "orderID": res.RefId,
		"subject":subject,
		"riderMessage": riderMessage,
		"userMessage": userMessage,
		"userEmail": user.Email,
		"riderEmail": riderUser.Email,
		"userUsername": user.Username,
	}
	// start rabbit channel
	ch := h.ch
	// declare queue
	q, err := ch.QueueDeclare(
		"order_notification",
		false,   
		false,  
		false,
		false,
		nil, 
	  )

	if err != nil {
		log.Println("Failed to declare a queue", err)
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

	// send mail with notification service using message broker

	utils.WriteJSON(w, http.StatusCreated, "success", respData, "Order Created")

}

func (h *OrderClientHandler)HandleGetOrders(w http.ResponseWriter, r *http.Request) {
	// get user from request
	ctx := r.Context()
    userID := auth.GetUserIDFromContext(ctx)
    if userID == -1 {
        auth.Forbidden(w)
        return
    }
	// get orders
	orders, err := h.client.GetOrders(ctx, &pb.AllOderRequest{
		UserId: int64(userID),
	})
	if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get orders")
        return
    }
	utils.WriteJSON(w, http.StatusOK, "success", orders, "Orders retrieved successfully")
}

func (h *OrderClientHandler)HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	// get order id for param
	params := mux.Vars(r)
	ctx := r.Context()
    orderId, err := strconv.Atoi(params["id"])
    if err!= nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid order id")
        return
    }
    // get order
    order, err := h.client.GetOrder(ctx, &pb.GetOrderRequest{
        Id: int64(orderId),
    })
    if err!= nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to get order")
        return
    }
    utils.WriteJSON(w, http.StatusOK, "success", order, "Order retrieved successfully")
}

// TODO: Implement UpdateAcknowledgeStatus - restricted to only Riders
// TODO: UpdateDeliveryStatus - restricted to only users

