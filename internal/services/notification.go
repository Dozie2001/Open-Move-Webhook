package services

import (
	"context"
	"errors"

	"github.com/uptrace/bun"

	"github.com/open-move/intercord/internal/models"
)

type NotificationService struct {
	db          *bun.DB
	teamService *TeamService
}

func NewNotificationService(db *bun.DB, teamService *TeamService) *NotificationService {
	return &NotificationService{
		db:          db,
		teamService: teamService,
	}
}

func (s *NotificationService) GetByID(ctx context.Context, id int64, userID int64) (*models.Notification, error) {
	notification := new(models.Notification)
	err := s.db.NewSelect().Model(notification).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, errors.New("notification not found")
	}

	subscription := new(models.Subscription)
	err = s.db.NewSelect().Model(subscription).Where("id = ?", notification.SubscriptionID).Scan(ctx)
	if err != nil {
		return nil, errors.New("subscription not found")
	}

	if subscription.UserID != userID {
		if subscription.TeamID == nil {
			return nil, errors.New("you don't have permission to view this notification")
		}

		_, err = s.teamService.GetMembership(ctx, *subscription.TeamID, userID)
		if err != nil {
			return nil, errors.New("you don't have permission to view this notification")
		}
	}

	return notification, nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID int64, limit, offset int) ([]models.Notification, int, error) {

	var subscriptions []models.Subscription
	err := s.db.NewSelect().
		Model(&subscriptions).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, 0, err
	}

	var teamMemberships []models.TeamMembership
	err = s.db.NewSelect().
		Model(&teamMemberships).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, 0, err
	}

	var teamIDs []int64
	for _, membership := range teamMemberships {
		teamIDs = append(teamIDs, membership.TeamID)
	}

	if len(subscriptions) == 0 && len(teamIDs) == 0 {
		return []models.Notification{}, 0, nil
	}

	var subscriptionIDs []int64
	for _, subscription := range subscriptions {
		subscriptionIDs = append(subscriptionIDs, subscription.ID)
	}

	if len(teamIDs) > 0 {
		var teamSubscriptions []models.Subscription
		err = s.db.NewSelect().
			Model(&teamSubscriptions).
			Where("team_id IN (?)", bun.In(teamIDs)).
			Scan(ctx)

		if err != nil {
			return nil, 0, err
		}

		for _, subscription := range teamSubscriptions {
			subscriptionIDs = append(subscriptionIDs, subscription.ID)
		}
	}

	if len(subscriptionIDs) == 0 {
		return []models.Notification{}, 0, nil
	}

	count, err := s.db.NewSelect().
		Model((*models.Notification)(nil)).
		Where("subscription_id IN (?)", bun.In(subscriptionIDs)).
		Count(ctx)

	if err != nil {
		return nil, 0, err
	}

	var notifications []models.Notification
	err = s.db.NewSelect().
		Model(&notifications).
		Where("subscription_id IN (?)", bun.In(subscriptionIDs)).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, 0, err
	}

	return notifications, count, nil
}
