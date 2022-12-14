package entity

import "time"

type Evidances struct {
	SalesID      int       `json:"sales_id"`
	CustomerID   int       `json:"customer_id"`
	SubmitDate   time.Time `json:"submit_date"`
	DueDate      time.Time `json:"due_date"`
	Content      string    `json:"content"`
	Comment      string    `json:"comment"`
	TypeEvidance string    `json:"type_evidance"`
}

type EvidancesUpdate struct {
	SubmitDate time.Time `json:"submit_date"`
	Content    string    `json:"content"`
	Comment    string    `json:"comment"`
}

type InsertEvidances struct {
	SalesID      int       `json:"sales_id"`
	CustomerID   int       `json:"customer_id"`
	DueDate      time.Time `json:"due_date"`
	TypeEvidance string    `json:"type_evidance"`
}
