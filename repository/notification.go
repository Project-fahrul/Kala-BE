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

func (n *NotificationRepositoryImpl) Delete(e entity.Evidances) {
	err := n.db.Where("sales_id = ? and customer_id = ? and type_notification = ?", e.SalesID, e.CustomerID, e.TypeEvidance).Delete(entity.Notifications{})
	exception.ResponseStatusError_New(err.Error)
}

func (n *NotificationRepositoryImpl) ListAllNotificationBySalesID(id int) ([]model.ListNotifiation, error) {

	var notif []model.ListNotifiation
	err := n.db.Raw(fmt.Sprintf("select n.message, n.sales_id , n.customer_id , c.name, n.type_notification from kala.notifications n inner join kala.customers c on n.customer_id = c.id where n.sales_id = %d", id)).Scan(&notif)

	return notif, err.Error
}
