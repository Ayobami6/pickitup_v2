package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StatusType string

type PaymentStatus string

const (
	Paid   PaymentStatus = "Paid"
	Unpaid PaymentStatus = "Unpaid"
)

const (
	Delivered  StatusType = "Delivered"
	Canceled   StatusType = "Canceled"
	Pending    StatusType = "Pending"
	InDelivery StatusType = "InDelivery"
)

type Order struct {
	ID             uint          `json:"id"`
	RiderID        uint          `json:"rider_id"`
	UserID         uint          `json:"user_id"`
	RefID          string        `json:"ref_id" gorm:"uniqueIndex;size:30"`
	Status         StatusType    `json:"status" gorm:"default:Pending"`
	Item           string        `json:"item"`
	Quantity       int           `json:"quantity"`
	Charge         float64       `json:"price"`
	PaymentStatus  PaymentStatus `json:"payment_status"`
	PickUpAddress  string        `json:"pickup_address"`
	DropOffAddress string        `json:"dropoff_address"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Acknowledge    bool          `json:"acknowledge" gorm:"default:false"`
}

func (u *Order) BeforeSave(tx *gorm.DB) error {
	u.Status = Pending
	if u.Status != Delivered && u.Status != Canceled && u.Status != Pending && u.Status != InDelivery {
		return fmt.Errorf("invalid status")
	}
	u.PaymentStatus = Unpaid
	if u.PaymentStatus != Paid && u.PaymentStatus != Unpaid {
		return fmt.Errorf("invalid payment status")
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	Id, err := generateRandomID(4)
	if err != nil {
		return err
	}
	id := strings.ToUpper(Id)
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, _ := now.Clock()
	o.RefID = fmt.Sprintf("PICK%s%d%d%d%d%d", id, year, month, day, hour, minute)
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = time.Now()
	if o.Acknowledge {
		fmt.Println("I got here")
		o.Status = InDelivery
	}
	return nil
}

func generateRandomID(n int) (string, error) {
    bytes := make([]byte, n)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return strings.ToUpper(hex.EncodeToString(bytes)[:n]), nil
}
