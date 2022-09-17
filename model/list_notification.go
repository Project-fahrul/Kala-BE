package model

type ListNotifiation struct {
	Message          string `json:"message"`
	SalesID          int    `json:"sales_id"`
	CustomerID       int    `json:"customer_id"`
	Name             string `json:"customer_name"`
	TypeNotification string `json:"type"`
}
