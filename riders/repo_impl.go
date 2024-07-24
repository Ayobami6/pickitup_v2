package main

import (
	"log"
	"math"

	"gorm.io/gorm"
)

type RiderRepoImpl struct {
	db *gorm.DB
}

func NewRiderRepoImpl(db *gorm.DB) *RiderRepoImpl {
	err := db.AutoMigrate(&Rider{}, &Review{})
    if err!= nil {
        log.Fatal(err)
    }
    return &RiderRepoImpl{db: db}
}

func (r *RiderRepoImpl) CreateRider(rider *Rider) error {
	return r.db.Create(rider).Error
}

func (r *RiderRepoImpl) GetRiders() ([]Rider, error) {
	riders := []Rider{}
	res := r.db.Find(&riders)
	if res.Error != nil {
		return nil, res.Error
	}
	return riders, nil
}

func (r *RiderRepoImpl) GetRider(id uint) (*Rider, error) {
	var rider Rider
	res := r.db.First(&rider, id)
	if res.Error != nil {
		return nil, res.Error
	}
	// get reviews
	reviews, err := r.GetRiderReviews(id)
	if err != nil {
		log.Println(err)
	}
	rider.Reviews = reviews
	return &rider, nil
}

func (r *RiderRepoImpl) GetRiderReviews(riderID uint) ([]Review, error) {
	var reviews []Review
    res := r.db.Where(Review{RiderID: riderID}).Find(&reviews)
    if res.Error!= nil {
        return nil, res.Error
    }
    return reviews, nil
}

func (r *RiderRepoImpl)GetRiderByUserID(userID uint) (*Rider,error){
	var rider Rider
    res := r.db.Where(Rider{UserID: userID}).First(&rider)
    if res.Error!= nil {
        return nil, res.Error
    }
    return &rider, nil
}


func (r *RiderRepoImpl) GetRiderByID(id uint)(*Rider, error) {
	var rider Rider
    res := r.db.First(&rider, id)
    if res.Error!= nil {
        return &Rider{}, res.Error
    }
    return &rider, nil
}


func (r *RiderRepoImpl)UpdateRating(riderID uint) error {
	var rider Rider
	res := r.db.Where(Rider{UserID: riderID}).First(&rider)
	if res.Error!= nil {
        return res.Error
    }
	// get all reviews
	var reviews []Review
    res = r.db.Where(Review{RiderID: riderID}).Find(&reviews)
    if res.Error!= nil {
        return res.Error
    }
    var totalRating float64 = 0
    for _, review := range reviews {
        totalRating += review.Rating
    }
    newRating := totalRating / float64(len(reviews))
    // round to 1 decimal place
    newRating = math.Round(newRating*10) / 10

    // update rider rating
    rider.Rating = newRating
	res = r.db.Save(&rider)
    if res.Error!= nil {
        return res.Error
    }
    return  nil
}


func (r *RiderRepoImpl)UpdateMinAndMaxCharge(minCharge float64, maxCharge float64, userID uint) error {
	res := r.db.Where(Rider{UserID: userID}).Updates(Rider{MinimumCharge: minCharge, MaximumCharge: maxCharge})
	if res.Error!= nil {
        return res.Error
    }
	return nil
}

func (r *RiderRepoImpl) UpdateRiderAvailability(riderId uint, status string) error {
	res := r.db.Where(Rider{UserID: riderId}).Updates(Rider{AvailabilityStatus: RiderAvailabilityStatus(status)})
    if res.Error!= nil {
        return res.Error
    }
    return nil
}

func (r *RiderRepoImpl)UpdateRiderSuccessfulRides(riderId uint) error {
    rider, err := r.GetRiderByID(riderId)
    if err!= nil {
        log.Printf("Unable to get rider %d\n", riderId)
        return err
    }
    rider.UpdateSuccessfulRides(r.db)
    return nil
}