package main

import (
	"log"
	"path/filepath"

	"pengyou-chinese/backend/internal/handlers"
	"pengyou-chinese/backend/internal/middleware"
	"pengyou-chinese/backend/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Initialize database service
	dbPath := filepath.Join(".", "words.db")
	db, err := service.NewDBService(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize handlers
	dashboardHandler := handlers.NewDashboardHandler(db)
	wordsHandler := handlers.NewWordsHandler(db)
	groupsHandler := handlers.NewGroupsHandler(db)
	studyHandler := handlers.NewStudyHandler(db)

	// Create a default Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api")
	{
		// Dashboard routes
		api.GET("/dashboard/last_study_session", dashboardHandler.GetLastStudySession)
		api.GET("/dashboard/study_progress", dashboardHandler.GetStudyProgress)
		api.GET("/dashboard/quick-stats", dashboardHandler.GetQuickStats)

		// Words routes
		api.GET("/words", wordsHandler.GetWords)
		api.GET("/words/:id", wordsHandler.GetWord)
		api.POST("/study_sessions/:id/words/:word_id/review", wordsHandler.AddWordReview)

		// Groups routes
		api.GET("/groups", groupsHandler.GetGroups)
		api.GET("/groups/:id", groupsHandler.GetGroup)
		api.GET("/groups/:id/words", groupsHandler.GetGroupWords)

		// Study sessions routes
		api.GET("/study_sessions", studyHandler.GetStudySessions)
		api.GET("/study_sessions/:id", studyHandler.GetStudySession)
		api.GET("/study_sessions/:id/words", studyHandler.GetStudySessionWords)
		api.POST("/study_sessions", studyHandler.CreateStudySession)

		// Study activities routes
		api.GET("/study_activities/:id", studyHandler.GetStudyActivity)
		api.GET("/study_activities/:id/study_sessions", studyHandler.GetStudyActivitySessions)
	}

	// Start the server
	log.Printf("Starting server on :8080")
	log.Fatal(router.Run(":8080"))
}
