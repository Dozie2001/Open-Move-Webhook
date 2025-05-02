package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/open-move/intercord/internal/config"
)

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(config *config.EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	var emailBody bytes.Buffer
	emailTemplate := `From: {{.From}} <{{.FromEmail}}>
To: {{.To}}
Subject: {{.Subject}}
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

{{.Body}}
`
	t := template.Must(template.New("email").Parse(emailTemplate))
	err := t.Execute(&emailBody, struct {
		From      string
		FromEmail string
		To        string
		Subject   string
		Body      string
	}{
		From:      s.config.FromName,
		FromEmail: s.config.FromEmail,
		To:        to,
		Subject:   subject,
		Body:      body,
	})

	if err != nil {
		return err
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.FromEmail,
		[]string{to},
		emailBody.Bytes(),
	)
}

func (s *EmailService) SendVerificationEmail(to, token string, baseURL string) error {
	subject := "Verify your email address"
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", baseURL, token)
	body := fmt.Sprintf(`
	<h1>Verify your email address</h1>
	<p>Thank you for registering with Intercord. Please click the link below to verify your email address:</p>
	<p><a href="%s">Verify Email</a></p>
	<p>If you did not register for an account, please ignore this email.</p>
	`, verificationLink)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(to, token string, baseURL string) error {
	subject := "Reset your password"
	resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", baseURL, token)
	body := fmt.Sprintf(`
	<h1>Reset your password</h1>
	<p>You have requested to reset your password. Please click the link below to reset your password:</p>
	<p><a href="%s">Reset Password</a></p>
	<p>If you did not request a password reset, please ignore this email.</p>
	`, resetLink)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendTeamInviteEmail(to, inviterName, teamName, inviteLink string) error {
	subject := fmt.Sprintf("Invitation to join %s team", teamName)
	body := fmt.Sprintf(`
	<h1>Team Invitation</h1>
	<p>%s has invited you to join the %s team on Intercord.</p>
	<p><a href="%s">Accept Invitation</a></p>
	<p>If you do not wish to join this team, please ignore this email.</p>
	`, inviterName, teamName, inviteLink)

	return s.SendEmail(to, subject, body)
}