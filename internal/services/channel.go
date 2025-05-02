package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uptrace/bun"

	"github.com/open-move/intercord/internal/models"
)

type ChannelService struct {
	db          *bun.DB
	teamService *TeamService
}

func NewChannelService(db *bun.DB, teamService *TeamService) *ChannelService {
	return &ChannelService{
		db:          db,
		teamService: teamService,
	}
}

type ChannelConfig struct {
	WebhookURL     string `json:"webhook_url,omitempty"`
	EmailAddress   string `json:"email_address,omitempty"`
	TelegramChatID string `json:"telegram_chat_id,omitempty"`
	DiscordWebhook string `json:"discord_webhook,omitempty"`
}

type CreateChannelInput struct {
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description"`
	Type        string        `json:"type" binding:"required,oneof=webhook email telegram discord"`
	Config      ChannelConfig `json:"config" binding:"required"`
	TeamID      *int64        `json:"team_id"`
}

type UpdateChannelInput struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Config      *ChannelConfig `json:"config"`
}

type SubscribeChannelInput struct {
	ChannelID      int64 `json:"channel_id" binding:"required"`
	SubscriptionID int64 `json:"subscription_id" binding:"required"`
}

func (s *ChannelService) Create(ctx context.Context, input CreateChannelInput, userID int64) (*models.Channel, error) {

	if input.TeamID != nil && *input.TeamID != 0 {
		membership, err := s.teamService.GetMembership(ctx, *input.TeamID, userID)
		if err != nil {
			return nil, errors.New("you are not a member of this team")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return nil, errors.New("you don't have permission to create channels for this team")
		}
	}

	switch models.ChannelType(input.Type) {
	case models.ChannelTypeWebhook:
		if input.Config.WebhookURL == "" {
			return nil, errors.New("webhook URL is required for webhook channels")
		}
	case models.ChannelTypeEmail:
		if input.Config.EmailAddress == "" {
			return nil, errors.New("email address is required for email channels")
		}
	case models.ChannelTypeTelegram:
		if input.Config.TelegramChatID == "" {
			return nil, errors.New("chat ID is required for telegram channels")
		}
	case models.ChannelTypeDiscord:
		if input.Config.DiscordWebhook == "" {
			return nil, errors.New("webhook URL is required for discord channels")
		}
	default:
		return nil, errors.New("invalid channel type")
	}

	configJSON, err := json.Marshal(input.Config)
	if err != nil {
		return nil, err
	}

	channel := &models.Channel{
		Name:        input.Name,
		Description: input.Description,
		Type:        models.ChannelType(input.Type),
		Config:      string(configJSON),
		TeamID:      input.TeamID,
		UserID:      userID,
	}

	_, err = s.db.NewInsert().Model(channel).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *ChannelService) GetByID(ctx context.Context, id int64) (*models.Channel, error) {
	channel := new(models.Channel)
	err := s.db.NewSelect().Model(channel).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (s *ChannelService) GetUserChannels(ctx context.Context, userID int64) ([]models.Channel, error) {
	var channels []models.Channel
	err := s.db.NewSelect().
		Model(&channels).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (s *ChannelService) GetTeamChannels(ctx context.Context, teamID int64) ([]models.Channel, error) {
	var channels []models.Channel
	err := s.db.NewSelect().
		Model(&channels).
		Where("team_id = ?", teamID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (s *ChannelService) Update(ctx context.Context, id int64, input UpdateChannelInput, userID int64) (*models.Channel, error) {
	channel := new(models.Channel)
	err := s.db.NewSelect().Model(channel).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, errors.New("channel not found")
	}

	if channel.UserID != userID {
		if channel.TeamID == nil {
			return nil, errors.New("you don't have permission to update this channel")
		}

		membership, err := s.teamService.GetMembership(ctx, *channel.TeamID, userID)
		if err != nil {
			return nil, errors.New("you don't have permission to update this channel")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return nil, errors.New("you don't have permission to update this channel")
		}
	}

	if input.Name != "" {
		channel.Name = input.Name
	}

	if input.Description != "" {
		channel.Description = input.Description
	}

	if input.Config != nil {

		switch channel.Type {
		case models.ChannelTypeWebhook:
			if input.Config.WebhookURL == "" {
				return nil, errors.New("webhook URL is required for webhook channels")
			}
		case models.ChannelTypeEmail:
			if input.Config.EmailAddress == "" {
				return nil, errors.New("email address is required for email channels")
			}
		case models.ChannelTypeTelegram:
			if input.Config.TelegramChatID == "" {
				return nil, errors.New("chat ID is required for telegram channels")
			}
		case models.ChannelTypeDiscord:
			if input.Config.DiscordWebhook == "" {
				return nil, errors.New("webhook URL is required for discord channels")
			}
		}

		configJSON, err := json.Marshal(input.Config)
		if err != nil {
			return nil, err
		}

		channel.Config = string(configJSON)
	}

	_, err = s.db.NewUpdate().Model(channel).
		Column("name", "description", "config", "updated_at").
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *ChannelService) Delete(ctx context.Context, id int64, userID int64) error {
	channel := new(models.Channel)
	err := s.db.NewSelect().Model(channel).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return errors.New("channel not found")
	}

	if channel.UserID != userID {
		if channel.TeamID == nil {
			return errors.New("you don't have permission to delete this channel")
		}

		membership, err := s.teamService.GetMembership(ctx, *channel.TeamID, userID)
		if err != nil {
			return errors.New("you don't have permission to delete this channel")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return errors.New("you don't have permission to delete this channel")
		}
	}

	_, err = s.db.NewDelete().
		Model((*models.SubscriptionChannel)(nil)).
		Where("channel_id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	_, err = s.db.NewDelete().
		Model(channel).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (s *ChannelService) SubscribeToSubscription(ctx context.Context, input SubscribeChannelInput, userID int64) error {

	subscription := new(models.Subscription)
	err := s.db.NewSelect().Model(subscription).Where("id = ?", input.SubscriptionID).Scan(ctx)
	if err != nil {
		return errors.New("subscription not found")
	}

	if subscription.UserID != userID {
		if subscription.TeamID == nil {
			return errors.New("you don't have permission to modify this subscription")
		}

		membership, err := s.teamService.GetMembership(ctx, *subscription.TeamID, userID)
		if err != nil {
			return errors.New("you don't have permission to modify this subscription")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return errors.New("you don't have permission to modify this subscription")
		}
	}

	channel := new(models.Channel)
	err = s.db.NewSelect().Model(channel).Where("id = ?", input.ChannelID).Scan(ctx)
	if err != nil {
		return errors.New("channel not found")
	}

	if channel.UserID != userID {
		if channel.TeamID == nil {
			return errors.New("you don't have permission to use this channel")
		}

		membership, err := s.teamService.GetMembership(ctx, *channel.TeamID, userID)
		if err != nil {
			return errors.New("you don't have permission to use this channel")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return errors.New("you don't have permission to use this channel")
		}
	}

	if subscription.TeamID != nil && channel.TeamID != nil {
		if *subscription.TeamID != *channel.TeamID {
			return errors.New("subscription and channel must belong to the same team")
		}
	}

	existingSubscription := new(models.SubscriptionChannel)
	err = s.db.NewSelect().
		Model(existingSubscription).
		Where("subscription_id = ?", input.SubscriptionID).
		Where("channel_id = ?", input.ChannelID).
		Scan(ctx)

	if err == nil {
		return errors.New("channel is already subscribed to this subscription")
	}

	subscriptionChannel := &models.SubscriptionChannel{
		SubscriptionID: input.SubscriptionID,
		ChannelID:      input.ChannelID,
	}

	_, err = s.db.NewInsert().Model(subscriptionChannel).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ChannelService) UnsubscribeFromSubscription(ctx context.Context, input SubscribeChannelInput, userID int64) error {

	subscription := new(models.Subscription)
	err := s.db.NewSelect().Model(subscription).Where("id = ?", input.SubscriptionID).Scan(ctx)
	if err != nil {
		return errors.New("subscription not found")
	}

	if subscription.UserID != userID {
		if subscription.TeamID == nil {
			return errors.New("you don't have permission to modify this subscription")
		}

		membership, err := s.teamService.GetMembership(ctx, *subscription.TeamID, userID)
		if err != nil {
			return errors.New("you don't have permission to modify this subscription")
		}

		if membership.Role != models.TeamRoleOwner && membership.Role != models.TeamRoleAdmin {
			return errors.New("you don't have permission to modify this subscription")
		}
	}

	_, err = s.db.NewDelete().
		Model((*models.SubscriptionChannel)(nil)).
		Where("subscription_id = ?", input.SubscriptionID).
		Where("channel_id = ?", input.ChannelID).
		Exec(ctx)

	if err != nil {
		return errors.New("failed to unsubscribe channel")
	}

	return nil
}
