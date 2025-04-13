package services

import (
	"github.com/Dozie2001/Open-Move-Webhook/internal/db"
	"github.com/Dozie2001/Open-Move-Webhook/internal/models"
	// "github.com/Dozie2001/Open-Move-Webhook/pkg/types"
	// "errors"
	// "fmt"

	// "gorm.io/gorm"
)


func CreateSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	if err := db.DB.Create(subscription).Error; err != nil {
		return nil, err
	}


	return subscription, nil
}