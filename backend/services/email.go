package services

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client      *resend.Client
	fromAddress string
}

type EmailTemplateData struct {
	Name             string
	Email            string
	DashboardURL     string
	VerificationLink string
	ResetLink        string
	ExpiryHours      int
	FamilyName       string
	InviteLink       string
	SenderName       string
}

func NewEmailService(apiKey string, fromAddress string) *EmailService {
	return &EmailService{
		client:      resend.NewClient(apiKey),
		fromAddress: fromAddress,
	}
}

// LoadTemplate loads and parses an email template
func (es *EmailService) loadTemplate(templateName string) (*template.Template, error) {
	templatePath := fmt.Sprintf("templates/%s.html", templateName)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// renderTemplate renders a template with data
func (es *EmailService) renderTemplate(templateName string, data EmailTemplateData) (string, error) {
	tmpl, err := es.loadTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// SendWelcomeEmail sends a welcome email to a new member
func (es *EmailService) SendWelcomeEmail(to string, name string, dashboardURL string) error {
	data := EmailTemplateData{
		Name:         name,
		Email:        to,
		DashboardURL: dashboardURL,
	}

	html, err := es.renderTemplate("welcome", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    es.fromAddress,
		To:      []string{to},
		Subject: "Welcome to Breakthrough Palace Mission",
		Html:    html,
	}

	_, err = es.client.Emails.Send(params)
	return err
}

// SendVerificationEmail sends an email verification link
func (es *EmailService) SendVerificationEmail(to string, name string, verificationLink string) error {
	data := EmailTemplateData{
		Name:             name,
		Email:            to,
		VerificationLink: verificationLink,
	}

	html, err := es.renderTemplate("verify-email", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    es.fromAddress,
		To:      []string{to},
		Subject: "Verify Your Email Address",
		Html:    html,
	}

	_, err = es.client.Emails.Send(params)
	return err
}

// SendPasswordResetEmail sends a password reset link
func (es *EmailService) SendPasswordResetEmail(to string, name string, resetLink string, expiryHours int) error {
	data := EmailTemplateData{
		Name:        name,
		Email:       to,
		ResetLink:   resetLink,
		ExpiryHours: expiryHours,
	}

	html, err := es.renderTemplate("reset-password", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    es.fromAddress,
		To:      []string{to},
		Subject: "Reset Your Password",
		Html:    html,
	}

	_, err = es.client.Emails.Send(params)
	return err
}

// SendInvitationEmail sends an invitation to join a family
func (es *EmailService) SendInvitationEmail(to string, name string, familyName string, senderName string, inviteLink string, expiryHours int) error {
	data := EmailTemplateData{
		Name:        name,
		Email:       to,
		FamilyName:  familyName,
		SenderName:  senderName,
		InviteLink:  inviteLink,
		ExpiryHours: expiryHours,
	}

	html, err := es.renderTemplate("family-invitation", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    es.fromAddress,
		To:      []string{to},
		Subject: fmt.Sprintf("You're invited to join %s on Breakthrough Palace Mission", familyName),
		Html:    html,
	}

	_, err = es.client.Emails.Send(params)
	return err
}
