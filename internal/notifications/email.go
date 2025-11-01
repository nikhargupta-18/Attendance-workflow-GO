package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"attendance-workflow/pkg/config"
	"attendance-workflow/pkg/db"

	"github.com/hibiken/asynq"
)

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func loadEmailConfig() EmailConfig {
	return EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}
}

type EmailSender interface {
	Send(to, subject, body string) error
}

type SMTPEmailSender struct {
	config EmailConfig
}

func NewSMTPEmailSender(config EmailConfig) *SMTPEmailSender {
	return &SMTPEmailSender{config: config}
}

func (s *SMTPEmailSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.config.From, to, subject, body)

	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(msg))
}

func (s *NotificationService) InitializeEmailWorker() {
	emailConfig := loadEmailConfig()
	emailSender := NewSMTPEmailSender(emailConfig)

	mux := asynq.NewServeMux()

	// Handler for processing email notifications
	mux.HandleFunc("email:send", func(ctx context.Context, t *asynq.Task) error {
		var notification db.EmailNotification
		if err := json.Unmarshal(t.Payload(), &notification); err != nil {
			return fmt.Errorf("failed to unmarshal email notification: %v", err)
		}

		// Get user email
		var user db.User
		if err := s.DB.First(&user, notification.UserID).Error; err != nil {
			return fmt.Errorf("failed to find user: %v", err)
		}

		// Send email
		if err := emailSender.Send(user.Email, notification.Subject, notification.Body); err != nil {
			notification.Status = "failed"
			notification.Error = err.Error()
		} else {
			notification.Status = "sent"
			now := time.Now()
			notification.SentAt = &now
		}

		// Update notification status
		if err := s.DB.Save(&notification).Error; err != nil {
			return fmt.Errorf("failed to update notification status: %v", err)
		}

		return nil
	})

	// Start the email worker
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.AppConfig.Redis.Addr},
		asynq.Config{Concurrency: 10},
	)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("Failed to run email worker server: %v", err)
		}
	}()
}

// QueueEmailNotification queues an email to be sent asynchronously
func (s *NotificationService) QueueEmailNotification(userID uint, subject, body string) error {
	notification := db.EmailNotification{
		UserID:  userID,
		Subject: subject,
		Body:    body,
		Status:  "pending",
	}

	if err := s.DB.Create(&notification).Error; err != nil {
		return fmt.Errorf("failed to create email notification: %v", err)
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal email notification: %v", err)
	}

	task := asynq.NewTask("email:send", payload)
	if _, err := s.client.Enqueue(task); err != nil {
		return fmt.Errorf("failed to enqueue email task: %v", err)
	}

	return nil
}
