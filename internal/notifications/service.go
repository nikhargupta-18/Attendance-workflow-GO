package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"attendance-workflow/pkg/config"
	"attendance-workflow/pkg/db"

	"github.com/hibiken/asynq"
)

const (
	TypeLeaveStatusUpdate = "leave:status_update"
	TypeAttendanceMarked  = "attendance:marked"
	TypeReminderEmail     = "email:reminder"
)

type NotificationService struct {
	client *asynq.Client
	server *asynq.Server
	wg     sync.WaitGroup
	DB     db.GormDB
}

var (
	notificationService *NotificationService
	once                sync.Once
)

func GetNotificationService() *NotificationService {
	once.Do(func() {
		redisOpt := asynq.RedisClientOpt{
			Addr:     config.AppConfig.Redis.Addr,
			Password: config.AppConfig.Redis.Password,
			DB:       config.AppConfig.Redis.DB,
		}
		client := asynq.NewClient(redisOpt)
		server := asynq.NewServer(
			redisOpt,
			asynq.Config{Concurrency: 10},
		)

		notificationService = &NotificationService{
			client: client,
			server: server,
			DB:     db.DB,
		}

		// Start processing jobs
		notificationService.startProcessor()
	})
	return notificationService
}

func (s *NotificationService) startProcessor() {
	mux := asynq.NewServeMux()

	// Handler for leave status updates
	mux.HandleFunc(TypeLeaveStatusUpdate, s.handleLeaveStatusUpdate)
	mux.HandleFunc(TypeAttendanceMarked, s.handleAttendanceMarked)
	mux.HandleFunc(TypeReminderEmail, s.handleReminderEmail)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.Run(mux); err != nil {
			log.Fatalf("Could not run async server: %v", err)
		}
	}()
}

func (s *NotificationService) Stop() {
	s.server.Stop()
	s.wg.Wait()
	s.client.Close()
}

type LeaveStatusUpdatePayload struct {
	LeaveID    uint   `json:"leave_id"`
	StudentID  uint   `json:"student_id"`
	Status     string `json:"status"`
	ApprovedBy string `json:"approved_by,omitempty"`
	Remarks    string `json:"remarks,omitempty"`
}

type AttendancePayload struct {
	StudentID uint      `json:"student_id"`
	Date      time.Time `json:"date"`
	Present   bool      `json:"present"`
	MarkedBy  string    `json:"marked_by"`
}

func (s *NotificationService) QueueLeaveStatusNotification(ctx context.Context, payload LeaveStatusUpdatePayload) error {
	taskBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal leave status payload: %v", err)
	}

	task := asynq.NewTask(TypeLeaveStatusUpdate, taskBytes)
	_, err = s.client.EnqueueContext(ctx, task)
	return err
}

func (s *NotificationService) QueueAttendanceNotification(ctx context.Context, payload AttendancePayload) error {
	taskBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal attendance payload: %v", err)
	}

	task := asynq.NewTask(TypeAttendanceMarked, taskBytes)
	_, err = s.client.EnqueueContext(ctx, task)
	return err
}

func (s *NotificationService) handleLeaveStatusUpdate(ctx context.Context, t *asynq.Task) error {
	var payload LeaveStatusUpdatePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal leave status payload: %v", err)
	}

	// Create notification
	notification := db.Notification{
		UserID:  payload.StudentID,
		Type:    "leave_status",
		Title:   "Leave Request Update",
		Message: fmt.Sprintf("Your leave request has been %s. %s", payload.Status, payload.Remarks),
		IsRead:  false,
	}

	if err := s.DB.Create(&notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %v", err)
	}

	// Send email
	// TODO: Implement actual email sending
	log.Printf("Sending email to student %d: Leave request %d is now %s",
		payload.StudentID, payload.LeaveID, payload.Status)

	return nil
}

func (s *NotificationService) handleAttendanceMarked(ctx context.Context, t *asynq.Task) error {
	var payload AttendancePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal attendance payload: %v", err)
	}

	status := "present"
	if !payload.Present {
		status = "absent"
	}

	notification := db.Notification{
		UserID: payload.StudentID,
		Type:   "attendance",
		Title:  "Attendance Update",
		Message: fmt.Sprintf("Your attendance has been marked as %s for %s by %s",
			status, payload.Date.Format("2006-01-02"), payload.MarkedBy),
		IsRead: false,
	}

	if err := s.DB.Create(&notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %v", err)
	}

	// Send email for absence
	if !payload.Present {
		// TODO: Implement actual email sending
		log.Printf("Sending absence notification email to student %d for date %s",
			payload.StudentID, payload.Date.Format("2006-01-02"))
	}

	return nil
}

func (s *NotificationService) handleReminderEmail(ctx context.Context, t *asynq.Task) error {
	// Send reminders for pending leave requests
	var pendingLeaves []db.LeaveRequest
	if err := s.DB.Where("status = ?", db.StatusPending).Find(&pendingLeaves).Error; err != nil {
		return fmt.Errorf("failed to fetch pending leaves: %v", err)
	}

	for _, leave := range pendingLeaves {
		notification := db.Notification{
			UserID: leave.StudentID,
			Type:   "reminder",
			Title:  "Leave Request Reminder",
			Message: fmt.Sprintf("Your leave request from %s to %s is still pending approval",
				leave.StartDate.Format("2006-01-02"), leave.EndDate.Format("2006-01-02")),
			IsRead: false,
		}

		if err := s.DB.Create(&notification).Error; err != nil {
			log.Printf("Failed to create reminder notification for leave %d: %v", leave.ID, err)
			continue
		}

		// TODO: Implement actual email sending
		log.Printf("Sending reminder email for pending leave request %d to student %d",
			leave.ID, leave.StudentID)
	}

	return nil
}
