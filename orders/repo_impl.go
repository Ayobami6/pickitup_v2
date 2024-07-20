package main

import (
	"log"

	"gorm.io/gorm"
)

type OrderRepoImpl struct {
	db *gorm.DB
}


func NewOrderRepoImpl(db *gorm.DB) *OrderRepoImpl {
	err := db.AutoMigrate(&Order{})
	if err!= nil {
        log.Println(err)
    }
	return &OrderRepoImpl{db: db}
}

func (o *OrderRepoImpl)CreateOrder(order *Order) (error) {
	return o.db.Create(order).Error
}

func (o *OrderRepoImpl) GetOrders(userID uint)([]Order, error) {
	var orders []Order
    return orders, o.db.Where("user_id =?", userID).Find(&orders).Error
}

func (o *OrderRepoImpl) GetOrderByID(id uint)(*Order, error) {
    var order Order
    return &order, o.db.First(&order, id).Error
}

func (o *OrderRepoImpl) UpdateDeliveryStatus(orderID uint, status StatusType) (error) {
    return o.db.Model(&Order{}).Where("id =?", orderID).Update("status", status).Error
}

func (o *OrderRepoImpl) UpdateAcknowledgeStatus(orderID uint) (error) {
    return o.db.Model(&Order{}).Where("id = ?", orderID).Update("acknowledge", true).Error
}