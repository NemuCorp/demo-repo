package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
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

	if err := handler.SeedAdminUser(database.Auth); err != nil {
		logger.Error.Fatalf("Failed to seed admin user: %v", err)
	}

	authHandler := handler.NewAuthHandler(database.Auth)
	cartHandler := handler.NewCartHandler(database.Cart)
	productHandler := handler.NewProductHandler(database.Product)
	trackingHandler := handler.NewTrackingHandler(database.Tracking)
	authMiddleware := handler.AuthMiddleware(database.Auth)
	adminMiddleware := handler.AdminMiddleware()

	r := gin.Default()

	allowedOrigins := corsAllowedOrigins()
	if len(allowedOrigins) > 0 {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     allowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authMiddleware, authHandler.Logout)
			auth.GET("/me", authMiddleware, authHandler.Me)
		}

		products := api.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.Get)
			products.POST("", authMiddleware, adminMiddleware, productHandler.Create)
			products.PUT("/:id", authMiddleware, adminMiddleware, productHandler.Update)
			products.DELETE("/:id", authMiddleware, adminMiddleware, productHandler.Delete)
		}

		cart := api.Group("/cart", authMiddleware)
		{
			cart.GET("", cartHandler.View)
			cart.POST("", cartHandler.Add)
			cart.PUT("/:productId", cartHandler.Update)
			cart.DELETE("/:productId", cartHandler.Remove)
		}

		tracking := api.Group("/track", handler.OptionalAuthMiddleware(database.Auth))
		{
			tracking.POST("", trackingHandler.Track)
		}

		admin := api.Group("/admin", authMiddleware, adminMiddleware)
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

func corsAllowedOrigins() []string {
	val := os.Getenv("CORS_ALLOWED_ORIGINS")
	if val == "" {
		val = "http://localhost:3000,http://admin.localhost:3000"
	}
	var origins []string
	for _, o := range strings.Split(val, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}
