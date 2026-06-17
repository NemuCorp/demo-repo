package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/handler"
	"github.com/NemuCorp/demo-repo/server/logger"
)

var (
	database *db.DB
)

func init() {
	logger.Init(logger.ModeDevelopment)

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/demorepo?sslmode=disable"
	}

	var err error
	database, err = db.Open(connStr)
	if err != nil {
		logger.Error.Fatalf("Failed to open database: %v", err)
	}
}

func main() {
	defer database.Close()

	if len(os.Args) > 1 {
		runCmd(database)
		return
	}

	if err := database.PrepareStatements(); err != nil {
		logger.Error.Fatalf("Failed to initialize database: %v", err)
	}

	authHandler := handler.NewAuthHandler(database.Auth)
	cartHandler := handler.NewCartHandler(database.Cart)
	productHandler := handler.NewProductHandler(database.Product)
	trackingHandler := handler.NewTrackingHandler(database.Tracking)
	authMiddleware := handler.AuthMiddleware(database.Auth)

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authMiddleware, authHandler.Logout)
		}

		products := api.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.Get)
			products.POST("", productHandler.Create)
		}

		cart := api.Group("/cart", authMiddleware)
		{
			cart.GET("", cartHandler.View)
			cart.POST("", cartHandler.Add)
			cart.PUT("/:productId", cartHandler.Update)
			cart.DELETE("/:productId", cartHandler.Remove)
		}

		tracking := api.Group("/track")
		{
			tracking.POST("", trackingHandler.Track)
		}

		admin := api.Group("/admin", authMiddleware)
		{
			admin.GET("/stats", trackingHandler.Dashboard)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info.Printf("Server starting on :%s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Error.Fatalf("Failed to start server: %v", err)
	}
}
