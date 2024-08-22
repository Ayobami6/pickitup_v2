package main

type OrderRepo interface {
	CreateOrder(order *Order) error
	GetOrders(userID uint) ([]Order, error)
	UpdateDeliveryStatus(orderId uint, status StatusType) error
	UpdateAcknowledgeStatus(orderId uint) error
	GetOrderByID(orderId uint) (*Order, error)
	CancelOrder(orderId uint) error
}