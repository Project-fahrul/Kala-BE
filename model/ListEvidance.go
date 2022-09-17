package model

import "time"

type ListEvidance struct {
	Name         string    `json:"name"`
	TypeEvidance string    `json:"type"`
	SalesName    string    `json:"sales_name"`
	DueDate      time.Time `json:"due"`
	SubmitDate   bool      `json:"submit"`
	SalesID      int       `json:"sales_id"`
	CustomerID   int       `json:"customer_id"`
}
