package users

import (
	"net/http"
	"strconv"

	"attendance-workflow/pkg/db"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB db.GormDB
}

func NewUserHandler() *UserHandler {
	return &UserHandler{DB: db.DB}
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get paginated list of all users with optional filters
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number" default(1)
// @Param        limit   query     int     false  "Items per page" default(10)
// @Param        role    query     string  false  "Filter by role"
// @Param        dept    query     string  false  "Filter by department"
// @Success      200     {object}  object{data=array,page=int,limit=int,total=int64,total_pages=int64}
// @Failure      401     {object}  object{error=string}
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var users []db.User
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	role := c.Query("role")
	dept := c.Query("dept")

	offset := (page - 1) * limit
	query := h.DB.Model(&db.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if dept != "" {
		query = query.Where("dept = ?", dept)
	}

	var total int64
	query.Count(&total)
	query.Limit(limit).Offset(offset).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data":        users,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get user details by user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int     true  "User ID"
// @Success      200     {object}  object{data=object}
// @Failure      404     {object}  object{error=string}
// @Failure      401     {object}  object{error=string}
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user db.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int     true  "User ID"
// @Param        request body      object  true  "User update data"
// @Success      200     {object}  object{message=string,data=object}
// @Failure      400     {object}  object{error=string}
// @Failure      404     {object}  object{error=string}
// @Failure      401     {object}  object{error=string}
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user db.User
	if err := h.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "data": user})
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int     true  "User ID"
// @Success      200     {object}  object{message=string}
// @Failure      404     {object}  object{error=string}
// @Failure      401     {object}  object{error=string}
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.DB.Delete(&db.User{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
