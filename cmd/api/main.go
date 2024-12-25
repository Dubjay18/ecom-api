// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Dubjay18/ecom-api/docs" // This is for swagger
	"github.com/Dubjay18/ecom-api/internal/config"
	"github.com/Dubjay18/ecom-api/internal/container"
	"github.com/Dubjay18/ecom-api/internal/handler"
	"github.com/Dubjay18/ecom-api/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           E-commerce API
// @version         1.0
// @description     A RESTful API for an e-commerce application
// @host           localhost:8080
// @BasePath       /api/v1

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	log.Printf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	// Set Gin mode
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize container (DB, repositories, and services)
	c, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatal("cannot initialize container:", err)
	}
	defer c.Close()

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes group
	api := router.Group("/api/v1")

	loggerInit := config.InitLog()
	api.Use(middleware.LoggerMiddleware(loggerInit))

	// Initialize handlers
	handler.NewUserHandler(api, c.UserService, loggerInit, cfg.JWT.SecretKey)
	handler.NewProductHandler(api, c.ProductService)
	handler.NewOrderHandler(api, c.OrderService)

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         cfg.Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		log.Println("shutting down server...")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Start server

	log.Printf("server is starting on %s\n", cfg.Server.GetServerAddress())
	fmt.Println(`
 /$$  / $$            $$                       
| $$  | $$          | $$                       
| $$  | $$  /$$$$$$ | $$$$$$$   /$$$$$$        
| $$$$$$$$|/$$__  $$| $$__  $$ /$$__  $$       
| $$__  $$| $$$$$$$$| $$  \ $$| $$$$$$$$       
| $$  | $$| $$_____/| $$  | $$| $$_____/       
| $$  | $$|  $$$$$$$| $$  | $$|  $$$$$$$       
|__/  |__/ \_______/|__/  |__/ \_______/      

 WELCOME TO JAY'S SERVER!
`)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("cannot start server:", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Println("server stopped")
}
