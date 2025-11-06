package leaves

import (
	"log"
	"net/http"
	"strconv"

	"attendance-workflow/internal/dto"
	"attendance-workflow/internal/notifications"
	"attendance-workflow/pkg/db"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	DB db.GormDB
}

func NewLeaveHandler() *LeaveHandler {
	return &LeaveHandler{DB: db.DB}
}

// ApplyLeave godoc
// @Summary      Apply for leave
// @Description  Submit a leave request
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.ApplyLeaveRequest  true  "Leave request details"
// @Success      201      {object}  object{message=string,data=object}
// @Failure      400      {object}  object{error=string}
// @Failure      401      {object}  object{error=string}
// @Router       /leaves/apply [post]
func (h *LeaveHandler) ApplyLeave(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req dto.ApplyLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EndDate.Time.Before(req.StartDate.Time) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before start date"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	leaveReq := db.LeaveRequest{
		StudentID: uid,
		LeaveType: db.LeaveType(req.LeaveType),
		Reason:    req.Reason,
		StartDate: req.StartDate.Time,
		EndDate:   req.EndDate.Time,
		Status:    db.StatusPending,
	}

	if err := h.DB.Create(&leaveReq).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Leave request submitted successfully",
		"data":    leaveReq,
	})
}

// GetMyLeaves godoc
// @Summary      Get my leaves
// @Description  Get leave requests for the authenticated user
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number" default(1)
// @Param        limit   query     int     false  "Items per page" default(10)
// @Param        status  query     string  false  "Filter by status"
// @Success      200     {object}  object{data=array,page=int,limit=int,total=int64,total_pages=int64}
// @Failure      401     {object}  object{error=string}
// @Router       /leaves/my [get]
func (h *LeaveHandler) GetMyLeaves(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var leaves []db.LeaveRequest
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	status := c.Query("status")

	offset := (page - 1) * limit
	query := h.DB.Where("student_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	countQuery := h.DB.Model(&db.LeaveRequest{}).Where("student_id = ?", userID)
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	countQuery.Count(&total)
	query.Limit(limit).Offset(offset).Find(&leaves)

	c.JSON(http.StatusOK, gin.H{
		"data":        leaves,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// GetPendingLeaves godoc
// @Summary      Get pending leaves
// @Description  Get all pending leave requests (admin/faculty/warden only)
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number" default(1)
// @Param        limit   query     int     false  "Items per page" default(10)
// @Success      200     {object}  object{data=array,page=int,limit=int,total=int64,total_pages=int64}
// @Failure      401     {object}  object{error=string}
// @Router       /leaves/pending [get]
func (h *LeaveHandler) GetPendingLeaves(c *gin.Context) {
	role, _ := c.Get("role")

	var leaves []db.LeaveRequest
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	offset := (page - 1) * limit
	query := h.DB.Where("status = ?", db.StatusPending)

	if role == "faculty" || role == "warden" {
		// Faculty and Warden can see all pending
		query = query.Preload("Student")
	}

	var total int64
	query.Count(&total)
	query.Limit(limit).Offset(offset).Find(&leaves)

	c.JSON(http.StatusOK, gin.H{
		"data":        leaves,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// ApproveLeave godoc
// @Summary      Approve or reject leave
// @Description  Approve or reject a leave request (admin/faculty/warden only)
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int                    true  "Leave ID"
// @Param        request body      dto.ApproveLeaveRequest true "Approval decision"
// @Success      200     {object}  object{message=string,data=object}
// @Failure      400     {object}  object{error=string}
// @Failure      404     {object}  object{error=string}
// @Failure      401     {object}  object{error=string}
// @Router       /leaves/{id}/approve [put]
func (h *LeaveHandler) ApproveLeave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID"})
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req dto.ApproveLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var leave db.LeaveRequest
	if err := h.DB.First(&leave, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		return
	}

	if leave.Status != db.StatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave request already processed"})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	approvedBy := uint(uid)
	leave.ApprovedBy = &approvedBy
	leave.Remarks = req.Remarks

	if req.Approved {
		leave.Status = db.StatusApproved
	} else {
		leave.Status = db.StatusRejected
	}

	if err := h.DB.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leave request"})
		return
	}

	// Queue async notification
	notifService := notifications.GetNotificationService()
	if err := notifService.QueueLeaveStatusNotification(c.Request.Context(), notifications.LeaveStatusUpdatePayload{
		LeaveID:    uint(id),
		StudentID:  leave.StudentID,
		Status:     string(leave.Status),
		ApprovedBy: req.Remarks,
		Remarks:    req.Remarks,
	}); err != nil {
		// Log error but don't fail the request
		log.Printf("Failed to queue notification: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Leave request " + string(leave.Status),
		"data":    leave,
	})
}

// GetAllLeaves godoc
// @Summary      Get all leaves
// @Description  Get all leave requests (admin only)
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number" default(1)
// @Param        limit   query     int     false  "Items per page" default(10)
// @Param        status  query     string  false  "Filter by status"
// @Success      200     {object}  object{data=array,page=int,limit=int,total=int64,total_pages=int64}
// @Failure      401     {object}  object{error=string}
// @Router       /leaves [get]
func (h *LeaveHandler) GetAllLeaves(c *gin.Context) {
	var leaves []db.LeaveRequest
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	status := c.Query("status")

	offset := (page - 1) * limit
	query := h.DB.Preload("Student")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)
	query.Limit(limit).Offset(offset).Find(&leaves)

	c.JSON(http.StatusOK, gin.H{
		"data":        leaves,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// DeleteLeave godoc
// @Summary      Delete leave
// @Description  Delete a leave request (admin only)
// @Tags         leaves
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int     true  "Leave ID"
// @Success      200     {object}  object{message=string}
// @Failure      404     {object}  object{error=string}
// @Failure      401     {object}  object{error=string}
// @Router       /leaves/{id} [delete]
func (h *LeaveHandler) DeleteLeave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leave ID"})
		return
	}

	if err := h.DB.Delete(&db.LeaveRequest{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave request deleted successfully"})
}
