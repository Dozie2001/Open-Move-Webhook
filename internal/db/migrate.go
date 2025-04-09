package db

import "github.com/Dozie2001/Open-Move-Webhook/internal/models"

func Migrate() error {
	err := DB.AutoMigrate(
		&models.Webhook{},
		&models.Subscription{},
		&models.DeliveryLog{},
	)
	return err
}