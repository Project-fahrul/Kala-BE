package repository

import (
	"kala/config"
	"kala/repository/entity"
	"time"

	"gorm.io/gorm"
)

type CustomerRepositoryImpl struct {
	db               *gorm.DB
	intervalAngsuran int
	intervalSTNK     int
	intervalService  int
}

type CustomerRepository interface {
	CreateCustomer(cus *entity.Customers) error
	UpdateCustomer(cus *entity.Customers) error
	UpdateCustomerAngsuranForThisYear(cus *entity.Customers) error
	UpdateCustomerSTNKForThisYear(cus *entity.Customers) error
	UpdateCustomerServiceForThisYear(cus *entity.Customers) error
	DeleteCustomer(id int) error
	FindCustomerBySalesID(sales_id int) ([]entity.Customers, error)
	FindCustomerByID(customer_id int) (*entity.Customers, error)
	FindUserBirthdayBy(date *time.Time) ([]entity.Customers, error)
	FindUserDeadLineAngsuranBy(date *time.Time) ([]entity.Customers, error)
	FindUserDeadLineSTNKBy(date *time.Time) ([]entity.Customers, error)
	FindUserDeadLineServiceBy(date *time.Time) ([]entity.Customers, error)
}

var customerRepository *CustomerRepositoryImpl = nil

func CustomerRepository_New() CustomerRepository {
	if customerRepository == nil {
		customerRepository = &CustomerRepositoryImpl{
			db:               config.DataSource_New(),
			intervalAngsuran: 3,
			intervalSTNK:     7,
			intervalService:  3,
		}
	}
	return customerRepository
}

func (c *CustomerRepositoryImpl) CreateCustomer(cus *entity.Customers) error {
	err := c.db.Create(cus)
	return err.Error
}

func (c *CustomerRepositoryImpl) UpdateCustomerAngsuranForThisYear(cus *entity.Customers) error {
	cus.TglAngsuran = cus.TglAngsuran.AddDate(0, 1, 0)
	err := c.db.Model(entity.Users{}).Where("id = ?", cus.ID).Save(cus)
	return err.Error
}

func (c *CustomerRepositoryImpl) UpdateCustomerServiceForThisYear(cus *entity.Customers) error {
	cus.TglAngsuran = cus.TglAngsuran.AddDate(0, 6, 0)
	err := c.db.Model(entity.Users{}).Where("id = ?", cus.ID).Save(cus)
	return err.Error
}

func (c *CustomerRepositoryImpl) UpdateCustomerSTNKForThisYear(cus *entity.Customers) error {
	cus.TglSTNK = cus.TglSTNK.AddDate(1, 0, 0)
	err := c.db.Model(entity.Users{}).Where("id = ?", cus.ID).Save(cus)
	return err.Error
}

func (c *CustomerRepositoryImpl) UpdateCustomer(cus *entity.Customers) error {
	err := c.db.Model(entity.Users{}).Where("id = ?", cus.ID).Save(cus)
	return err.Error
}

func (c *CustomerRepositoryImpl) DeleteCustomer(id int) error {
	err := c.db.Where("id = ?", id).Delete(&entity.Customers{})
	return err.Error
}

func (c *CustomerRepositoryImpl) FindCustomerBySalesID(sales_id int) ([]entity.Customers, error) {
	cus := make([]entity.Customers, 0)
	err := c.db.Where("sales_id = ?", sales_id).Find(&cus)
	return cus, err.Error
}

func (c *CustomerRepositoryImpl) FindCustomerByID(customer_id int) (*entity.Customers, error) {
	cus := entity.Customers{}
	err := c.db.Where("id = ?", customer_id).First(&cus)
	return &cus, err.Error
}

func (c *CustomerRepositoryImpl) FindUserBirthdayBy(date *time.Time) ([]entity.Customers, error) {
	cus := make([]entity.Customers, 0)
	sql := "SELECT * FROM kala.customers cus WHERE EXTRACT( DAY FROM cus.tgl_lahir) = EXTRACT( DAY FROM NOW()) AND EXTRACT( MONTH FROM cus.tgl_lahir) = EXTRACT( MONTH FROM NOW())"
	err := c.db.Raw(sql).Scan(&cus)
	return cus, err.Error
}

func (c *CustomerRepositoryImpl) FindUserDeadLineAngsuranBy(date *time.Time) ([]entity.Customers, error) {
	cus := make([]entity.Customers, 0)
	sql := "SELECT * FROM kala.customers cus WHERE cus.tgl_angsuran BETWEEN NOW() AND NOW() + interval '3 day'"
	err := c.db.Raw(sql).Scan(&cus)
	return cus, err.Error
}

func (c *CustomerRepositoryImpl) FindUserDeadLineSTNKBy(date *time.Time) ([]entity.Customers, error) {
	cus := make([]entity.Customers, 0)
	sql := "SELECT * FROM kala.customers cus WHERE cus.tgl_stnk BETWEEN NOW() AND NOW() + interval '7 day'"
	err := c.db.Raw(sql).Scan(&cus)
	return cus, err.Error
}

func (c *CustomerRepositoryImpl) FindUserDeadLineServiceBy(date *time.Time) ([]entity.Customers, error) {
	cus := make([]entity.Customers, 0)
	sql := "SELECT * FROM kala.customers cus WHERE " +
		"(cus.new_customer == true AND cus.tgl_service + INTERVAL '1 MONTH' BETWEEN NOW() AND NOW() + INTERVAL '3 day')" +
		" OR " +
		"(cus.tgl_service BETWEEN NOW() AND NOW() + INTERVAL '3 day'))"
	err := c.db.Raw(sql).Scan(&cus)
	return cus, err.Error
}
