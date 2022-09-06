package entity

import "time"

type Customers struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	NoHp          string    `json:"no_hp"`
	TglDec        time.Time `json:"tgl_dec"`
	TglLahir      time.Time `json:"tgl_lahir"`
	TglSTNK       time.Time `json:"tgl_stnk"`
	TglAngsuran   time.Time `json:"tgl_angsuran"`
	NoRangka      string    `json:"no_rangka"`
	TypeKendaraan string    `json:"type_kendaraan"`
	Leasing       string    `json:"leasing"`
	SalesID       int       `json:"sales_id"`
	NewCustomer   bool      `json:"new_customer"`
}
