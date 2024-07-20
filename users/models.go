package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserName      string    `json:"username"`
	Email         string    `json:"email required" gorm:"unique"`
	PhoneNumber   string    `json:"phone_number" gorm:"unique"`
	WalletBalance float64   `json:"wallet_balance"`
	AccountNumber string    `json:"account_number" gorm:""`
	AccountName   string    `json:"account_name"`
	BankName      string    `json:"bank_name"`
	Password      string    `json:"password"`
	RiderID  	uint `json:"rider_id"`
	Verified      bool      `json:"verified" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (u *User) Debit(tx *gorm.DB, amount float64) error {
	// Debit user wallet balance and save user again
	u.WalletBalance = u.WalletBalance - amount
	return tx.Save(u).Error

}

func (u *User) Credit(tx *gorm.DB, amount float64) error {
	// fund user wallet
	u.WalletBalance = u.WalletBalance + amount
	return tx.Save(u).Error
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    u.UpdatedAt = time.Now()
    return nil
}