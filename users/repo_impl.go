package main

import (
	"log"

	"gorm.io/gorm"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) *UserRepoImpl {
	// automigrate 
	err := db.AutoMigrate(&User{})
    if err!= nil {
        log.Fatal(err)
    }
    return &UserRepoImpl{db: db}
}


func (u *UserRepoImpl) CreateUser(user *User) (error) {
	res := u.db.Create(&user)
	if res.Error!= nil {
        return res.Error
    }
	return nil
}


func (u *UserRepoImpl) GetUserByEmail(email string) (*User, error) {
	result := &User{}
    err := u.db.Where("email =?", email).First(&result).Error
    if err!= nil {
        return nil, err
    }
    return result, nil
}

func (u *UserRepoImpl) GetUserByID(id uint) (*User, error) {
	result := &User{}
    res := u.db.First(&result, id)
    if res.Error!= nil {
        return nil, res.Error
    }
    return result, nil
}