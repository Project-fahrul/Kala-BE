package repository

import (
	"fmt"
	"kala/config"
	"kala/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	ListAllNotificationBySalesID(id int) ([]model.ListNotifiation, error)
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

func (n *NotificationRepositoryImpl) ListAllNotificationBySalesID(id int) ([]model.ListNotifiation, error) {

	var notif []model.ListNotifiation
	err := n.db.Raw(fmt.Sprintf("select n.message, n.sales_id , n.customer_id , c.name, n.type_notification from kala.notifications n inner join kala.customers c on n.customer_id = c.id where n.sales_id = %d", id)).Scan(&notif)

	return notif, err.Error
}
