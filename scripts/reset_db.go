package main

import (
	"attendance-workflow/pkg/config"
	"attendance-workflow/pkg/db"
	"fmt"
	"log"
	"time"
)

func main() {
	// Load configuration
	config.Load()

	// Connect to database
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get database instance
	database := db.DB

	// Drop all tables
	log.Println("Dropping all tables...")
	if err := database.Migrator().DropTable(
		&db.User{},
		&db.LeaveRequest{},
		&db.Attendance{},
		&db.Notification{},
	); err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	// Recreate tables
	log.Println("Recreating tables...")
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to recreate tables: %v", err)
	}

	log.Println("Creating sample data...")

	// Create users with different roles
	admin := &db.User{
		Name:     "admin",
		Email:    "admin@college.edu",
		Password: "$2a$10$DemoHashedPasswordAA", // Demo password: admin123
		Role:     db.RoleAdmin,
	}

	warden := &db.User{
		Name:     "warden",
		Email:    "warden@college.edu",
		Password: "$2a$10$DemoHashedPasswordBB", // Demo password: warden123
		Role:     db.RoleWarden,
	}

	faculty := &db.User{
		Name:     "faculty",
		Email:    "faculty@college.edu",
		Password: "$2a$10$DemoHashedPasswordCC", // Demo password: faculty123
		Role:     db.RoleFaculty,
		Dept:     "Computer Science",
	}

	students := []*db.User{
		{
			Name:     "stud1",
			Email:    "stud1@college.edu",
			Password: "$2a$10$DemoHashedPasswordDD", // Demo password: student123
			Role:     db.RoleStudent,
			Dept:     "Computer Science",
		},
		{
			Name:     "stud2",
			Email:    "stud2@college.edu",
			Password: "$2a$10$DemoHashedPasswordEE", // Demo password: student123
			Role:     db.RoleStudent,
			Dept:     "Computer Science",
		},
	}

	// Create all users
	users := []*db.User{admin, warden, faculty}
	users = append(users, students...)

	for _, user := range users {
		if err := database.Create(user).Error; err != nil {
			log.Fatalf("Failed to create user %s: %v", user.Name, err)
		}
	}

	// Create sample leave requests
	leaveRequests := []*db.LeaveRequest{
		{
			StudentID:  students[0].ID,
			LeaveType:  db.LeaveTypeMedical,
			Reason:     "Medical appointment",
			StartDate:  time.Now().AddDate(0, 0, 1),
			EndDate:    time.Now().AddDate(0, 0, 2),
			Status:     db.StatusApproved,
			ApprovedBy: &warden.ID,
		},
		{
			StudentID: students[1].ID,
			LeaveType: db.LeaveTypePersonal,
			Reason:    "Family function",
			StartDate: time.Now().AddDate(0, 0, 5),
			EndDate:   time.Now().AddDate(0, 0, 7),
			Status:    db.StatusPending,
		},
	}

	for _, leave := range leaveRequests {
		if err := database.Create(leave).Error; err != nil {
			log.Fatalf("Failed to create leave request: %v", err)
		}
	}

	// Create sample attendance records
	startDate := time.Now().AddDate(0, 0, -10)
	for _, student := range students {
		for i := 0; i < 10; i++ {
			date := startDate.AddDate(0, 0, i)
			attendance := &db.Attendance{
				StudentID: student.ID,
				Date:      date,
				Present:   i%2 == 0, // Alternate between present and absent
				MarkedBy:  faculty.ID,
			}
			if err := database.Create(attendance).Error; err != nil {
				log.Fatalf("Failed to create attendance record: %v", err)
			}
		}
	}

	// Create sample notifications
	notifications := []*db.Notification{
		{
			UserID:  students[0].ID,
			Type:    "leave_approved",
			Title:   "Leave Request Approved",
			Message: "Your medical leave request has been approved",
			IsRead:  false,
		},
		{
			UserID:  warden.ID,
			Type:    "new_leave_request",
			Title:   "New Leave Request",
			Message: "New leave request from Bob Johnson",
			IsRead:  false,
		},
	}

	for _, notification := range notifications {
		if err := database.Create(notification).Error; err != nil {
			log.Fatalf("Failed to create notification: %v", err)
		}
	}

	fmt.Println("\nDatabase populated successfully!")
	fmt.Println("\nDemo Accounts:")
	fmt.Println("==========================")

	// Print credentials for all created users
	demoUsers := []*db.User{admin, warden, faculty}
	demoUsers = append(demoUsers, students...)

	for _, user := range demoUsers {
		fmt.Printf("\n%s User: %s", user.Role, user.Name)
		fmt.Printf("\nEmail: %s", user.Email)
		switch user.Role {
		case db.RoleAdmin:
			fmt.Printf("\nPassword: admin123")
		case db.RoleWarden:
			fmt.Printf("\nPassword: warden123")
		case db.RoleFaculty:
			fmt.Printf("\nPassword: faculty123")
		case db.RoleStudent:
			fmt.Printf("\nPassword: student123")
		}
		fmt.Println("\n-----------------------------------------------------------------")
	}

	fmt.Println("\nNote: Use these credentials to log in via the API and get your bearer token.")
	fmt.Printf("Example: curl -X POST http://localhost:8080/api/v1/auth/login -H \"Content-Type: application/json\" -d '{\"email\":\"admin@college.edu\",\"password\":\"admin123\"}'\n")
}
