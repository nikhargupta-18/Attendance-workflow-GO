// @title           Attendance Workflow API
// @version         1.0
// @description     API for managing attendance, leaves, and notifications in an educational institution
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@attendance-workflow.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"
	"log"
	"os"

	docs "attendance-workflow/docs"
	"attendance-workflow/internal/api"
	"attendance-workflow/internal/auth"
	"attendance-workflow/internal/notifications"
	"attendance-workflow/pkg/config"
	"attendance-workflow/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.Load()

	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.GinMode)

	// Connect to database
	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default admin user if not exists
	createDefaultAdmin()

	// Initialize notification service
	notifService := notifications.GetNotificationService()
	defer notifService.Stop()

	// Start cron service
	cronService := notifications.NewCronService()
	go func() {
		if err := cronService.Start(); err != nil {
			log.Printf("Failed to start cron service: %v", err)
		}
	}()
	defer cronService.Stop()

	// Setup routes
	router := api.SetupRoutes()

	// Start server
	port := config.AppConfig.Server.Port
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	// Configure Swagger
	docs.SwaggerInfo.Title = "Attendance Workflow API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.Schemes = []string{"http"}

	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("API: http://localhost:" + port)
	fmt.Println("Health Check: http://localhost:" + port + "/health")
	fmt.Println("Swagger UI: http://localhost:" + port + "/swagger/index.html")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func createDefaultAdmin() {
	var count int64
	db.DB.Model(&db.User{}).Count(&count)

	if count == 0 {
		hashedPassword, err := auth.HashPassword("admin123")
		if err != nil {
			log.Println("Failed to hash admin password:", err)
			return
		}
		admin := db.User{
			Name:     "Admin User",
			Email:    "admin@university.edu",
			Password: hashedPassword,
			Role:     db.RoleAdmin,
			Dept:     "Administration",
		}

		if err := db.DB.Create(&admin).Error; err == nil {
			log.Println("Default admin user created (admin@university.edu / admin123)")
		} else {
			log.Println("Failed to create default admin user:", err)
		}
	}
}
