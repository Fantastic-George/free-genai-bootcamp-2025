package validation

// CreateStudySessionRequest represents the request to create a study session
type CreateStudySessionRequest struct {
	GroupID         int64 `json:"group_id" binding:"required,min=1"`
	StudyActivityID int64 `json:"study_activity_id" binding:"required,min=1"`
}

// AddWordReviewRequest represents the request to add a word review
type AddWordReviewRequest struct {
	Correct bool `json:"correct" binding:"required"`
}

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

// GetDefaultPagination returns default pagination values if not provided
func GetDefaultPagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
