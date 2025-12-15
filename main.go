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

func main() {
	_ = godotenv.Load()
	fmt.Println("Environment variables loaded successfully")

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initiate Database: %v", err)

	}
	defer db.Close()
	fmt.Println("Server is running...")

	cache.Init()
	defer cache.Close()
	if db.DB == nil {
		log.Fatal("Database connection failed")
	}

	router := gin.Default()

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	roomRepo := repository.NewRoomRepository(db.DB)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	roomMaintenanceRepo := repository.NewRoomMaintenanceRepository(db.DB)
	roomMaintenanceService := service.NewRoomMaintenanceService(roomMaintenanceRepo)
	roomMaintenanceHandler := handler.NewRoomMaintenanceHandler(roomMaintenanceService)

	bookingRepo := repository.NewBookingRepository(db.DB)
	bookingService := service.NewBookingService(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	paymentRepo := repository.NewPaymentRepository(db.DB)
	paymentService := service.NewPaymentRepository(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/auth")
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.LoginUser)
		users.GET("/fetch-users", userHandler.GetUserList)
		users.GET("/fetch-user-by-id/:id", userHandler.GetUserByID)
		users.PUT("/update-user-status/:id", userHandler.UpdateUserStatus)

		rooms := v1.Group("/rooms")
		rooms.POST("/add", roomHandler.AddRoom)
		rooms.GET("/allRoomsList", roomHandler.GetRoomsList)
		rooms.GET("/availableRoomsList", roomHandler.GetAvailableRooms)

		roomMaintenance := v1.Group("/roomMaintenance")
		roomMaintenance.POST("/add", roomMaintenanceHandler.AddRoomMaintenance)

		booking := v1.Group("/bookings")
		booking.POST("/add", bookingHandler.AddBooking)

		payment := v1.Group("/payments")
		payment.POST("/initiate", paymentHandler.InitiatePayment)
		payment.PUT("/update-payment", paymentHandler.UpdatePayment)

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)

}
