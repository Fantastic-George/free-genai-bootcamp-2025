package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pengyou-chinese/backend/internal/service"
)

// DashboardHandler handles dashboard-related routes
type DashboardHandler struct {
	db *service.DBService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(db *service.DBService) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetLastStudySession returns information about the most recent study session
func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	session, err := h.db.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No study sessions found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudyProgress returns study progress statistics
func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	progress, err := h.db.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetQuickStats returns quick overview statistics
func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	stats, err := h.db.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
} 