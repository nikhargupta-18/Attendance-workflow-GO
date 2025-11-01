package analytics

import (
	"net/http"
	"time"

	"attendance-workflow/pkg/db"

	"github.com/gin-gonic/gin"
)

type AttendanceSummary struct {
	TotalDays      int64   `json:"total_days"`
	PresentDays    int64   `json:"present_days"`
	AbsentDays     int64   `json:"absent_days"`
	AttendanceRate float64 `json:"attendance_rate"`
}

type LeaveSummary struct {
	Total     int64                 `json:"total"`
	ByType    map[string]int64      `json:"by_type"`
	ByStatus  map[string]int64      `json:"by_status"`
	MonthWise []MonthlyLeaveSummary `json:"month_wise"`
}

type MonthlyLeaveSummary struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	// Get date range from query params with defaults
	endDate := time.Now()
	startDate := endDate.AddDate(0, -6, 0) // Last 6 months by default

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsedDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsedDate
		}
	}

	// Get department-wise attendance
	type DeptAttendance struct {
		Dept          string  `json:"dept"`
		AttendanceSum float64 `json:"attendance_rate"`
	}
	var deptAttendance []DeptAttendance

	h.DB.Table("attendances").
		Select("users.dept, (COUNT(CASE WHEN attendances.present THEN 1 END) * 100.0 / COUNT(*)) as attendance_sum").
		Joins("INNER JOIN users ON users.id = attendances.student_id").
		Where("attendances.date BETWEEN ? AND ?", startDate, endDate).
		Where("users.dept IS NOT NULL").
		Group("users.dept").
		Scan(&deptAttendance)

	// Get overall attendance trends
	var monthlyAttendance []struct {
		Month         string  `json:"month"`
		AttendanceSum float64 `json:"attendance_rate"`
	}

	h.DB.Table("attendances").
		Select("DATE_FORMAT(date, '%Y-%m') as month, (COUNT(CASE WHEN present THEN 1 END) * 100.0 / COUNT(*)) as attendance_sum").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Group("month").
		Order("month").
		Scan(&monthlyAttendance)

	// Get leave summary
	var leavesByType map[string]int64
	var leavesByStatus map[string]int64
	var monthlyLeaves []MonthlyLeaveSummary

	h.DB.Model(&db.LeaveRequest{}).
		Select("leave_type, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("leave_type").
		Scan(&leavesByType)

	h.DB.Model(&db.LeaveRequest{}).
		Select("status, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("status").
		Scan(&leavesByStatus)

	h.DB.Model(&db.LeaveRequest{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("month").
		Order("month").
		Scan(&monthlyLeaves)

	// Calculate absentee trends
	type AbsenteeTrend struct {
		StudentID uint    `json:"student_id"`
		Name      string  `json:"name"`
		Dept      string  `json:"dept"`
		Rate      float64 `json:"absence_rate"`
	}
	var absenteeTrends []AbsenteeTrend

	h.DB.Table("attendances").
		Select("users.id as student_id, users.name, users.dept, (COUNT(CASE WHEN NOT attendances.present THEN 1 END) * 100.0 / COUNT(*)) as rate").
		Joins("INNER JOIN users ON users.id = attendances.student_id").
		Where("attendances.date BETWEEN ? AND ?", startDate, endDate).
		Group("users.id, users.name, users.dept").
		Having("rate > ?", 25.0). // Students with >25% absence rate
		Order("rate DESC").
		Limit(10).
		Scan(&absenteeTrends)

	c.JSON(http.StatusOK, gin.H{
		"period": gin.H{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"department_wise_attendance": deptAttendance,
		"monthly_attendance_trends":  monthlyAttendance,
		"leaves": gin.H{
			"by_type":    leavesByType,
			"by_status":  leavesByStatus,
			"month_wise": monthlyLeaves,
		},
		"high_absentee_students": absenteeTrends,
	})
}
