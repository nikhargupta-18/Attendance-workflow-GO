package attendance

import (
	"net/http"
	"strconv"
	"time"

	"attendance-workflow/internal/dto"
	"attendance-workflow/pkg/db"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	DB db.GormDB
}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{DB: db.DB}
}

// MarkAttendance godoc
// @Summary      Mark attendance
// @Description  Mark attendance for a student
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.MarkAttendanceRequest  true  "Attendance data"
// @Success      201      {object}  object{message=string,data=object}
// @Success      200      {object}  object{message=string,data=object}
// @Failure      400      {object}  object{error=string}
// @Failure      401      {object}  object{error=string}
// @Router       /attendance/mark [post]
func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	markedBy, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req dto.MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, ok := markedBy.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	attendance := db.Attendance{
		StudentID: req.StudentID,
		Date:      req.Date,
		Present:   req.Present,
		MarkedBy:  uint(uid),
	}

	// Check if attendance already marked
	var existing db.Attendance
	if err := h.DB.Where("student_id = ? AND date = ?", req.StudentID, req.Date).First(&existing).Error; err == nil {
		// Update existing
		existing.Present = req.Present
		if err := h.DB.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully", "data": existing})
		return
	}

	if err := h.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark attendance"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Attendance marked successfully", "data": attendance})
}

// GetStudentAttendance godoc
// @Summary      Get student attendance
// @Description  Get attendance records for a specific student
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      int     true   "Student ID"
// @Param        start_date query     string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date   query     string  false  "End date (YYYY-MM-DD)"
// @Success      200        {object}  object{data=array,stats=object}
// @Failure      401        {object}  object{error=string}
// @Router       /attendance/student/{id} [get]
func (h *AttendanceHandler) GetStudentAttendance(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var attendance []db.Attendance
	query := h.DB.Where("student_id = ?", studentID)

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	query.Order("date DESC").Find(&attendance)

	// Calculate stats
	var presentCount int64
	var totalCount int64
	h.DB.Model(&db.Attendance{}).Where("student_id = ?", studentID).Where("present = ?", true).Count(&presentCount)
	h.DB.Model(&db.Attendance{}).Where("student_id = ?", studentID).Count(&totalCount)

	var percentage float64
	if totalCount > 0 {
		percentage = float64(presentCount) / float64(totalCount) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"data": attendance,
		"stats": gin.H{
			"present_days":          presentCount,
			"total_days":            totalCount,
			"attendance_percentage": percentage,
		},
	})
}

// GetMyAttendance godoc
// @Summary      Get my attendance
// @Description  Get attendance records for the authenticated student
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        start_date query     string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date   query     string  false  "End date (YYYY-MM-DD)"
// @Success      200        {object}  object{data=array,stats=object}
// @Failure      401        {object}  object{error=string}
// @Router       /attendance/my [get]
func (h *AttendanceHandler) GetMyAttendance(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var attendance []db.Attendance
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := h.DB.Where("student_id = ?", userID)

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	query.Order("date DESC").Find(&attendance)

	// Calculate stats
	var presentCount int64
	var totalCount int64
	h.DB.Model(&db.Attendance{}).Where("student_id = ?", userID).Where("present = ?", true).Count(&presentCount)
	h.DB.Model(&db.Attendance{}).Where("student_id = ?", userID).Count(&totalCount)

	var percentage float64
	if totalCount > 0 {
		percentage = float64(presentCount) / float64(totalCount) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"data": attendance,
		"stats": gin.H{
			"present_days":          presentCount,
			"total_days":            totalCount,
			"attendance_percentage": percentage,
		},
	})
}

// GetDailyAttendance godoc
// @Summary      Get daily attendance
// @Description  Get attendance records for a specific date
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        date   query     string  false  "Date (YYYY-MM-DD), defaults to today"
// @Success      200    {object}  object{date=string,data=array}
// @Failure      401    {object}  object{error=string}
// @Router       /attendance/daily [get]
func (h *AttendanceHandler) GetDailyAttendance(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var attendance []db.Attendance
	h.DB.Where("date = ?", date).Preload("Student").Find(&attendance)

	c.JSON(http.StatusOK, gin.H{
		"date": date,
		"data": attendance,
	})
}
