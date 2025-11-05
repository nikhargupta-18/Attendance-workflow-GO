package dto

import "time"

type MarkAttendanceRequest struct {
	StudentID uint      `json:"student_id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Present   bool      `json:"present"`
}


