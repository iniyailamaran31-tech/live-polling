package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"live-polling-backend/handlers"
    "live-polling-backend/middleware"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")

	if mongoURI == "" {
		panic("MONGODB_URI is not set")
	}

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoURI),
	)
	if err != nil {
		panic("Failed to create MongoDB client: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	println("MongoDB connected successfully!")

	// Database
	databaseName := os.Getenv("MONGODB_DATABASE")

	if databaseName == "" {
		databaseName = "livepoll"
	}

	db := client.Database(databaseName)
	// Redis connection
redisURL := os.Getenv("REDIS_URL")

if redisURL == "" {
	panic("REDIS_URL is not set")
}

redisOptions, err := redis.ParseURL(redisURL)
if err != nil {
	panic("Invalid Redis URL: " + err.Error())
}

redisClient := redis.NewClient(redisOptions)

redisCtx, redisCancel := context.WithTimeout(
	context.Background(),
	10*time.Second,
)
defer redisCancel()

err = redisClient.Ping(redisCtx).Err()
if err != nil {
	panic("Failed to connect to Redis: " + err.Error())
}

println("Redis connected successfully!")

	// Gin router
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "http://localhost:5173",
        "https://live-polling-wheat-theta.vercel.app",
    },
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: false,
}))

	// Health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Live Polling Backend is running!",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "MongoDB connected",
		})
	})

	// Authentication
	authHandler := &handlers.AuthHandler{
		DB: db,
	}

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
	// Poll routes
pollHandler := &handlers.PollHandler{
	DB:    db,
	Redis: redisClient,
}
realtimeHandler := &handlers.RealtimeHandler{
	Redis: redisClient,
}
// Protected route - only logged-in users can create polls
// Public poll routes
router.GET("/polls/:id", pollHandler.GetPoll)
router.POST("/polls/:id/vote", pollHandler.Vote)
// WebSocket live updates
router.GET("/ws/polls/:id", realtimeHandler.PollUpdates)

// Protected poll creation
protectedPolls := router.Group("/polls")
protectedPolls.Use(middleware.AuthRequired())
{
	protectedPolls.POST("", pollHandler.CreatePoll)
}

	// Start server
	router.Run(":8080")
}