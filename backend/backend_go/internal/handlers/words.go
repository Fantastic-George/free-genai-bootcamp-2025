package handlers

import (
	"net/http"
	"strconv"

	"pengyou-chinese/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// WordsHandler handles word-related routes
type WordsHandler struct {
	db *service.DBService
}

// NewWordsHandler creates a new words handler
func NewWordsHandler(db *service.DBService) *WordsHandler {
	return &WordsHandler{db: db}
}

// GetWords returns a paginated list of words
func (h *WordsHandler) GetWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	words, total, err := h.db.GetWords(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// GetWord returns a single word by ID
func (h *WordsHandler) GetWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := h.db.GetWord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	c.JSON(http.StatusOK, word)
}

// AddWordReview adds a review for a word
func (h *WordsHandler) AddWordReview(c *gin.Context) {
	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var review struct {
		Correct bool `json:"correct" binding:"required"`
	}

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.db.AddWordReview(wordID, sessionID, review.Correct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"word_id":          wordID,
		"study_session_id": sessionID,
		"correct":          review.Correct,
	})
}
