package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/mocbotau/mocbot-archive/internal/database"
	"github.com/mocbotau/mocbot-archive/internal/handlers"
	"github.com/mocbotau/mocbot-archive/internal/middleware"
	"github.com/mocbotau/mocbot-archive/internal/utils"
)

const (
	port = "9000"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading the .env file: %v... continuing", err)
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDBName := os.Getenv("MYSQL_DATABASE")
	passwordFilePath := os.Getenv("MYSQL_PASSWORD_FILE")

	secretManager := utils.NewSecretManager(0)

	mysqlPassword, err := secretManager.GetSecretFromFile(passwordFilePath)
	if err != nil {
		secretManager.Stop()
		log.Fatalf("Failed to read MySQL password from file: %v", err)
	}

	db, err := database.NewMySQL(mysqlUser, mysqlPassword, mysqlHost, mysqlDBName)
	if err != nil {
		secretManager.Stop()
		log.Fatalf("Failed to initialize database: %v", err)
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		secretManager.Stop()
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	defer func() {
		secretManager.Stop()

		_ = sqlDB.Close()
	}()

	handler := handlers.NewHandler(db)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	setupRoutes(r, handler, secretManager)

	log.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Panic("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine, handler *handlers.Handler, secretManager *utils.SecretManager) {
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1", middleware.EnsureToken(secretManager, os.Getenv("API_KEY_FILE")))
	{
		sessions := api.Group("/sessions")
		{
			sessions.GET("", handler.GetSessionsByGuilds)     // GET /api/v1/sessions?guildIds=123,456&limit=50
			sessions.POST("", handler.StartSession)           // POST /api/v1/sessions
			sessions.GET("/:sessionId", handler.GetSession)   // GET /api/v1/sessions/{sessionId}
			sessions.PATCH("/:sessionId", handler.EndSession) // PATCH /api/v1/sessions/{sessionId}

			sessions.GET("/:sessionId/tracks", handler.GetTracksBySession) // GET /api/v1/sessions/{sessionId}/tracks
			sessions.POST("/:sessionId/tracks", handler.StartTrack)        // POST /api/v1/sessions/{sessionId}/tracks
		}

		tracks := api.Group("/tracks")
		{
			tracks.GET("/:trackPlayId", handler.GetTrackPlay)                // GET /api/v1/tracks/{trackPlayId}
			tracks.PATCH("/:trackPlayId", handler.EndTrack)                  // PATCH /api/v1/tracks/{trackPlayId}
			tracks.GET("/:trackPlayId/listeners", handler.GetTrackListeners) // GET /api/v1/tracks/{trackPlayId}/listeners
		}

		guilds := api.Group("/guilds")
		{
			guilds.GET("/:guildId/tracks/recent", handler.GetRecentTracksByGuild) // GET /api/v1/guilds/{guildId}/tracks/recent?limit=50
		}

		users := api.Group("/users")
		{
			users.GET("/:userId/tracks/recent", handler.GetRecentTracksByUser) // GET /api/v1/users/{userId}/tracks/recent?limit=50
		}
	}
}
