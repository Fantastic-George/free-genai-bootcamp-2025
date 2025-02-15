package handlers

import (
	"net/http"
	"strconv"

	"pengyou-chinese/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// GroupsHandler handles group-related routes
type GroupsHandler struct {
	db *service.DBService
}

// NewGroupsHandler creates a new groups handler
func NewGroupsHandler(db *service.DBService) *GroupsHandler {
	return &GroupsHandler{db: db}
}

// GetGroups returns a paginated list of groups
func (h *GroupsHandler) GetGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	groups, total, err := h.db.GetGroups(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": groups,
		"pagination": gin.H{
			"current_page":   page,
			"total_pages":    (total + pageSize - 1) / pageSize,
			"total_items":    total,
			"items_per_page": pageSize,
		},
	})
}

// GetGroup returns a single group by ID
func (h *GroupsHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := h.db.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupWords returns words for a specific group
func (h *GroupsHandler) GetGroupWords(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))

	words, total, err := h.db.GetGroupWords(groupID, page, pageSize)
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
