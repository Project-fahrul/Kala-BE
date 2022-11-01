package repository

import (
	"fmt"
	"kala/config"
	"kala/exception"
	"kala/model"
	"kala/repository/entity"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	ListAllNotificationBySalesID(id int) ([]model.ListNotifiation, error)
	Delete(e entity.Evidances)
	DeleteByCustomerID(id int)
	InsertMany(notif []entity.Notifications)
	RemoveExpired()
}

type NotificationRepositoryImpl struct {
	db *gorm.DB
}

var notif *NotificationRepositoryImpl = nil

func Notification_New() NotificationRepository {
	if notif == nil {
		notif = &NotificationRepositoryImpl{
			db: config.DataSource_New(),
		}
	}

	return notif
}

func (n *NotificationRepositoryImpl) InsertMany(notif []entity.Notifications) {
	var test entity.Notifications
	for _, note := range notif {
		res := n.db.Where("customer_id = ? and due_date = ? and sales_id = ? and type_notification = ?", note.CustomerID, note.DueDate, note.SalesID, note.TypeNotification).First(&test)
		if res.Error != nil && res.Error.Error() == "record not found" {
			n.db.Create(note)
		}
	}
}

func (n *NotificationRepositoryImpl) RemoveExpired() {
	n.db.Raw("DELETE FROM kala.notifications WHERE due_date < NOW() - interval '8 day'")
}

func (n *NotificationRepositoryImpl) Delete(e entity.Evidances) {
	err := n.db.Where("sales_id = ? and customer_id = ? and type_notification = ?", e.SalesID, e.CustomerID, e.TypeEvidance).Delete(entity.Notifications{})
	exception.ResponseStatusError_New(err.Error)
}

func (n *NotificationRepositoryImpl) DeleteByCustomerID(id int) {
	err := n.db.Where("customer_id = ?", id).Delete(entity.Notifications{})
	exception.ResponseStatusError_New(err.Error)
}

func (n *NotificationRepositoryImpl) ListAllNotificationBySalesID(id int) ([]model.ListNotifiation, error) {

	var notif []model.ListNotifiation
	err := n.db.Raw(fmt.Sprintf("select n.message, n.sales_id , n.customer_id , c.name, n.type_notification from kala.notifications n inner join kala.customers c on n.customer_id = c.id where n.sales_id = %d", id)).Scan(&notif)

	return notif, err.Error
}
