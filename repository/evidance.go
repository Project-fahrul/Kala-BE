package repository

import (
	"kala/config"
	"kala/model"
	"kala/repository/entity"
	"time"

	"gorm.io/gorm"
)

type EvidanceRepositoryImpl struct {
	db *gorm.DB
}

var db *EvidanceRepositoryImpl = nil

func EvidanceRepository_New() *EvidanceRepositoryImpl {
	if db == nil {
		db = &EvidanceRepositoryImpl{
			db: config.DataSource_New(),
		}
	}

	return db
}

func (d *EvidanceRepositoryImpl) UploadFile(e entity.Evidances) error {
	err := d.db.Table("kala.evidances").Model(entity.EvidancesUpdate{}).Where("sales_id = ? and customer_id = ? and due_date = ? and type_evidance = ?", e.CustomerID, e.SalesID, e.DueDate, e.TypeEvidance).Updates(e)
	return err.Error
}

func (d *EvidanceRepositoryImpl) ListEvidance() ([]model.ListEvidance, error) {
	data := make([]model.ListEvidance, 0)

	err := d.db.Raw("select case when e.submit_date notnull then true else false end as submit_date, e.due_date , e.sales_id , e.customer_id , u.name as sales_name, c.name as name, e.type_evidance from kala.evidances e inner join kala.users u on u.id = e.sales_id inner join kala.customers c on c.id = e.customer_id").
		Scan(&data)

	return data, err.Error
}

func (d *EvidanceRepositoryImpl) Evidance(sales int, customer int, due time.Time, typeEvidance string) (entity.Evidances, error) {
	var s entity.Evidances
	err := d.db.Where("sales_id = ? AND customer_id = ? AND due_date = ? AND type_evidance = ?",
		sales, customer, due, typeEvidance).First(&s)

	return s, err.Error
}

func (d *EvidanceRepositoryImpl) Count() (*model.TotalEvidance, error) {
	var total model.TotalEvidance
	err := d.db.Raw("select  sum(case when e.submit_date = null then 0 else 1 end) as send, count(*) as total  from kala.evidances e").Scan(&total)
	return &total, err.Error
}
