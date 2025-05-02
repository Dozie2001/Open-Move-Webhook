package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/open-move/intercord/internal/models"
)

type TeamService struct {
	db           *bun.DB
	emailService *EmailService
}

func NewTeamService(db *bun.DB, emailService *EmailService) *TeamService {
	return &TeamService{
		db:           db,
		emailService: emailService,
	}
}

type CreateTeamInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type InviteToTeamInput struct {
	TeamID int64  `json:"team_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Role   string `json:"role" binding:"required,oneof=admin member"`
}

func (s *TeamService) Create(ctx context.Context, input CreateTeamInput, userID int64) (*models.Team, error) {
	team := &models.Team{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     userID,
	}

	_, err := s.db.NewInsert().Model(team).Exec(ctx)
	if err != nil {
		return nil, err
	}

	membership := &models.TeamMembership{
		TeamID: team.ID,
		UserID: userID,
		Role:   models.TeamRoleOwner,
	}

	_, err = s.db.NewInsert().Model(membership).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetByID(ctx context.Context, id int64) (*models.Team, error) {
	team := new(models.Team)
	err := s.db.NewSelect().Model(team).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) GetTeamsByUserID(ctx context.Context, userID int64) ([]models.Team, error) {
	var teams []models.Team
	err := s.db.NewSelect().
		Model(&teams).
		Join("JOIN team_memberships AS tm ON tm.team_id = team.id").
		Where("tm.user_id = ?", userID).
		Where("tm.deleted_at IS NULL").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (s *TeamService) GetMembership(ctx context.Context, teamID, userID int64) (*models.TeamMembership, error) {
	membership := new(models.TeamMembership)
	err := s.db.NewSelect().
		Model(membership).
		Where("team_id = ?", teamID).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return membership, nil
}

func (s *TeamService) InviteToTeam(ctx context.Context, input InviteToTeamInput, inviterID int64, baseURL string) error {

	inviterMembership := new(models.TeamMembership)
	err := s.db.NewSelect().
		Model(inviterMembership).
		Where("team_id = ?", input.TeamID).
		Where("user_id = ?", inviterID).
		Scan(ctx)

	if err != nil {
		return errors.New("you are not a member of this team")
	}

	if inviterMembership.Role != models.TeamRoleOwner && inviterMembership.Role != models.TeamRoleAdmin {
		return errors.New("you don't have permission to invite members")
	}

	user := new(models.User)
	err = s.db.NewSelect().Model(user).Where("email = ?", input.Email).Scan(ctx)

	if err != nil {

		return errors.New("user with this email does not exist")
	}

	existingMembership := new(models.TeamMembership)
	err = s.db.NewSelect().
		Model(existingMembership).
		Where("team_id = ?", input.TeamID).
		Where("user_id = ?", user.ID).
		Scan(ctx)

	if err == nil {
		return errors.New("user is already a member of this team")
	}

	team := new(models.Team)
	err = s.db.NewSelect().Model(team).Where("id = ?", input.TeamID).Scan(ctx)
	if err != nil {
		return err
	}

	inviter := new(models.User)
	err = s.db.NewSelect().Model(inviter).Where("id = ?", inviterID).Scan(ctx)
	if err != nil {
		return err
	}

	role := models.TeamRoleMember
	if input.Role == "admin" {
		role = models.TeamRoleAdmin
	}

	membership := &models.TeamMembership{
		TeamID: input.TeamID,
		UserID: user.ID,
		Role:   role,
	}

	_, err = s.db.NewInsert().Model(membership).Exec(ctx)
	if err != nil {
		return err
	}

	inviteLink := fmt.Sprintf("%s/teams/%d", baseURL, input.TeamID)
	inviterName := fmt.Sprintf("%s %s", inviter.FirstName, inviter.LastName)
	err = s.emailService.SendTeamInviteEmail(user.Email, inviterName, team.Name, inviteLink)
	if err != nil {
		return err
	}

	return nil
}

func (s *TeamService) JoinTeam(ctx context.Context, teamID, userID int64) error {

	team := new(models.Team)
	err := s.db.NewSelect().Model(team).Where("id = ?", teamID).Scan(ctx)
	if err != nil {
		return errors.New("team not found")
	}

	membership := new(models.TeamMembership)
	err = s.db.NewSelect().
		Model(membership).
		Where("team_id = ?", teamID).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err == nil {
		return errors.New("you are already a member of this team")
	}

	newMembership := &models.TeamMembership{
		TeamID: teamID,
		UserID: userID,
		Role:   models.TeamRoleMember,
	}

	_, err = s.db.NewInsert().Model(newMembership).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *TeamService) LeaveTeam(ctx context.Context, teamID, userID int64) error {

	membership := new(models.TeamMembership)
	err := s.db.NewSelect().
		Model(membership).
		Where("team_id = ?", teamID).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return errors.New("you are not a member of this team")
	}

	if membership.Role == models.TeamRoleOwner {

		var count int
		count, err = s.db.NewSelect().
			Model((*models.TeamMembership)(nil)).
			Where("team_id = ?", teamID).
			Count(ctx)

		if err != nil {
			return err
		}

		if count > 1 {
			return errors.New("as the owner, you cannot leave the team while other members exist; transfer ownership first or delete the team")
		}
	}

	_, err = s.db.NewDelete().
		Model(membership).
		Where("team_id = ?", teamID).
		Where("user_id = ?", userID).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, teamID, userID int64) error {

	membership := new(models.TeamMembership)
	err := s.db.NewSelect().
		Model(membership).
		Where("team_id = ?", teamID).
		Where("user_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return errors.New("you are not a member of this team")
	}

	if membership.Role != models.TeamRoleOwner {
		return errors.New("only the team owner can delete the team")
	}

	_, err = s.db.NewDelete().
		Model((*models.TeamMembership)(nil)).
		Where("team_id = ?", teamID).
		Exec(ctx)

	if err != nil {
		return err
	}

	_, err = s.db.NewDelete().
		Model((*models.Team)(nil)).
		Where("id = ?", teamID).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
