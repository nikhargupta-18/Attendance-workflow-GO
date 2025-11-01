package db

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleFaculty UserRole = "faculty"
	RoleWarden  UserRole = "warden"
	RoleStudent UserRole = "student"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      UserRole  `gorm:"not null" json:"role"`
	Dept      string    `json:"dept,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	LeaveRequests []LeaveRequest `gorm:"foreignKey:StudentID" json:"-"`
	Attendance    []Attendance   `gorm:"foreignKey:StudentID" json:"-"`
}

type LeaveStatus string

const (
	StatusPending  LeaveStatus = "pending"
	StatusApproved LeaveStatus = "approved"
	StatusRejected LeaveStatus = "rejected"
)

type LeaveType string

const (
	LeaveTypeMedical   LeaveType = "Medical"
	LeaveTypePersonal  LeaveType = "Personal"
	LeaveTypeEmergency LeaveType = "Emergency"
	LeaveTypeOther     LeaveType = "Other"
)

type LeaveRequest struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	StudentID  uint        `gorm:"not null;index" json:"student_id"`
	Student    User        `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	LeaveType  LeaveType   `gorm:"not null" json:"leave_type"`
	Reason     string      `gorm:"not null;type:text" json:"reason"`
	StartDate  time.Time   `gorm:"not null" json:"start_date"`
	EndDate    time.Time   `gorm:"not null" json:"end_date"`
	Status     LeaveStatus `gorm:"default:pending" json:"status"`
	ApprovedBy *uint       `gorm:"index" json:"approved_by,omitempty"`
	Approver   *User       `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Remarks    string      `json:"remarks,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type Attendance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID uint      `gorm:"not null;index" json:"student_id"`
	Student   User      `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Date      time.Time `gorm:"not null;type:date" json:"date"`
	Present   bool      `gorm:"default:false" json:"present"`
	MarkedBy  uint      `json:"marked_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `gorm:"type:text" json:"message"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func (LeaveRequest) TableName() string {
	return "leave_requests"
}

func (Attendance) TableName() string {
	return "attendances"
}
