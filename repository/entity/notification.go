package entity

import "time"

type Notifications struct {
	SalesID          int       `json:"sales_id"`
	CustomerID       int       `json:"customer_id"`
	Message          string    `json:"message"`
	TypeNotification string    `json:"type_notification"`
	DueDate          time.Time `json:"due_date"`
}
