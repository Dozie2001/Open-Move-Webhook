package services

import (
	"context"
	"errors"
	"time"

	"github.com/uptrace/bun"

	"github.com/open-move/intercord/internal/config"
	"github.com/open-move/intercord/internal/models"
	"github.com/open-move/intercord/internal/utils"
)

type UserService struct {
	db           *bun.DB
	jwtConfig    *config.JWTConfig
	emailService *EmailService
}

func NewUserService(db *bun.DB, jwtConfig *config.JWTConfig, emailService *EmailService) *UserService {
	return &UserService{
		db:           db,
		jwtConfig:    jwtConfig,
		emailService: emailService,
	}
}

type RegisterUserInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthResponse struct {
	User         models.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

func (s *UserService) Register(ctx context.Context, input RegisterUserInput, baseURL string) (*models.User, error) {

	existingUser := new(models.User)
	err := s.db.NewSelect().Model(existingUser).Where("email = ?", input.Email).Scan(ctx)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Verified:  false,
	}

	_, err = s.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	verification := &models.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	_, err = s.db.NewInsert().Model(verification).Exec(ctx)
	if err != nil {
		return nil, err
	}

	err = s.emailService.SendVerificationEmail(user.Email, token, baseURL)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	user := new(models.User)
	err := s.db.NewSelect().Model(user).Where("email = ?", input.Email).Scan(ctx)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.jwtConfig)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	verification := new(models.EmailVerification)
	err := s.db.NewSelect().Model(verification).
		Where("token = ?", token).
		Where("used = ?", false).
		Where("expires_at > ?", time.Now()).
		Scan(ctx)

	if err != nil {
		return errors.New("invalid or expired token")
	}

	_, err = s.db.NewUpdate().Model(verification).
		Set("used = ?", true).
		Where("id = ?", verification.ID).
		Exec(ctx)

	if err != nil {
		return err
	}

	_, err = s.db.NewUpdate().Model(&models.User{}).
		Set("verified = ?", true).
		Where("id = ?", verification.UserID).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) RequestPasswordReset(ctx context.Context, email, baseURL string) error {
	user := new(models.User)
	err := s.db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil
	}

	token, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	reset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	_, err = s.db.NewInsert().Model(reset).Exec(ctx)
	if err != nil {
		return err
	}

	return s.emailService.SendPasswordResetEmail(user.Email, token, baseURL)
}

func (s *UserService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	reset := new(models.PasswordReset)
	err := s.db.NewSelect().Model(reset).
		Where("token = ?", input.Token).
		Where("used = ?", false).
		Where("expires_at > ?", time.Now()).
		Scan(ctx)

	if err != nil {
		return errors.New("invalid or expired token")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}

	_, err = s.db.NewUpdate().Model(&models.User{}).
		Set("password = ?", hashedPassword).
		Where("id = ?", reset.UserID).
		Exec(ctx)

	if err != nil {
		return err
	}

	_, err = s.db.NewUpdate().Model(reset).
		Set("used = ?", true).
		Where("id = ?", reset.ID).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := new(models.User)
	err := s.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
