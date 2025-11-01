package db

import (
	"time"
)

// AnalyticsSummary represents the analytics summary model
type AnalyticsSummary struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Date           time.Time `gorm:"type:date;index" json:"date"`
	TotalPresent   int64     `json:"total_present"`
	TotalAbsent    int64     `json:"total_absent"`
	LeavesApproved int64     `json:"leaves_approved"`
	LeavesPending  int64     `json:"leaves_pending"`
	LeavesRejected int64     `json:"leaves_rejected"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (AnalyticsSummary) TableName() string {
	return "analytics_summaries"
}

// EmailNotification represents the email notification model
type EmailNotification struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Subject   string     `gorm:"not null" json:"subject"`
	Body      string     `gorm:"type:text;not null" json:"body"`
	Status    string     `gorm:"default:pending" json:"status"` // pending, sent, failed
	Error     string     `json:"error,omitempty"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (EmailNotification) TableName() string {
	return "email_notifications"
}
