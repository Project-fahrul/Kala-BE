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
