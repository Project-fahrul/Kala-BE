package model

type UserSales struct {
	Name          string `json:"name"`
	ID            int    `json:"id"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	TotalEvidance int    `json:"total"`
	Progress      int    `json:"progress"`
}
