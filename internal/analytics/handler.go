package analytics

import (
	"net/http"

	"attendance-workflow/pkg/db"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	DB db.GormDB
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{DB: db.DB}
}

func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	// Total users by role
	var totalUsers int64
	var students int64
	var faculty int64
	var wardens int64
	var admins int64

	h.DB.Model(&db.User{}).Count(&totalUsers)
	h.DB.Model(&db.User{}).Where("role = ?", db.RoleStudent).Count(&students)
	h.DB.Model(&db.User{}).Where("role = ?", db.RoleFaculty).Count(&faculty)
	h.DB.Model(&db.User{}).Where("role = ?", db.RoleWarden).Count(&wardens)
	h.DB.Model(&db.User{}).Where("role = ?", db.RoleAdmin).Count(&admins)

	// Leave requests by status
	var pending int64
	var approved int64
	var rejected int64
	h.DB.Model(&db.LeaveRequest{}).Where("status = ?", db.StatusPending).Count(&pending)
	h.DB.Model(&db.LeaveRequest{}).Where("status = ?", db.StatusApproved).Count(&approved)
	h.DB.Model(&db.LeaveRequest{}).Where("status = ?", db.StatusRejected).Count(&rejected)

	// Total attendance records
	var totalAttendance int64
	h.DB.Model(&db.Attendance{}).Count(&totalAttendance)

	c.JSON(http.StatusOK, gin.H{
		"users": gin.H{
			"total":    totalUsers,
			"students": students,
			"faculty":  faculty,
			"wardens":  wardens,
			"admins":   admins,
		},
		"leaves": gin.H{
			"pending":  pending,
			"approved": approved,
			"rejected": rejected,
		},
		"attendance": totalAttendance,
	})
}

func (h *AnalyticsHandler) GetLeaveBreakdown(c *gin.Context) {
	type Result struct {
		LeaveType string
		Count     int64
	}

	var results []Result
	h.DB.Model(&db.LeaveRequest{}).
		Select("leave_type, count(*) as count").
		Group("leave_type").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (h *AnalyticsHandler) GetDepartmentStats(c *gin.Context) {
	dept := c.Query("dept")
	if dept == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department parameter required"})
		return
	}

	// Students in department
	var students int64
	h.DB.Model(&db.User{}).Where("role = ? AND dept = ?", db.RoleStudent, dept).Count(&students)

	// Total attendance
	var presentDays int64
	var totalDays int64
	h.DB.Table("attendances").
		Joins("INNER JOIN users ON users.id = attendances.student_id").
		Where("users.dept = ? AND attendances.present = ?", dept, true).
		Count(&presentDays)
	h.DB.Table("attendances").
		Joins("INNER JOIN users ON users.id = attendances.student_id").
		Where("users.dept = ?", dept).
		Count(&totalDays)

	var percentage float64
	if totalDays > 0 {
		percentage = float64(presentDays) / float64(totalDays) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"department": dept,
		"students":   students,
		"attendance": gin.H{
			"present_days":          presentDays,
			"total_days":            totalDays,
			"attendance_percentage": percentage,
		},
	})
}

func (h *AnalyticsHandler) GetFrequentAbsentees(c *gin.Context) {
	type Result struct {
		StudentID  uint
		Name       string
		AbsentDays int64
	}

	var results []Result
	h.DB.Table("attendances").
		Select("attendances.student_id, users.name, count(*) as absent_days").
		Joins("INNER JOIN users ON users.id = attendances.student_id").
		Where("attendances.present = ?", false).
		Group("attendances.student_id, users.name").
		Order("absent_days DESC").
		Limit(10).
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}
