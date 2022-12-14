package repository

import (
	"kala/config"
	"kala/model"
	"kala/repository/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.Users) error
	UpdateUser(user *entity.Users) error
	DeleteUsers(id int) error
	FindUserByEmail(email string) (*entity.Users, error)
	FindUserByID(id int) (*entity.Users, error)
	FindAll(offset int, limit int, role string) ([]entity.Users, error)
	FindAllSales(limit int, offset int) ([]model.UserSales, error)
	FindAllSalesNotVerified() ([]entity.Users, error)
	Total() int
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

var userRepository *UserRepositoryImpl = nil

func User_New() UserRepository {
	if userRepository == nil {
		userRepository = &UserRepositoryImpl{
			db: config.DataSource_New(),
		}
	}
	return userRepository
}

func (u *UserRepositoryImpl) Total() int {
	var d struct {
		Total int
	}
	u.db.Raw("select count(xxx.*) as Total from (SELECT u.id , u.name, u.email, u.phone_number, count(e.*) as total_evidance, sum(case when e.submit_date notnull then 1 else 0 end) as progress  FROM kala.users u left join kala.evidances e on e.sales_id = u.id  where  u.role = 'sales' group by u.id) xxx").
		First(&d)
	return d.Total
}

func (u *UserRepositoryImpl) FindAllSales(limit int, offset int) ([]model.UserSales, error) {
	m := make([]model.UserSales, 0)

	err := u.db.Raw("select * from (SELECT u.id , u.name, u.email, u.phone_number, count(e.*) as total_evidance, sum(case when e.submit_date notnull then 1 else 0 end) as progress  FROM kala.users u left join kala.evidances e on e.sales_id = u.id  where  u.role = 'sales' group by u.id) xxx order by xxx.id desc limit ? offset ?", limit, offset).Scan(&m)
	return m, err.Error
}

func (u *UserRepositoryImpl) FindAllSalesNotVerified() ([]entity.Users, error) {
	m := make([]entity.Users, 0)

	err := u.db.Where("role = 'admin' and verified = false").Find(&m)
	return m, err.Error
}

func (u *UserRepositoryImpl) CreateUser(user *entity.Users) error {
	err := u.db.Create(user)
	if err.Error != nil {
		u.db.Raw("alter sequence kala.user_seq INCREMENT BY -1")
	}
	return err.Error
}

func (u *UserRepositoryImpl) UpdateUser(user *entity.Users) error {
	err := u.db.Model(entity.Users{}).Where("id = ?", user.ID).Updates(user)
	return err.Error
}

func (u *UserRepositoryImpl) DeleteUsers(id int) error {
	err := u.db.Where("id = ?", id).Delete(&entity.Users{})
	return err.Error
}
func (u *UserRepositoryImpl) FindUserByEmail(email string) (*entity.Users, error) {
	var user entity.Users
	err := u.db.Where("email = ?", email).First(&user)
	return &user, err.Error
}

func (u *UserRepositoryImpl) FindUserByID(id int) (*entity.Users, error) {
	var user entity.Users
	err := u.db.Where("id = ?", id).First(&user)
	return &user, err.Error
}

func (u *UserRepositoryImpl) FindAll(offset int, limit int, role string) ([]entity.Users, error) {
	var users []entity.Users
	err := u.db.Where("role = ?", role).Offset(offset).Limit(limit).Find(&users)
	return users, err.Error
}
