package handlers

import (
	"net/http"
	"strconv"

	"pengyou-chinese/backend/internal/service"
	"pengyou-chinese/backend/internal/validation"

	"github.com/gin-gonic/gin"
)

// StudyHandler handles study session and activity related routes
type StudyHandler struct {
	db *service.DBService
}

// NewStudyHandler creates a new study handler
func NewStudyHandler(db *service.DBService) *StudyHandler {
	return &StudyHandler{db: db}
}

// GetStudySessions returns a paginated list of study sessions
func (h *StudyHandler) GetStudySessions(c *gin.Context) {
	var pagination validation.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	page, pageSize := validation.GetDefaultPagination(pagination.Page, pagination.PageSize)
	sessions, total, err := h.db.GetStudySessions(page, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": sessions,
		"pagination": gin.H{
			"current_page":   page,
			"total_pages":    (total + pageSize - 1) / pageSize,
			"total_items":    total,
			"items_per_page": pageSize,
		},
	})
}

// GetStudySession returns a single study session by ID
func (h *StudyHandler) GetStudySession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := h.db.GetStudySession(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudySessionWords returns words reviewed in a study session
func (h *StudyHandler) GetStudySessionWords(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var pagination validation.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	page, pageSize := validation.GetDefaultPagination(pagination.Page, pagination.PageSize)
	words, total, err := h.db.GetStudySessionWords(sessionID, page, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": words,
		"pagination": gin.H{
			"current_page":   page,
			"total_pages":    (total + pageSize - 1) / pageSize,
			"total_items":    total,
			"items_per_page": pageSize,
		},
	})
}

// GetStudyActivity returns a study activity by ID
func (h *StudyHandler) GetStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := h.db.GetStudyActivity(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if activity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// GetStudyActivitySessions returns study sessions for a specific activity
func (h *StudyHandler) GetStudyActivitySessions(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	var pagination validation.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	page, pageSize := validation.GetDefaultPagination(pagination.Page, pagination.PageSize)
	sessions, total, err := h.db.GetStudyActivitySessions(activityID, page, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": sessions,
		"pagination": gin.H{
			"current_page":   page,
			"total_pages":    (total + pageSize - 1) / pageSize,
			"total_items":    total,
			"items_per_page": pageSize,
		},
	})
}

// CreateStudySession creates a new study session
func (h *StudyHandler) CreateStudySession(c *gin.Context) {
	var request validation.CreateStudySessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, err := h.db.CreateStudySession(request.GroupID, request.StudyActivityID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, session)
}
