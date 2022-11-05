package service

import (
	"fmt"
	"kala/exception"
	"kala/repository"
	"kala/repository/entity"
	"time"
)

func GenerateAllEvidance() {
	allEvidance := make([]entity.Evidances, 0)

	allCustomer, err := repository.CustomerRepository_New().FindUserBirthday()
	exception.ResponseStatusError_New(err)
	allEvidance = append(allEvidance, sliceCustomerToEvidance(allCustomer, "birthday")...)

	allCustomer, err = repository.CustomerRepository_New().FindUserDeadLineAngsuran()
	exception.ResponseStatusError_New(err)
	allEvidance = append(allEvidance, sliceCustomerToEvidance(allCustomer, "angsuran")...)

	allCustomer, err = repository.CustomerRepository_New().FindUserDeadLineSTNK()
	exception.ResponseStatusError_New(err)
	allEvidance = append(allEvidance, sliceCustomerToEvidance(allCustomer, "stnk")...)

	allCustomer, err = repository.CustomerRepository_New().FindUserDeadLineService()
	exception.ResponseStatusError_New(err)
	allEvidance = append(allEvidance, sliceCustomerToEvidance(allCustomer, "service")...)

	repository.EvidanceRepository_New().InsertEvidance(allEvidance)
	repository.Notification_New().InsertMany(sliceEvidanceToNotif(allEvidance))

	repository.EvidanceRepository_New().RemoveExpired()
	repository.Notification_New().RemoveExpired()

	//sync evidance and notif
	fmt.Println("Sync")
	notifs := repository.Notification_New().All()
	notSync := make([]entity.Notifications, 0)
	for _, n := range notifs {
		_, err := repository.EvidanceRepository_New().SelectBySalesIdAndCustomerIdAndTypeEvidance(entity.Evidances{
			SalesID:      n.SalesID,
			CustomerID:   n.CustomerID,
			TypeEvidance: n.TypeNotification,
		})

		// notif and evidance not sync
		if err != nil {
			notSync = append(notSync, n)
		}
	}

	for _, n := range notSync {
		repository.Notification_New().Delete(n)
	}
}

func RemoveExpiredEvidance() {
	repository.Notification_New()
}

func sliceCustomerToEvidance(cus []entity.Customers, evidanceType string) []entity.Evidances {
	evidance := make([]entity.Evidances, 0)
	now := time.Now()
	for _, d := range cus {
		evidance = append(evidance, entity.Evidances{
			SalesID:      d.SalesID,
			CustomerID:   d.ID,
			DueDate:      now,
			TypeEvidance: evidanceType,
		})
	}
	return evidance
}

func sliceEvidanceToNotif(evidance []entity.Evidances) []entity.Notifications {
	notif := make([]entity.Notifications, 0)

	for _, e := range evidance {
		notif = append(notif, entity.Notifications{
			SalesID:          e.SalesID,
			CustomerID:       e.CustomerID,
			Message:          generateMessage(e.TypeEvidance),
			DueDate:          e.DueDate,
			TypeNotification: e.TypeEvidance,
		})
	}
	return notif
}

func generateMessage(typeEvidance string) string {
	if typeEvidance == "stnk" {
		return "Jatuh tempo STNK"
	} else if typeEvidance == "angsuran" {
		return "Jatuh tempo angsuran"
	} else if typeEvidance == "service" {
		return "Waktunya service"
	} else { // == birthday
		return "Hari ini ulang tahun"
	}
}
