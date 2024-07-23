package main

type RiderRepo interface {
	CreateRider(rider *Rider) error
	GetRiders() ([]Rider, error)
	GetRider(id uint) (*Rider, error)
	GetRiderByUserID(userID uint) (*Rider, error)
	GetRiderByID(id uint) (*Rider, error)
	UpdateRating(riderId uint) error
	GetRiderReviews(riderID uint) ([]Review, error)
	UpdateMinAndMaxCharge(minCharge float64, maxCharge float64, userID uint) error
	UpdateRiderAvailability(riderID uint, status string) error
	UpdateRiderSuccessfulRides(riderID uint) error
}