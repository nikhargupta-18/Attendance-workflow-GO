package api

import (
	"attendance-workflow/internal/analytics"
	"attendance-workflow/internal/attendance"
	"attendance-workflow/internal/auth"
	"attendance-workflow/internal/leaves"
	"attendance-workflow/internal/notifications"
	"attendance-workflow/internal/users"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize handlers
	authHandler := auth.NewAuthHandler()
	userHandler := users.NewUserHandler()
	leaveHandler := leaves.NewLeaveHandler()
	attendanceHandler := attendance.NewAttendanceHandler()
	notificationHandler := notifications.NewNotificationHandler()
	analyticsHandler := analytics.NewAnalyticsHandler()

	// Public routes
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(auth.AuthMiddleware())
	{
		// Users
		usersGroup := protected.Group("/users")
		{
			usersGroup.GET("", auth.RoleMiddleware("admin"), userHandler.GetAllUsers)
			usersGroup.GET("/:id", userHandler.GetUserByID)
			usersGroup.PUT("/:id", userHandler.UpdateUser)
			usersGroup.DELETE("/:id", auth.RoleMiddleware("admin"), userHandler.DeleteUser)
		}

		// Leaves
		leavesGroup := protected.Group("/leaves")
		{
			leavesGroup.POST("/apply", auth.RoleMiddleware("student"), leaveHandler.ApplyLeave)
			leavesGroup.GET("/my", leaveHandler.GetMyLeaves)
			leavesGroup.GET("/pending", auth.RoleMiddleware("faculty", "warden", "admin"), leaveHandler.GetPendingLeaves)
			leavesGroup.PUT("/:id/approve", auth.RoleMiddleware("faculty", "warden", "admin"), leaveHandler.ApproveLeave)
			leavesGroup.GET("", auth.RoleMiddleware("admin"), leaveHandler.GetAllLeaves)
			leavesGroup.DELETE("/:id", auth.RoleMiddleware("admin"), leaveHandler.DeleteLeave)
		}

		// Attendance
		attendanceGroup := protected.Group("/attendance")
		{
			attendanceGroup.POST("/mark", attendanceHandler.MarkAttendance)
			attendanceGroup.GET("/student/:id", attendanceHandler.GetStudentAttendance)
			attendanceGroup.GET("/my", auth.RoleMiddleware("student"), attendanceHandler.GetMyAttendance)
			attendanceGroup.GET("/daily", attendanceHandler.GetDailyAttendance)
		}

		// Notifications
		notificationsGroup := protected.Group("/notifications")
		{
			notificationsGroup.GET("/my", notificationHandler.GetMyNotifications)
			notificationsGroup.PUT("/:id/read", notificationHandler.MarkAsRead)
			notificationsGroup.GET("/unread-count", notificationHandler.GetUnreadCount)
		}

		// Analytics
		analyticsGroup := protected.Group("/analytics")
		analyticsGroup.Use(auth.RoleMiddleware("admin"))
		{
			analyticsGroup.GET("/dashboard", analyticsHandler.GetDashboardStats)
			analyticsGroup.GET("/leave-breakdown", analyticsHandler.GetLeaveBreakdown)
			analyticsGroup.GET("/department", analyticsHandler.GetDepartmentStats)
			analyticsGroup.GET("/absentees", analyticsHandler.GetFrequentAbsentees)
		}
	}

	return router
}

