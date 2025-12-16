// Package main initializes and runs the hotel booking API server.
// It sets up all the necessary dependencies, routes, and starts the Gin web server.
package main

import (
	"fmt"
	"industry-api/db"
	"industry-api/internal/cache"
	"industry-api/internal/handler"
	"industry-api/internal/repository"
	"industry-api/internal/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main is the entry point of the application.
// It performs the following operations:
// 1. Loads environment variables from .env file
// 2. Initializes database connection
// 3. Initializes Redis cache
// 4. Creates all repositories (data access layer)
// 5. Creates all services (business logic layer)
// 6. Creates all handlers (HTTP request handlers)
// 7. Sets up all API routes
// 8. Starts the HTTP server on the configured port
func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()
	fmt.Println("Environment variables loaded successfully")

	// Initialize database connection
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initiate Database: %v", err)

	}
	// Ensure database connection is closed when the application exits
	defer db.Close()
	fmt.Println("Server is running...")

	// Initialize Redis cache client
	cache.Init()
	// Ensure cache connection is closed when the application exits
	defer cache.Close()
	// Validate that the database connection was successful
	if db.DB == nil {
		log.Fatal("Database connection failed")
	}

	// Initialize Gin router for handling HTTP requests
	router := gin.Default()

	// ========== User Management Setup ==========
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// ========== Room Management Setup ==========
	roomRepo := repository.NewRoomRepository(db.DB)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	// ========== Room Maintenance Setup ==========
	roomMaintenanceRepo := repository.NewRoomMaintenanceRepository(db.DB)
	roomMaintenanceService := service.NewRoomMaintenanceService(roomMaintenanceRepo)
	roomMaintenanceHandler := handler.NewRoomMaintenanceHandler(roomMaintenanceService)

	// ========== Booking Management Setup ==========
	bookingRepo := repository.NewBookingRepository(db.DB)
	bookingService := service.NewBookingService(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// ========== Payment Management Setup ==========
	paymentRepo := repository.NewPaymentRepository(db.DB)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// ========== Route Configuration ==========
	// All API routes are prefixed with /api/v1 for versioning
	v1 := router.Group("/api/v1")
	{
		// Authentication and user management routes
		users := v1.Group("/auth")
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.LoginUser)
		users.GET("/fetch-users", userHandler.GetUserList)
		users.GET("/fetch-user-by-id/:id", userHandler.GetUserByID)
		users.PUT("/update-user-status/:id", userHandler.UpdateUserStatus)

		// Room management routes
		rooms := v1.Group("/rooms")
		rooms.POST("/add", roomHandler.AddRoom)
		rooms.GET("/allRoomsList", roomHandler.GetRoomsList)
		rooms.GET("/availableRoomsList", roomHandler.GetAvailableRooms)

		// Room maintenance routes
		roomMaintenance := v1.Group("/roomMaintenance")
		roomMaintenance.POST("/add", roomMaintenanceHandler.AddRoomMaintenance)

		// Booking management routes
		booking := v1.Group("/bookings")
		booking.POST("/add", bookingHandler.AddBooking)

		// Payment processing routes
		payment := v1.Group("/payments")
		payment.POST("/initiate", paymentHandler.InitiatePayment)
		payment.PUT("/update-payment", paymentHandler.UpdatePayment)

	}

	// Get the port from environment variables or use default port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the HTTP server and listen for incoming requests
	router.Run(":" + port)

}
