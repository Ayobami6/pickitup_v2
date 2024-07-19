package main

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"gorm.io/gorm"
)

type RiderAvailabilityStatus string

const (
	Available   RiderAvailabilityStatus = "Available"
	Unavailable RiderAvailabilityStatus = "Unavailable"
	OnBreak     RiderAvailabilityStatus = "On Break"
	Busy        RiderAvailabilityStatus = "Busy"
)

type Rider struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	RiderID             string    `json:"rider_id" gorm:"uniqueIndex;size:8"`
	BikeNumber          string    `json:"bike_number"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	UserID              uint      `json:"user_id"`
	NextOfKinName       string    `json:"next_of_kin_name"`
	NextOfKinPhone      string    `json:"next_of_kin_phone"`
	DriverLicenseNumber string    `json:"driver_license_number"`
	NextOfKinAddress    string    `json:"next_of_kin_address"`
	Address             string    `json:"address"`
	SuccessfulRides     int64     `json:"successful_rides"`
	Rating              float64   `json:"rating"`
	Level               string    `json:"level"`
	CurrentLocation     string    `json:"current_location"`
	Reviews             []Review  `json:"reviews" gorm:"foreignKey:RiderID"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	// add minimum and maximum charge to this
	MinimumCharge      float64                 `json:"minimum_charge"`
	MaximumCharge      float64                 `json:"maximum_charge"`
	AvailabilityStatus RiderAvailabilityStatus `json:"availability_status" gorm:"default:Available"`
}

func (u *Rider) BeforeCreate(tx *gorm.DB) (err error) {
	u.RiderID, err = generateRandomID(8)
	if err != nil {
		return err
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (r *Rider) UpdateSuccessfulRides(tx *gorm.DB) error {
	r.SuccessfulRides++
	return tx.Save(r).Error
}

func (r *Rider) BeforeUpdate(tx *gorm.DB) (err error) {
	r.UpdatedAt = time.Now()
	return nil
}

type Review struct {
	ID      uint    `json:"id"`
	RiderID uint    `json:"rider_id"`
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}

func generateRandomID(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)[:n]), nil
}
