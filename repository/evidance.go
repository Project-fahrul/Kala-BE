package repository

import (
	"fmt"
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

func (n *EvidanceRepositoryImpl) RemoveExpired() {
	n.db.Raw("DELETE FROM kala.evidances WHERE due_date < NOW() - interval '8 day'")
}

func (e *EvidanceRepositoryImpl) InsertMany(d []entity.Evidances) error {
	return e.db.Create(d).Error
}

func (e *EvidanceRepositoryImpl) SelectBySalesIdAndCustomerIdAndTypeEvidance(evi entity.Evidances) (*entity.Evidances, error) {
	var test entity.Evidances
	res := e.db.Where("sales_id = ? and customer_id = ? and type_evidance = ?", evi.SalesID, evi.CustomerID, evi.TypeEvidance).First(&test)
	return &test, res.Error
}

func (e *EvidanceRepositoryImpl) SelectBySalesIdAndCustomerIdAndTypeEvidanceAndSubmit(evi entity.Evidances) (*entity.Evidances, error) {
	var test entity.Evidances
	res := e.db.Where("sales_id = ? and customer_id = ? and type_evidance = ? and submit_date notnull", evi.SalesID, evi.CustomerID, evi.TypeEvidance).First(&test)
	return &test, res.Error
}

func (e *EvidanceRepositoryImpl) InsertEvidance(d []entity.Evidances) error {
	var test entity.Evidances
	var evi entity.Evidances
	for _, evi = range d {
		res := e.db.Where("sales_id = ? and customer_id = ? and type_evidance = ?", evi.SalesID, evi.CustomerID, evi.TypeEvidance).First(&test)
		if res.Error != nil && res.Error.Error() == "record not found" {
			e.db.Table("kala.evidances").Create(entity.InsertEvidances{
				SalesID:      evi.SalesID,
				CustomerID:   evi.CustomerID,
				DueDate:      evi.DueDate,
				TypeEvidance: evi.TypeEvidance,
			})
		}
	}
	return nil
}

func (d *EvidanceRepositoryImpl) DeleteByCustomerID(id int) {
	d.db.Where("customer_id = ?", id).Delete(&entity.Evidances{})
}

func (d *EvidanceRepositoryImpl) UploadFile(e entity.Evidances) error {
	err := d.db.Table("kala.evidances").Model(entity.EvidancesUpdate{}).Where("sales_id = ? and customer_id = ? and type_evidance = ?", e.SalesID, e.CustomerID, e.TypeEvidance).Updates(e)
	return err.Error
}

func (d *EvidanceRepositoryImpl) ListEvidance() ([]model.ListEvidance, error) {
	data := make([]model.ListEvidance, 0)

	err := d.db.Raw("select case when e.submit_date notnull then true else false end as submit_date, e.due_date , e.sales_id , e.customer_id , u.name as sales_name, c.name as name, e.type_evidance from kala.evidances e inner join kala.users u on u.id = e.sales_id inner join kala.customers c on c.id = e.customer_id").
		Scan(&data)

	return data, err.Error
}

func (d *EvidanceRepositoryImpl) ListEvidanceWithLimit(limit int, offset int) ([]model.ListEvidance, error) {
	data := make([]model.ListEvidance, 0)
	sql := fmt.Sprintf("select case when e.submit_date notnull then true else false end as submit_date, e.due_date , e.sales_id , e.customer_id , u.name as sales_name, c.name as name, e.type_evidance from kala.evidances e inner join kala.users u on u.id = e.sales_id inner join kala.customers c on c.id = e.customer_id"+
		" limit %d offset %d", limit, offset)
	err := d.db.Raw(sql).
		Scan(&data)

	return data, err.Error
}

func (d *EvidanceRepositoryImpl) Total() int {
	var l struct {
		Total int
	}
	d.db.Raw("SELECT COUNT(*) as total FROM kala.evidances").First(&l)
	return l.Total
}

func (d *EvidanceRepositoryImpl) Evidance(sales int, customer int, due time.Time, typeEvidance string) (entity.Evidances, error) {
	var s entity.Evidances
	err := d.db.Where("sales_id = ? AND customer_id = ? AND due_date = ? AND type_evidance = ?",
		sales, customer, due, typeEvidance).First(&s)

	return s, err.Error
}

func (d *EvidanceRepositoryImpl) Count() (*model.TotalEvidance, error) {
	var total model.TotalEvidance
	err := d.db.Raw("select  sum(case when e.submit_date notnull then 1 else 0 end) as send, count(*) as total, sum(case when e.due_date < now() - interval '6 day' and e.submit_date  is null  then 1 else 0 end) as notsend  from kala.evidances e").Scan(&total)
	return &total, err.Error
}
