package services

import (
	"context"
	"errors"

	"github.com/uptrace/bun"

	"github.com/open-move/intercord/internal/models"
)

type SubscriptionService struct {
	db          *bun.DB
	teamService *TeamService
}

func NewSubscriptionService(db *bun.DB, teamService *TeamService) *SubscriptionService {
	return &SubscriptionService{
		db:          db,
		teamService: teamService,
	}
}

type CreateSubscriptionInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	EventType   string `json:"event_type" binding:"required"`
	TeamID      *int64 `json:"team_id"`
}

type UpdateSubscriptionInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	EventType   string `json:"event_type"`
	IsActive    *bool  `json:"is_active"`
}

func (s *SubscriptionService) Create(ctx context.Context, input CreateSubscriptionInput, userID int64) (*models.Subscription, error) {

	if input.TeamID != nil && *input.TeamID != 0 {
		membership, err := s.teamService.GetMembership(ctx, *input.TeamID, userID)
		if err != nil {
			return nil, errors.New("you are not a member of this team")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return nil, errors.New("you don't have permission to create subscriptions for this team")
		}
	}

	subscription := &models.Subscription{
		Name:        input.Name,
		Description: input.Description,
		EventType:   input.EventType,
		TeamID:      input.TeamID,
		UserID:      userID,
		IsActive:    true,
	}

	_, err := s.db.NewInsert().Model(subscription).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int64) (*models.Subscription, error) {
	subscription := new(models.Subscription)
	err := s.db.NewSelect().Model(subscription).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *SubscriptionService) GetUserSubscriptions(ctx context.Context, userID int64) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := s.db.NewSelect().
		Model(&subscriptions).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionService) GetTeamSubscriptions(ctx context.Context, teamID int64) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := s.db.NewSelect().
		Model(&subscriptions).
		Where("team_id = ?", teamID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionService) Update(ctx context.Context, id int64, input UpdateSubscriptionInput, userID int64) (*models.Subscription, error) {
	subscription := new(models.Subscription)
	err := s.db.NewSelect().Model(subscription).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, errors.New("subscription not found")
	}

	if subscription.UserID != userID {
		if subscription.TeamID == nil {
			return nil, errors.New("you don't have permission to update this subscription")
		}

		membership, err := s.teamService.GetMembership(ctx, *subscription.TeamID, userID)
		if err != nil {
			return nil, errors.New("you don't have permission to update this subscription")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return nil, errors.New("you don't have permission to update this subscription")
		}
	}

	if input.Name != "" {
		subscription.Name = input.Name
	}

	if input.Description != "" {
		subscription.Description = input.Description
	}

	if input.EventType != "" {
		subscription.EventType = input.EventType
	}

	if input.IsActive != nil {
		subscription.IsActive = *input.IsActive
	}

	_, err = s.db.NewUpdate().Model(subscription).
		Column("name", "description", "event_type", "is_active", "updated_at").
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64, userID int64) error {
	subscription := new(models.Subscription)
	err := s.db.NewSelect().Model(subscription).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return errors.New("subscription not found")
	}

	if subscription.UserID != userID {
		if subscription.TeamID == nil {
			return errors.New("you don't have permission to delete this subscription")
		}

		membership, err := s.teamService.GetMembership(ctx, *subscription.TeamID, userID)
		if err != nil {
			return errors.New("you don't have permission to delete this subscription")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return errors.New("you don't have permission to delete this subscription")
		}
	}

	_, err = s.db.NewDelete().
		Model((*models.SubscriptionChannel)(nil)).
		Where("subscription_id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	_, err = s.db.NewDelete().
		Model(subscription).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
