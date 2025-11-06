package dto

type MarkAttendanceRequest struct {
	StudentID uint   `json:"student_id" binding:"required"`
	Date      Date   `json:"date" binding:"required"`
	Present   bool   `json:"present"`
}



